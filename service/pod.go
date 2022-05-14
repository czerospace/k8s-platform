package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

// PodsResp 定义列表的返回内容，Items 是 pod 元素列表， total 是元素数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items"`
}

// GetPods 获取 pod 列表，支持过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	// context.TODO() 用于声明一个空的 context 上下文，用于 List 方法内置这个请求的超时
	// metav1.ListOptions{}) 用于过滤 List 数据，如使用label，field 等
	// 比如 kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// logger 用于打印日志
		// return 用于返回 response 内容
		logger.Info("获取Pod列表失败: " + err.Error())
		return nil, errors.New(err.Error())
	}
	// 实例化 dataSelector 结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelect: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 再排序和分页
	data := filtered.Sort().Paginate()
	//将 DataCell 类型转成 Pod
	pods := p.fromCells(data.GenericDataList)

	/*
		// 处理后的数据和原始数据的比较
		// 处理后的数据
		fmt.Println("============处理后的数据=========")
		for _, pod := range pods {
			fmt.Println(pod.Name, pod.CreationTimestamp)
		}
		// 原始数据
		fmt.Println("============原始数据=========")
		for _, pod := range podList.Items {
			fmt.Println(pod.Name, pod.CreationTimestamp)
		}
	*/

	// 将 []DataCell 类型的 pod 列表转为 v1.pod 列表
	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil
}

// 定义 DataCell 与 Pod 类型之间转换的方法 corev1.pod -> DataCell , DataCell -> corev1.pod
// toCells 方法用于将 pod 类型数组，转换成 DataCell 类型数组
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// fromCells 方法用于将 DataCell 类型数组，转换成 pod 类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		// cells[i].(podCell) 用断言将 DataCell 类型转成 podCell 类型
		// 最后转成换 Pod 类型，放入 pods 中
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
