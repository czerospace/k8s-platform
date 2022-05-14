package controller

import (
	"github.com/gin-gonic/gin"
)

// Router 初始化 router 类型的对象，首字母大写，用于跨包调用
var Router router

// router 生命一个 router 的结构体
type router struct{}

// InitApiRouter 初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/api/k8s/pods", Pod.GetPods)
}
