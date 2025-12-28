package skill

import (
	"skill-editor/node"
	"skill-editor/tool"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SkillEditor 负责技能编辑画布
type SkillEditor struct {
	skill       *Skill
	canvas      *fyne.Container
	viewport    *Viewport
	nodeManager *node.NodeManager
	scale       float32
}

// NewSkillEditor 创建编辑器
func NewSkillEditor(nodeMgr *node.NodeManager) *SkillEditor {
	canvas := container.NewWithoutLayout()
	viewport := NewViewport(canvas)

	editor := &SkillEditor{
		canvas:      canvas,
		viewport:    viewport,
		nodeManager: nodeMgr,
		scale:       1.0,
	}
	return editor
}

func (e *SkillEditor) Root() fyne.CanvasObject {
	return e.viewport
}

func (e *SkillEditor) Refresh() {
	e.canvas.Objects = nil
	if e.skill == nil {
		return
	}

	// 主节点
	main := widget.NewCard(e.skill.Name, "", widget.NewLabel("技能主节点"))
	main.Resize(fyne.NewSize(200, 120))
	main.Move(fyne.NewPos(100, 50))
	e.canvas.Add(main)

	// 技能节点
	for i := range e.skill.Nodes {
		n := &e.skill.Nodes[i]
		w := NewNodeWidget(n, e)
		w.Resize(fyne.NewSize(180, 100))
		w.Move(fyne.NewPos(n.X, n.Y))
		e.canvas.Add(w)
	}

	e.canvas.Refresh()
}

func (e *SkillEditor) SetSkill(skill *Skill) {
	e.skill = skill
	e.Refresh()
}

// 添加节点到画布（鼠标点击位置）
func (e *SkillEditor) AddNode(n *node.Node, pos fyne.Position) {
	sn := &SkillNode{
		ID:       tool.GenerateUniqueID(),
		NodeID:   n.ID,
		NodeData: *n,
		X:        pos.X,
		Y:        pos.Y,
		Inputs:   map[string]string{},
		Outputs:  map[string][]string{},
	}
	e.skill.Nodes = append(e.skill.Nodes, *sn)
	e.Refresh()
}

// ---------------- 右键菜单 ----------------
func (e *SkillEditor) ShowNodeRightClickMenu(pos fyne.Position) {
	createTrigger := fyne.NewMenuItem("创建触发器", nil)
	createEffector := fyne.NewMenuItem("创建效果器", nil)

	createTrigger.ChildMenu = &fyne.Menu{
		Items: buildNodeMenuItems(e.nodeManager, "trigger", e, pos),
	}
	createEffector.ChildMenu = &fyne.Menu{
		Items: buildNodeMenuItems(e.nodeManager, "effector", e, pos),
	}

	menu := fyne.NewMenu("", createTrigger, createEffector)
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	widget.ShowPopUpMenuAtPosition(menu, win.Canvas(), pos)
}

func buildNodeMenuItems(manager *node.NodeManager, nodeType string, editor *SkillEditor, pos fyne.Position) []*fyne.MenuItem {
	var items []*fyne.MenuItem
	for i := range manager.Nodes {
		n := manager.Nodes[i]
		if n.Type != nodeType {
			continue
		}
		nodeCopy := n
		item := fyne.NewMenuItem(nodeCopy.Name, func() {
			editor.AddNode(&nodeCopy, pos)
		})
		items = append(items, item)
	}
	return items
}
