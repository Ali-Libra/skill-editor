package skill

import "skill-editor/node"

// SkillNode 是技能编辑器中的节点实例
type SkillNode struct {
	ID       string // 唯一实例 ID
	NodeID   string // 对应 Node.ID
	NodeData node.Node

	Inputs  map[string]string   // inputPortID -> from SkillNodeID
	Outputs map[string][]string // outputPortID -> to SkillNodeIDs
}
