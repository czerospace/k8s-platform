package service

import (
	corev1 "k8s.io/api/core/v1"
)

var Pod pod

type pod struct{}

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
