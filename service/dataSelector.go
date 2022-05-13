package service

import (
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// 数组的排序、过滤、分页

// 一、定义数据结构
// dataSelector 用于封装排序、过滤、分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell 接口，用于各种资源 List 的类型转换，转换后可以使用 dataSelector 的排序、过滤、分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体，过滤:Name ,分页: Limit和Page
// Limit 是单页的数据条数
// Page 是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

// 二、排序
// 实现自定义结构的排序，需要重写Len、Swap、Less方法

// Len 方法用于获取数组长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap 方法用于数组中的元素在比较大小后的位置交换，可定义升序或降序
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less 方法用于定义数组中元素排序的“大小”的比较方式,按创建时间进行比较
func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

// Sort 重写以上3个方法用使用sort.Sort进行自定义排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 三、过滤

// Filter 方法用于过滤元素，比较元素的Name属性，若包含，再返回
func (d *dataSelector) Filter() *dataSelector {
	// 若 Name 的传参为空，则返回所有元素
	if "" == d.DataSelect.FilterQuery.Name {
		return d
	}
	// 若 Name 的传参不为空，则按照入参 Name 进行过滤
	// 声明一个新的数组，若 Name 包含，则把数据放进数组返回出去
	var filteredList []DataCell
	for _, value := range d.GenericDataList {
		// 定义一个匹配标签看 Name 是否在 d.GenericDataList 中
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.FilterQuery.Name) {
			// 如果不包含 objName,则 matches 为 false
			matches = false
			continue
		}
		if matches {
			filteredList = append(filteredList, value)
		}
	}
	d.GenericDataList = filteredList
	return d
}

// 四、分页

// Paginate 方法用于数组分页，根据Limit和Page的传参，取一定范围内的数据,返回数据
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.DataSelect.PaginateQuery.Limit
	page := d.DataSelect.PaginateQuery.Page
	// 检验参数的合法性
	if limit <= 0 || page <= 0 {
		return d
	}
	// 定义分页数据范围需要的 startIndex 和 endIndex
	// 举例，有25个元素的数组，limit 是10
	// page 是 1，startIndex 是 0，endIndex 是 9
	// page 是 2，startIndex 是 10，endIndex 是 19
	// page 是 3，startIndex 是 20，endIndex 是 25
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1
	// 处理最后一页 endIndex
	if endIndex > len(d.GenericDataList) {
		endIndex = len(d.GenericDataList) - 1
	}
	d.GenericDataList = d.GenericDataList[startIndex : endIndex+1]
	return d
}

// 定义 podCell 类型，重写 GetCreation 和 GetName 方法后，可以进行数据转换

type podCell corev1.Pod

// corev1.Pod -> podCell -> DataCell
// appsv1.Deployment -> deployCell -> DataCell

// GetCreation 重写 DataCell 接口的两个方法
func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
