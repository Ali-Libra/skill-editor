package skill

import (
	"fyne.io/fyne/v2"
)

// PortLocation 表示一个端口的位置
type PortLocation struct {
	SkillNodeID string // 所属 SkillNode 的 ID
	PortName    string // 端口名称
	PortType    string // "input" 或 "output"
	Position    fyne.Position
}

// ConnectionManager 管理所有的连接
type ConnectionManager struct {
	editor          *SkillEditor
	dragSourcePort  *PortLocation   // 正在拖动的源端口
	availableTarget map[string]bool // 可用的目标端口
	isDrawingLine   bool
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(editor *SkillEditor) *ConnectionManager {
	return &ConnectionManager{
		editor:          editor,
		availableTarget: make(map[string]bool),
	}
}

// StartDragFromPort 开始从端口拖动
func (cm *ConnectionManager) StartDragFromPort(port *PortLocation) {
	cm.dragSourcePort = port
	cm.isDrawingLine = true
	cm.updateAvailableTargets()
}

// updateAvailableTargets 更新可用的目标端口
func (cm *ConnectionManager) updateAvailableTargets() {
	cm.availableTarget = make(map[string]bool)

	if cm.dragSourcePort == nil || cm.editor.skill == nil {
		return
	}

	// 找到源节点的索引
	var sourceNodeIndex int = -1
	for i := range cm.editor.skill.Nodes {
		if cm.editor.skill.Nodes[i].ID == cm.dragSourcePort.SkillNodeID {
			sourceNodeIndex = i
			break
		}
	}

	if sourceNodeIndex == -1 {
		return
	}

	// 只能连接到相邻节点
	if cm.dragSourcePort.PortType == "output" {
		// 输出端口只能连接到下一个节点的同名输入端口
		if sourceNodeIndex+1 < len(cm.editor.skill.Nodes) {
			nextNode := &cm.editor.skill.Nodes[sourceNodeIndex+1]
			for _, input := range nextNode.NodeData.Inputs {
				if input.Name == cm.dragSourcePort.PortName {
					portID := nextNode.ID + ":input:" + input.Name
					cm.availableTarget[portID] = true
				}
			}
		}
	} else if cm.dragSourcePort.PortType == "input" {
		// 输入端口只能连接到前一个节点的同名输出端口
		if sourceNodeIndex > 0 {
			prevNode := &cm.editor.skill.Nodes[sourceNodeIndex-1]
			for _, output := range prevNode.NodeData.Outputs {
				if output.Name == cm.dragSourcePort.PortName {
					portID := prevNode.ID + ":output:" + output.Name
					cm.availableTarget[portID] = true
				}
			}
		}
	}
}

// ConnectPorts 连接两个端口
func (cm *ConnectionManager) ConnectPorts(targetPort *PortLocation) bool {
	if cm.dragSourcePort == nil {
		return false
	}

	// 检查是否是有效的连接
	if cm.dragSourcePort.PortType == targetPort.PortType {
		return false // 同类型端口不能连接
	}

	// 确定源和目标
	var sourceNode, targetNode *SkillNode
	var sourcePortName, targetPortName string

	if cm.dragSourcePort.PortType == "output" {
		sourceNode = cm.getSkillNodeByID(cm.dragSourcePort.SkillNodeID)
		targetNode = cm.getSkillNodeByID(targetPort.SkillNodeID)
		sourcePortName = cm.dragSourcePort.PortName
		targetPortName = targetPort.PortName
	} else {
		sourceNode = cm.getSkillNodeByID(targetPort.SkillNodeID)
		targetNode = cm.getSkillNodeByID(cm.dragSourcePort.SkillNodeID)
		sourcePortName = targetPort.PortName
		targetPortName = cm.dragSourcePort.PortName
	}

	if sourceNode == nil || targetNode == nil {
		return false
	}

	// 保存连接关系
	if sourceNode.Outputs[sourcePortName] == nil {
		sourceNode.Outputs[sourcePortName] = []string{}
	}
	sourceNode.Outputs[sourcePortName] = append(
		sourceNode.Outputs[sourcePortName],
		targetNode.ID,
	)

	targetNode.Inputs[targetPortName] = sourceNode.ID

	cm.EndDrag()
	return true
}

// EndDrag 结束拖动
func (cm *ConnectionManager) EndDrag() {
	cm.dragSourcePort = nil
	cm.isDrawingLine = false
	cm.availableTarget = make(map[string]bool)
}

// getSkillNodeByID 通过 ID 获取 SkillNode
func (cm *ConnectionManager) getSkillNodeByID(id string) *SkillNode {
	if cm.editor.skill == nil {
		return nil
	}
	for i := range cm.editor.skill.Nodes {
		if cm.editor.skill.Nodes[i].ID == id {
			return &cm.editor.skill.Nodes[i]
		}
	}
	return nil
}

// IsPortAvailable 检查端口是否可用
func (cm *ConnectionManager) IsPortAvailable(portID string) bool {
	return cm.availableTarget[portID]
}
