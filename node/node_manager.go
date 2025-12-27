package node

type Node struct {
	ID   int
	Name string
}

type NodeManager struct {
	Nodes []Node
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		Nodes: []Node{
			{ID: 1, Name: "伤害节点"},
			{ID: 2, Name: "冷却节点"},
		},
	}
}
