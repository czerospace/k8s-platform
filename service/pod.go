package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"k8s-platform/config"

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

// PodsNp 用于返回 namespace 中的 pod 数量
type PodsNp struct {
	Namespace string
	PodNum    int
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

// GetPodDetail  获取 Pod 详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取Pod详情失败: " + err.Error())
		return nil, errors.New("获取Pod详情失败: " + err.Error())
	}

	return pod, nil
}

// DeletePod 删除 Pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除Pod失败: " + err.Error())
		return errors.New("删除Pod失败: " + err.Error())
	}

	return nil
}

// UpdatePod 更新 Pod
func (p *pod) UpdatePod(podName, namespace, content string) (err error) {
	// content 参数是请求中传入的 Pod 对象的 json 数据
	var pod = &corev1.Pod{}
	// 反序列化 Pod 对象
	err = json.Unmarshal([]byte(content), pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败: " + err.Error()))
		return errors.New("反序列化失败: " + err.Error())
	}
	// 更新 Pod
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新Pod失败: " + err.Error())
		return errors.New("更新Pod失败: " + err.Error())
	}

	return nil
}

// GetPodContainer 获取 Pod 的容器名列表
func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	// 获取 Pod 详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	// 从 Pod 对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}

// GetPodLog 获取 Pod 内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	// 设置日志的配置，容器名，获取的内容的配置
	LineLimit := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &LineLimit,
	}

	// 获取一个 request 实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	// 发起 Stream 链接，得到 Response.body
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error("获取podLog失败: " + err.Error())
		return "", errors.New("获取podLog失败: " + err.Error())
	}
	defer podLogs.Close()
	// 将 response body 写入到缓冲区，目的是为了转换成 string 类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error("复制podLog失败: " + err.Error())
		return "", errors.New("复制podLog失败: " + err.Error())
	}
	return buf.String(), nil
}

// GetPodNumPerNp 获取每个 namespace 里的 Pod 数量
func (p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	// 获取 namespace 列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		// 获取 pod 列表
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		// 组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum:    len(podList.Items),
		}
		// 添加到 podsNps 数组中
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil

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
