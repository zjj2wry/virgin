package config

import (
	"log"

)

var Conf *Config


func init() {
	Conf = &Config{}
	Conf.InitConfig("conf/app.conf")
	log.Println(Conf)

	log.Println("orderrpc = ", Conf.GetString("orderrpc"))
}
