package pritunl

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

// UpdatePublicAccessAddress 更新系统对外提供的公网地址，这个地址会被客户端连接配置文件使用，只需要服务端返回200即可
func UpdatePublicAccessAddress(c *Client, newAddress string) (*http.Response, error) {
	if net.ParseIP(newAddress) == nil {
		return nil, errors.New("地址格式不合法")
	}
	opts := RequestOpts{
		JSONBody: map[string]string{
			"public_address": newAddress,
		},
	}
	return c.Request("put", getServerSettingsPath(), &opts)
}

// AdminUser 管理账号配置
type AdminUser struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	AuthApi   bool   `json:"auth_api"`
	Token     string `json:"token"`
	Secret    string `json:"secret"`
	SuperUser bool   `json:"super_user"`
}

// GetAdminUserList 获取管理员账号列表
func GetAdminUserList(c *Client) ([]AdminUser, error) {
	var adminUsers []AdminUser
	// 设置json响应数据结构
	opts := RequestOpts{
		JSONResponse: &adminUsers,
	}
	if _, err := c.Request("get", getAdminListPath(), &opts); err != nil {
		return nil, err
	}
	return adminUsers, nil
}

// UpdateAdminUserAuthConfig 更新指定管理员账号配置, 如果adminUser的token和secret传了值，不管传什么值，服务端都是任意更新
// 响应体是adminUser
func UpdateAdminUserAuthConfig(c *Client, adminUser AdminUser) (*AdminUser, error) {
	adminUser.Token = "newToken"
	adminUser.Secret = "newSecret"

	opts := RequestOpts{
		JSONBody:     adminUser,
		JSONResponse: &adminUser,
	}
	if _, err := c.Request("put", getUpdateAdminUserPath(adminUser.Id), &opts); err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// VpnServer vpn server实例配置
type VpnServer struct {
	Name           string `json:"name,omitempty"`    // 不给的话由本包自动生成
	Id             string `json:"id,omitempty"`      // 创建server时不需要传递此参数
	Network        string `json:"network,omitempty"` // 不给的话由本包自动生成一个，必须满足[10,172,192].[0-255,16-31,168].[0-255].0/[8-24]
	Port           int    `json:"port,omitempty"`
	Protocol       string `json:"protocol,omitempty"` // 不给的话默认为udp
	Cipher         string `json:"cipher,omitempty"`   // 不给的话默认为aes128
	Hash           string `json:"hash,omitempty"`     // 不给的话默认为sha1
	RestrictRoutes bool   `json:"restrict_routes,omitempty"`
	NetworkMode    string `json:"network_mode,omitempty"` // 不给的话默认为tunnel
	Status         string `json:"status,omitempty"`       // 服务的状态
}

// CreateVpnServer 创建一个新的vpn server, 返回值是ServerCreateConfig
func CreateVpnServer(c *Client, server VpnServer) (*VpnServer, error) {
	if len(server.Name) == 0 {
		server.Name = generateVpnServerName()
	}
	if len(server.Network) == 0 {
		server.Network = "10.12.12.0/24"
	}
	if len(server.Protocol) == 0 {
		server.Protocol = "udp"
	}
	if len(server.Cipher) == 0 {
		server.Cipher = "aes128"
	}
	if len(server.Hash) == 0 {
		server.Hash = "sha1"
	}
	if len(server.NetworkMode) == 0 {
		server.NetworkMode = "tunnel"
	}

	opts := RequestOpts{
		JSONBody:     server,
		JSONResponse: &server,
	}
	if _, err := c.Request("post", getCreateServerPath(), &opts); err != nil {
		return nil, err
	}
	return &server, nil
}

// StartStopServer 启动或者停止vpn server
func StartStopServer(c *Client, serverId string, start bool) (*VpnServer, error) {
	var server VpnServer
	opts := RequestOpts{
		JSONResponse: &server,
	}
	operation := "start"
	if !start {
		operation = "stop"
	}
	if _, err := c.Request("put", getServerStartStopUrl(serverId, operation), &opts); err != nil {
		return nil, err
	}
	return &server, nil
}

// Organization 组织数据结构
type Organization struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	AuthApi    bool   `json:"auth_api"`
	AuthToken  string `json:"auth_token"`
	AuthSecret string `json:"auth_secret"`
	UserCount  int    `json:"user_count"`
}

// GetOrganizationList 获取组织列表
func GetOrganizationList(c *Client) ([]Organization, error) {
	var orgs []Organization
	opts := RequestOpts{
		JSONResponse: &orgs,
	}
	if _, err := c.Request("get", getOrganizationList(), &opts); err != nil {
		return nil, err
	}
	return orgs, nil
}

// AttachConf 添加组织的配置
type AttachConf struct {
	Id     string `json:"id"`             // 组织id
	Server string `json:"server"`         // server id
	Name   string `json:"name,omitempty"` // 组织名称
}

// AttachOrganizationToServer 为server添加一个组织，一个在pritunl中，一个vpn server必须要属于某个组织
func AttachOrganizationToServer(c *Client, conf AttachConf) (*AttachConf, error) {
	opts := RequestOpts{
		JSONBody:     conf,
		JSONResponse: &conf,
	}
	if _, err := c.Request("put", getAttachOrganizationUrl(conf.Server, conf.Id), &opts); err != nil {
		return nil, err
	}
	return &conf, nil
}

// RouteDetail 针对内部网段的路由配置信息
type RouteDetail struct {
	Id      string `json:"id,omitempty"` // 路由id, 添加路由时可为空
	Server  string `json:"server"`       // vpn server id
	Network string `json:"network"`      // 路由的网段
	Nat     bool   `json:"nat"`          // 针对此网段是否采用nat模式，否则就是路由模式
}

// RouteAddOpts 路由添加配置
type RouteAddOpts struct {
	Id      string `json:"id,omitempty"`
	Server  string `json:"server"`  // vpn server id
	Network string `json:"network"` // 路由的网段
	Nat     bool   `json:"nat"`     // 针对此网段是否采用nat模式，否则就是路由模式
}

// GetServerRouteList 获取指定vpn服务的路由列表
func GetServerRouteList(c *Client, serverId string) ([]RouteDetail, error) {
	var routes []RouteDetail
	opts := RequestOpts{
		JSONResponse: &routes,
	}
	if _, err := c.Request("get", getServerRoutesUrl(serverId), &opts); err != nil {
		return nil, err
	}
	return routes, nil
}

// DeleteRoute 删除指定路由
func DeleteRoute(c *Client, serverId, routeId string) error {
	if _, err := c.Request("delete", getDeleteRouteUrl(serverId, routeId), nil); err != nil {
		return err
	}
	return nil
}

// AddRoute 添加路由
func AddRoute(c *Client, route RouteAddOpts) (*RouteDetail, error) {
	var routeDetail RouteDetail
	opts := RequestOpts{
		JSONBody:     route,
		JSONResponse: &routeDetail,
	}
	if _, err := c.Request("post", getAddRouteUrl(route.Server), &opts); err != nil {
		return nil, err
	}
	return &routeDetail, nil
}

// RouteUpdateOpts 路由添加配置
type RouteUpdateOpts struct {
	Id      string `json:"id"`                // 路由id
	Server  string `json:"server"`            // vpn server id
	Network string `json:"network,omitempty"` // 路由的网段
	Nat     bool   `json:"nat,omitempty"`     // 针对此网段是否采用nat模式，否则就是路由模式
}

// UpdateRoute 更新路由配置
func UpdateRoute(c *Client, route RouteUpdateOpts) (*RouteDetail, error) {
	var routeDetail RouteDetail
	opts := RequestOpts{
		JSONBody:     route,
		JSONResponse: &routeDetail,
	}
	if _, err := c.Request("put", getUpdateRouteUrl(route.Server, route.Id), &opts); err != nil {
		return nil, err
	}
	return &routeDetail, nil
}

// UserDetail 用户详情
type UserDetail struct {
	Id               string `json:"id"`
	Organization     string `json:"organization"`
	OrganizationName string `json:"organization_name"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Disabled         bool   `json:"disabled"` // 是否被禁用
}

// UserAddOpts 用户添加配置
type UserAddOpts struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
}

// AddUser 向组织添加用户
func AddUser(c *Client, user UserAddOpts) ([]UserDetail, error) {
	var users []UserDetail
	opts := RequestOpts{
		JSONBody:     user,
		JSONResponse: &users,
	}
	if _, err := c.Request("post", getAddUserUrl(user.OrganizationId), &opts); err != nil {
		return nil, err
	}
	return users, nil
}

// UserUpdateOpts 用户更新选项
type UserUpdateOpts struct {
	UserId         string `json:"userId"`
	OrganizationId string `json:"organizationId"`
	Disabled       bool   `json:"disabled"` // 是否禁用
}

// EnableDisableUser 启用禁用用户
func EnableDisableUser(c *Client, conf UserUpdateOpts) (*UserDetail, error) {
	var userDetail UserDetail
	opts := RequestOpts{
		JSONBody:     map[string]bool{"disabled": conf.Disabled},
		JSONResponse: &userDetail,
	}
	if _, err := c.Request("put", getUpdateUserUrl(conf.OrganizationId, conf.UserId), &opts); err != nil {
		return nil, err
	}
	return &userDetail, nil
}

// ConnectionFile 连接配置文件
type ConnectionFile struct {
	Name    string `json:"name"`    // 连接文件名，以.ovpn为后缀，可直接在openvpn客户端导入
	Content []byte `json:"content"` // 连接文件的内容
}

// ExportUserConnectFile 导出用户的连接配置, 导出的是配置文件的tar包的内容，需要自己按需获取解压内容
func ExportUserConnectFile(c *Client, organizationId, userId string) (*ConnectionFile, error) {
	opts := RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := c.Request("get", getExportConnectFileUrl(organizationId, userId), &opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解压压缩包
	files, err := extractTar(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(files) != 1 {
		return nil, fmt.Errorf("export connection config file failed")
	}

	return &files[0], nil
}

// extractTar 解压tar文件内容
func extractTar(body io.Reader) ([]ConnectionFile, error) {
	files := []ConnectionFile{}

	// 创建一个tar读取器，并遍历tar文件的每个条目
	tr := tar.NewReader(body)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		// 检查是否是文件
		if header.Typeflag == tar.TypeReg {
			fileContent := &bytes.Buffer{}
			if _, err := io.Copy(fileContent, tr); err != nil {
				return nil, fmt.Errorf("failed to read file content: %w", err)
			}

			// 存储文件内容
			files = append(files, ConnectionFile{
				Name:    header.Name,
				Content: fileContent.Bytes(),
			})
		}
	}
	return files, nil
}
