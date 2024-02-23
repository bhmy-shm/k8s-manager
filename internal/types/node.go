package types

type NodeModel struct {
	Name       string
	IP         string
	HostName   string
	Labels     []string      //标签
	Capacity   *NodeCapacity //容量 包含了cpu 内存和pods数量
	Usage      *NodeUsage    //资源 使用情况
	CreateTime string
}

// NodeUsage pod使用数量
type NodeUsage struct {
	Pods   int
	Cpu    float64
	Memory float64
}

func NewNodeUsage(pods int, cpu float64, memory float64) *NodeUsage {
	return &NodeUsage{Pods: pods, Cpu: cpu, Memory: memory}
}

// NodeCapacity 节点容量
type NodeCapacity struct {
	Cpu    int64
	Memory int64
	Pods   int64
}

func NewNodeCapacity(cpu int64, memory int64, pods int64) *NodeCapacity {
	return &NodeCapacity{Cpu: cpu, Memory: memory, Pods: pods}
}
