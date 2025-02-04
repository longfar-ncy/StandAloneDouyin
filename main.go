// Code generated by hertz generator.

package main

import (
	"flag"

	"douyin/biz/handler/api"
	"douyin/biz/mw"
	"douyin/pkg/global"
	"douyin/pkg/initialize"
)

func Init() {
	flag.StringVar(&global.ConfigPath, "c", "./pkg/config/config.yml", "config file path")
	flag.Parse()

	initialize.Viper()
	initialize.MySQL()
	initialize.Redis()
	initialize.Global()
	mw.InitJWT()

	// initialize.Hertz() 需要保持在最下方，因为调用完后 Hertz 就启动完毕了
	go api.MannaClient.Run()
	initialize.Hertz()
}

func main() {
	Init()
}
