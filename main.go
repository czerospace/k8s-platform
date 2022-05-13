package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 k8s client
	service.K8s.Init() // 可以使用 service.K8s.clientset
	// 初始化 gin
	r := gin.Default()
	// 跨包调用 router 的初始化方法
	controller.Router.InitApiRouter(r)
	// 启动 gin server
	r.Run(config.ListenAddr)
}
