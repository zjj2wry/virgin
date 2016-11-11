package config

import (
	"config.1kf.com/server/base_go/log"
	"config.1kf.com/server/base_go/rpc/order"
	"config.1kf.com/server/base_go/rpc/user"
)

/***
定时任务控制器
*/

var Conf *Config

var (
	// 订单RPC
	OrderRPC *order.OrderRPCClient
	//用户RPC
	UserPRC *user.UserRPCClient
)

func init() {
	Conf = &Config{}
	Conf.InitConfig("conf/app.conf")
	log.Info(Conf)

	log.Info("orderrpc = ", Conf.GetString("orderrpc"))
	OrderRPC = order.NewOrderRPCClient(Conf.GetString("orderrpc"))

	UserPRCURL := Conf.GetString("UserPRCURL")
	log.Info("userrpc:", UserPRCURL)
	UserPRC = user.NewUserRPCClient(UserPRCURL)
}
