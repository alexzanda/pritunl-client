package pritunl

import "fmt"

// getServerSettingsPath 服务配置获取url
func getServerSettingsPath() string {
	return "/settings"
}

// getAdminListPath 获取管理员列表url
func getAdminListPath() string {
	return "/admin"
}

// getUpdateAdminUserPath 获取更新管理员配置的url
func getUpdateAdminUserPath(userId string) string {
	return fmt.Sprintf("/admin/%s", userId)
}

// getCreateServerPath 获取创建vpn server的url
func getCreateServerPath() string {
	return "/server"
}

// getOrganizationList 获取组织列表url
func getOrganizationList() string {
	return "/organization"
}

// getAttachOrganizationUrl 获取添加组织url
func getAttachOrganizationUrl(serverId, organizationId string) string {
	return fmt.Sprintf("/server/%s/organization/%s", serverId, organizationId)
}

// getServerRoutesUrl 获取vpn server路由列表url
func getServerRoutesUrl(serverId string) string {
	return fmt.Sprintf("/server/%s/route", serverId)
}

// getDeleteRouteUrl 获取删除路由的url
func getDeleteRouteUrl(serverId, routeId string) string {
	return fmt.Sprintf("/server/%s/route/%s", serverId, routeId)
}

// getAddRouteUrl 获取添加路由的url
func getAddRouteUrl(serverId string) string {
	return fmt.Sprintf("/server/%s/route", serverId)
}

// getUpdateRouteUrl 获取更新路由的url
func getUpdateRouteUrl(serverId, routeId string) string {
	return fmt.Sprintf("/server/%s/route/%s", serverId, routeId)
}

// getAddUserUrl 获取添加用户的url
func getAddUserUrl(organizationId string) string {
	return fmt.Sprintf("/user/%s", organizationId)
}

// getUpdateUserUrl 获取更新用户的url
func getUpdateUserUrl(organizationId, userId string) string {
	return fmt.Sprintf("/user/%s/%s", organizationId, userId)
}

// getServerStartStopUrl 获取vpn server启动停止url
func getServerStartStopUrl(serverId, operation string) string {
	return fmt.Sprintf("/server/%s/operation/%s", serverId, operation)
}

// getExportConnectFileUrl 获取导出用户连接配置文件的url
func getExportConnectFileUrl(organizationId, userId string) string {
	return fmt.Sprintf("/key/%s/%s.tar", organizationId, userId)
}
