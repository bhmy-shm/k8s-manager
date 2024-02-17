package model

type Deployment struct {
	Name       string   `json:"name"`
	NameSpace  string   `json:"nameSpace"`
	Replicas   [3]int32 `json:"replicas"`   //3个值，分别是总副本数，可用副本数 ，不可用副本数
	Images     string   `json:"images"`     //pod镜像名称
	IsComplete bool     `json:"isComplete"` //pod是否完成
	Message    string   `json:"message"`    //pod显示错误信息
	CreateTime string   `json:"createTime"` //pod创建时间
	Pods       []*Pod   `json:"pods"`
}
