package controller

import (
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Pod pod

type pod struct{}

// Controller 中的方法入参是 gin.Context ，用于从上下文中获取请求参数及定义响应内容
// 流程：绑定参数->调用 service 代码->根据调用结果响应具体内容

// GetPods 获取 Pod 列表，支持分页，过滤，排序
func (p *pod) GetPods(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于声明入参，get 请求为 form 格式，其他请求为 json 格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// 绑定参数，给匿名结构体中的属性赋值，值是入参
	// form格式使用 ctx.Bind 方法，json 格式使用 ctx.ShouldBindJSON 方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind绑定参数失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

// GetPodDetail 获取 Pod 详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod详情成功",
		"data": data,
	})
}

// DeletePod 删除 Pod
func (p *pod) DeletePod(ctx *gin.Context) {
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	// DELETE 请求，绑定参数方法改为 ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Pod成功",
		"data": nil,
	})
}

// UpdatePod 更新 Pod
func (p *pod) UpdatePod(ctx *gin.Context) {
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	// PUT 请求，绑定参数方法改为 ctx.ShouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Pod.UpdatePod(params.PodName, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Pod成功",
		"data": nil,
	})
}

// GetPodContainer 获取 Pod 容器
func (p *pod) GetPodContainer(ctx *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	// Get 请求，绑定参数方法为 ctx.Bind
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod容器成功",
		"data": data,
	})
}

// GetPodLog 获取 Pod 中的日志
func (p *pod) GetPodLog(ctx *gin.Context) {
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod中容器日志成功",
		"data": data,
	})
}

// GetPodNumPerNp 获取每个 namespace 的 pod 数量
func (p *pod) GetPodNumPerNp(ctx *gin.Context) {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个namespace的pod数量成功",
		"data": data,
	})
}
