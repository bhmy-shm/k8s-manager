package model

/*
	pod 状态
	pod 的 podConditions 取所有的状态均为 true
	PodScheduled： Pod已经被调度到某节点；
	ContainersReady：Pod中所有容器都已就绪；
	Initialized：所有的init容器都已经启动成功；
	Ready：Pod可以为请求提供服务。
*/

type Pod struct {
	Name       string
	NameSpace  string //新增一个命名空间
	Images     string
	NodeName   string
	IP         []string //第一个是 POD IP 第二个是 node ip
	Phase      string   //pod 当前所处的阶段(running)
	IsReady    bool     //判断pod 是否就绪
	Message    string
	CreateTime string
}
