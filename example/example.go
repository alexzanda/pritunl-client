package main

import (
	"fmt"
	"log"

	"nve/pkg/pritunl"
)

var apiToken = "lKULJmhyKNvzggiCFcoGBhgpWzgAkLx7"
var apiSecret = "c94vfEDDmfWvhLuqlcVDEhRCdeoNvfYb"
var host = "192.170.1.193"

func main() {
	//client, err := pritunl.NewClient(apiToken, apiSecret, host, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//// 更新服务端的public address
	//fmt.Println(pritunl.UpdatePublicAccessAddress(client, "192.170.1.111"))

	//// 获取管理账号列表
	//adminUsers, err := pritunl.GetAdminUserList(client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, adminUser := range adminUsers {
	//	if adminUser.Username == "pritunl" {
	//		fmt.Printf("---pritunl user id: %s\n", adminUser.Id)
	//	}
	//}

	////// 更新admin用户的认证配置
	//adminUser := pritunl.AdminUser{
	//	Id:        "6749534c2e05c6b19c0435a7",
	//	Username:  "pritunl",
	//	AuthApi:   true,
	//	SuperUser: true,
	//}
	//userConf, err := pritunl.UpdateAdminUserAuthConfig(client, adminUser)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("user token: %s, secret: %s\n", userConf.Token, userConf.Secret)

	//// 创建一个新的server，如果vpn私有网段不指定，就自动生成
	//server := pritunl.VpnServer{}
	//s, err := pritunl.CreateVpnServer(client, server)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("name: %s, server id: %s network: %s, port: %d\n", s.Name, s.Id, s.Network, s.Port)

	//// 获取组织列表，内置的pritunl镜像，默认会内置一个组织，名为default
	//orgs, err := pritunl.GetOrganizationList(client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, org := range orgs {
	//	fmt.Printf("org id: %s, name: %s\n", org.Id, org.Name)
	//}

	//// 为server添加一个组织
	//attachConf := pritunl.AttachConf{
	//	Id:     "674953b02e05c6b19c0435f7",
	//	Server: "674ee8e80d1fc18bf2c6339f",
	//}
	//c, err := pritunl.AttachOrganizationToServer(client, attachConf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(c)

	//// 获取指定server的路由列表
	//rs, err := pritunl.GetServerRouteList(client, "674e68150d1fc18bf2c5ce4f")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, route := range rs {
	//	fmt.Printf("route id: %s, network: %s\n", route.Id, route.Network)
	//}

	//// 删除默认的0.0.0.0/0的路由
	//if err := pritunl.DeleteRoute(client, "674e68150d1fc18bf2c5ce4f", "302e302e302e302f30"); err != nil {
	//	log.Fatal(err)
	//}

	//// 添加一个要放行的内网网段的路由
	//route := pritunl.RouteAddOpts{
	//	Server:  "674e68150d1fc18bf2c5ce4f",
	//	Network: "10.11.0.0/16",
	//	Nat:     false,
	//}
	//r, err := pritunl.AddRoute(client, route)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("route id: %s, network: %s, nat: %v\n", r.Id, r.Network, r.Nat)

	//// 更新一个路由的模式
	//opts := pritunl.RouteUpdateOpts{
	//	Id:     "31302e31312e302e302f3136",
	//	Server: "674e68150d1fc18bf2c5ce4f",
	//	Nat:    true,
	//}
	//r, err := pritunl.UpdateRoute(client, opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("route id: %s, is nat: %v\n", r.Id, r.Nat)

	//// 启动停止一个vpn server
	//s, err := pritunl.StartStopServer(client, "674e68150d1fc18bf2c5ce4f", false)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("server status: %s\n", s.Status)

	//// 新建一个用户
	//opts := pritunl.UserAddOpts{
	//	Name:           "aaaa",
	//	OrganizationId: "674953b02e05c6b19c0435f7",
	//}
	//users, err := pritunl.AddUser(client, opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, user := range users {
	//	fmt.Printf("user id: %s, name: %s, disabled: %v\n", user.Id, user.Name, user.Disabled)
	//}

	//// 启用禁用用户
	//opts := pritunl.UserUpdateOpts{
	//	UserId:         "674953b52e05c6b19c04361b",
	//	OrganizationId: "674953b02e05c6b19c0435f7",
	//	Disabled:       true,
	//}
	//u, err := pritunl.EnableDisableUser(client, opts)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("user disabled: %v\n", u.Disabled)

	//// 导出对应选手的连接配置信息
	//connFile, err := pritunl.ExportUserConnectFile(client, "674953b02e05c6b19c0435f7", "674953b52e05c6b19c043613")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("connection file, name: %s, content: %s\n", connFile.Name, connFile.Content)

	// 一键启动一个pritunl
	totalConf, err := pritunl.InitVpnServer(host, "192.170.1.163", "10.10.0.0/16", false, apiToken, apiSecret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(totalConf)
}
