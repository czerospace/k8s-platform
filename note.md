```shell
#知识点总结
   匿名结构体，用于声明入参
   get 请求为 form 格式，其他请求为 json 格式
	eg：用于 get 请求
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	使用ctx.Bind(params)绑定
	eg: 用于 其他 请求
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})
	使用ctx.ShouldBindJSON(params)绑定
```