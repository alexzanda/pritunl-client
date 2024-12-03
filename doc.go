package pritunl

// pritunl客户端实现，pritunl是一款开源的基于openvpn协议的vpn管理平台，提供openvpn server管理、用户管理等
// 使用pritunl时，初始配置步骤主要包含如下步骤：
// 1、更改默认的public access address，这个地址用在创建用户连接配置文件时的回连地址，新创建的pritunl服务端需要修改这个地址;
// 2、更新内置管理员账号的api token和api key，镜像启动时内置的默认管理账号已经配置了api token和api key，为了安全起见，请修改他;
//     2.1 获取管理员账户列表，找到pritunl账号，记录账号id；
//     2.2 使用账号id和默认的api token及api key，更新settings，服务端会返回最新的settings，请记住新的api key及api token；
// 3、创建一个server，即一个vpn服务端，可在创建时指定vpn私有网络的地址段；
// 4、获取组织列表
// 5、为server连接一个组织
// 6、获取server的默认路由表；
// 7、删除server默认路由表中的0.0.0.0这个路由，这个路由在客户端连上来时，会在vpn客户端默认添加一个0.0.0.0到vpn的路由，实际是不需要的，只需要
//   增加我们关心的内部网段即可
// 8、添加内部网段的路由，比如10.10.0.0/16网段，这个网段是pritunl server连接的网段，添加路由的时候可选择指定采用NAT模式还是ROUTE模式，对于
//   NAT模式，vpn客户端访问内部网络的流量，对于内部网络来说，看到的都是vpn server在内网的地址，如果ROUTE模式，则是vpn客户端的地址
// 9、当有客户端需要连接的时候，为每一个选手创建一个新的用户，便于审计
// 10、业务层面选择导出vpn连接配置
// 11、用户选择下载vpn客户端，导入连接配置
