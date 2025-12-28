package node

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Port struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value,omitempty"`
}

type Node struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Params  []Port `json:"params,omitempty"`
	Inputs  []Port `json:"inputs,omitempty"`
	Outputs []Port `json:"outputs,omitempty"`
}

type NodeManager struct {
	Nodes    []Node
	filePath string
}

// NewNodeManager 会读取本地 node.json，如果不存在就创建空数据
func NewNodeManager() *NodeManager {
	// 获取当前程序工作目录
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "." // 获取失败就使用当前相对目录
	}

	filePath := filepath.Join(cwd, "node.json")

	manager := &NodeManager{
		Nodes:    []Node{},
		filePath: filePath,
	}

	manager.load()
	return manager
}

// AddNode 新增节点
func (m *NodeManager) AddNode(id, name string) {
	m.Nodes = append(m.Nodes, Node{
		ID:   id,
		Name: name,
	})
}

// RemoveNode 删除节点
func (m *NodeManager) RemoveNode(target *Node) {
	for i := range m.Nodes {
		if &m.Nodes[i] == target {
			m.Nodes = append(m.Nodes[:i], m.Nodes[i+1:]...)
			return
		}
	}
}

// Save 将节点数据写入 JSON 文件
func (m *NodeManager) Save() error {
	data, err := json.MarshalIndent(m.Nodes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

// load 从 JSON 文件读取节点数据
func (m *NodeManager) load() {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return // 文件不存在，使用空数据
	}

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return
	}

	_ = json.Unmarshal(data, &m.Nodes)
}
