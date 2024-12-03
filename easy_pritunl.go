package pritunl

import (
	"fmt"
	"log"
	"net"
)

// 对pritunl各个基础方法的进一步封装，实现一键配置并启动pritunl的vpn环境，并返回一个配置实例

const (
	DEFAULT_ADMIN_USER   = "pritunl"
	DEFAULT_ORGANIZATION = "default"
	DEFAULT_ROUTE        = "0.0.0.0/0"
)

// PritunlTotalConfig 一个正常运行着的完整vpn服务所用到的配置
type PritunlTotalConfig struct {
	PublicAddress  string `json:"publicAddress"`  // 服务对外地址
	AdminUserId    string `json:"adminUserId"`    // 默认的管理员账号id
	ApiToken       string `json:"apiToken"`       // api token
	ApiSecret      string `json:"apiSecret"`      // api secret
	VpnServerName  string `json:"vpnServerName"`  // vpn server名称
	VpnServerId    string `json:"vpnServerId"`    // vpn server id
	VpnServerState string `json:"vpnServerState"` // vpn server的状态，online offline
	VpnNetwork     string `json:"vpnNetwork"`     // vpn网段
	VpnPort        int    `json:"vpnPort"`        // vpn端口
	OrganizationId string `json:"organizationId"` // 组织id
	Route          string `json:"route"`          // vpn连接的内部网络路由
	RouteId        string `json:"routeId"`        // vpn连接的内部网络的路由id
	RouteUseNat    bool   `json:"routeUseNat"`    // vpn连接的内部网络是否启用nat模式

}

// InitVpnServer 一键初始化一个vpn服务。包括修改默认的认证key，创建vpn server、配置组织、路由、启动服务等
func InitVpnServer(adminIp, publicAddr, network string, useNat bool, apiToken, apiSecret string) (*PritunlTotalConfig, error) {
	// 校验外网地址
	if net.ParseIP(publicAddr) == nil {
		return nil, fmt.Errorf("public addr is invalid")
	}
	if _, _, err := net.ParseCIDR(network); err != nil {
		return nil, fmt.Errorf("network is invalid")
	}

	client, err := NewClient(apiToken, apiSecret, adminIp, nil)
	if err != nil {
		return nil, fmt.Errorf("create new client failed, err: %w", err)
	}

	totalConf := PritunlTotalConfig{}

	// 获取管理账号列表
	adminUsers, err := GetAdminUserList(client)
	if err != nil {
		return nil, fmt.Errorf("get admin user list failed, err: %w", err)
	}
	var defaultAdmin AdminUser
	for _, adminUser := range adminUsers {
		if adminUser.Username == DEFAULT_ADMIN_USER {
			defaultAdmin = adminUser
			break
		}
	}
	totalConf.AdminUserId = defaultAdmin.Id

	// 更新admin用户的认证配置
	updateAdminOpts := AdminUser{
		Id:        defaultAdmin.Id,
		Username:  defaultAdmin.Username,
		AuthApi:   true,
		SuperUser: true,
	}
	userConf, err := UpdateAdminUserAuthConfig(client, updateAdminOpts)
	if err != nil {
		return nil, fmt.Errorf("update admin user config failed, err: %w", err)
	}
	totalConf.ApiToken = userConf.Token
	totalConf.ApiSecret = userConf.Secret

	// 创建新的client
	client, err = NewClient(totalConf.ApiToken, totalConf.ApiSecret, adminIp, nil)
	if err != nil {
		return nil, fmt.Errorf("create new client failed, err: %w", err)
	}

	// 创建一个新的server，如果vpn私有网段不指定，就自动生成
	server := VpnServer{}
	s, err := CreateVpnServer(client, server)
	if err != nil {
		return nil, fmt.Errorf("create vpn server failed, err: %w", err)
	}
	totalConf.VpnServerName = s.Name
	totalConf.VpnServerId = s.Id
	totalConf.VpnNetwork = s.Network
	totalConf.VpnPort = s.Port

	// 获取组织列表，内置的pritunl镜像，默认会内置一个组织，名为default
	orgs, err := GetOrganizationList(client)
	if err != nil {
		return nil, fmt.Errorf("get organizations failed, err: %w", err)
	}
	var defaultOrg Organization
	for _, org := range orgs {
		if org.Name == DEFAULT_ORGANIZATION {
			defaultOrg = org
			break
		}
	}

	// 为server指定一个组织
	attachConf := AttachConf{
		Id:     defaultOrg.Id,
		Server: s.Id,
	}
	if _, err = AttachOrganizationToServer(client, attachConf); err != nil {
		return nil, fmt.Errorf("attach organization to server failed, err: %w", err)
	}
	totalConf.OrganizationId = defaultOrg.Id

	// 获取server的路由列表
	rs, err := GetServerRouteList(client, s.Id)
	if err != nil {
		log.Fatal(err)
	}
	var defaultRoute RouteDetail
	for _, route := range rs {
		if route.Network == DEFAULT_ROUTE {
			defaultRoute = route
			break
		}
	}

	// 删除默认的路由
	if err = DeleteRoute(client, s.Id, defaultRoute.Id); err != nil {
		return nil, fmt.Errorf("delete default route failed, err: %w", err)
	}

	// 添加内网网段路由
	route := RouteAddOpts{
		Server:  s.Id,
		Network: network,
		Nat:     useNat,
	}
	r, err := AddRoute(client, route)
	if err != nil {
		return nil, fmt.Errorf("add internal route failed, err: %w", err)
	}
	totalConf.Route = network
	totalConf.RouteId = r.Id
	totalConf.RouteUseNat = useNat

	// 启动vpn server
	s, err = StartStopServer(client, s.Id, true)
	if err != nil {
		return nil, fmt.Errorf("start vpn server failed, err: %w", err.Error())
	}
	totalConf.VpnServerState = s.Status

	// 更新服务端的public address, 这一步会导致pritunl服务端重启，需要放在最后一步
	if _, err = UpdatePublicAccessAddress(client, publicAddr); err != nil {
		return nil, fmt.Errorf("update public addr failed, err: %w", err)
	}
	totalConf.PublicAddress = publicAddr

	return &totalConf, nil
}
