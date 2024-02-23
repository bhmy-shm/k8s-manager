package maps

import (
	corev1 "k8s.io/api/core/v1"
	"sync"
)

type NodeMap struct {
	data sync.Map // [nodename string] *v1.Node   注意里面不是切片
}

func (n *NodeMap) Get(name string) *corev1.Node {
	if node, ok := n.data.Load(name); ok {
		return node.(*corev1.Node)
	}
	return nil
}

func (n *NodeMap) Add(item *corev1.Node) {
	//直接覆盖
	n.data.Store(item.Name, item)
}

func (n *NodeMap) Update(item *corev1.Node) bool {
	n.data.Store(item.Name, item)
	return true
}

func (n *NodeMap) Delete(node *corev1.Node) {
	n.data.Delete(node.Name)
}

func (n *NodeMap) ListAll() []*corev1.Node {
	ret := []*corev1.Node{}
	n.data.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(*corev1.Node))
		return true
	})
	return ret //返回空列表
}
