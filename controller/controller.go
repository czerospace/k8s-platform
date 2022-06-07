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
	router.
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pods/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp)
}
