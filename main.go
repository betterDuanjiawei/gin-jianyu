package main

import (
	"context"
	"fmt"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/setting"
	"github.com/betterDuanjiawei/gin-jianyu/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// http.Server - Shutdown() 热更新版本
func main() {
	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("pid is: ", syscall.Getpid())
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("listen failed, err:%v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
