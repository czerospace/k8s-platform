package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 读取配置文件
	// 将 kubeconfig 文件转换成 rest.config 类型的对象
	conf, err := clientcmd.BuildConfigFromFlags("", "E:\\study\\config\\opsuser.kubeconfig")
	if err != nil {
		panic(err)
	}
	// 根据 rest.config 类型的对象，new 一个 clientset 出来
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err)
	}
	// 使用 clienset 获取 pod 列表
	podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, pod := range podList.Items {
		fmt.Printf("POD 名称:%v  命名空间:%v\n", pod.Name, pod.Namespace)
	}
}
