package main

import (
	"fmt"
	"github.com/betterDuanjiawei/gin-jianyu/models"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/logging"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/setting"
	"github.com/betterDuanjiawei/gin-jianyu/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

// endless 热更新版本
func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}
