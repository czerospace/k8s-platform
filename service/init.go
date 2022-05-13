package service

import (
	"k8s-platform/config"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

// Init 初始化
func (k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		panic("获取k8s client 配置失败: " + err.Error())
	}
	// 根据 rest.config 类型的对象， new 一个 clientset 出来
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建k8s client失败： " + err.Error())
	} else {
		logger.Info("k8s client 初始化成功!")
	}

	k.ClientSet = clientset
}
