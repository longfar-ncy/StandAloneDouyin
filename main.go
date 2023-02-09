// Code generated by hertz generator.

package main

import (
	"douyin/dal/db"
	"douyin/util"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Init() {
	db.Init()
}

func main() {
	h := server.Default()

	Init()
	util.ScheduledInit()

	register(h)
	h.Spin()
}
