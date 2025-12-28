package skill

import (
	"skill-editor/node"
	"skill-editor/skill/uihelp"
	"skill-editor/tool"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SkillEditor struct {
	skill       *Skill
	nodeManager *node.NodeManager

	viewport *Viewport
	canvas   *fyne.Container
	scale    float32

	listener *uihelp.AreaListener
}

func NewSkillEditor(nodeMgr *node.NodeManager) *SkillEditor {
	canvas := container.NewWithoutLayout()
	editor := &SkillEditor{
		nodeManager: nodeMgr,
		canvas:      canvas,
		scale:       1.0,
	}

	editor.listener = uihelp.NewAreaListener()
	editor.listener.OnRightClick = editor.ShowNodeRightClickMenu
	editor.listener.Resize(fyne.NewSize(2000, 2000))

	editor.canvas.Add(editor.listener)
	editor.viewport = NewViewport(canvas)
	return editor
}

func (e *SkillEditor) Root() fyne.CanvasObject {
	return e.viewport
}

func (e *SkillEditor) Refresh() {
	if e.skill == nil {
		return
	}

	// 清空 canvas 中原来的节点（保留 editorWidget）
	newObjects := []fyne.CanvasObject{e.listener}

	// 主节点
	main := widget.NewCard(e.skill.Name, "", widget.NewLabel("技能主节点"))
	main.Resize(fyne.NewSize(150, 200))
	main.Move(fyne.NewPos(100, 50))
	newObjects = append(newObjects, main)

	// 技能节点
	for i := range e.skill.Nodes {
		n := &e.skill.Nodes[i]
		w := NewNodeWidget(n, e)
		w.Resize(fyne.NewSize(180, 100))
		w.Move(fyne.NewPos(n.X, n.Y))
		newObjects = append(newObjects, w)
	}

	e.canvas.Objects = newObjects
	e.canvas.Refresh()
}

func (e *SkillEditor) SetSkill(skill *Skill) {
	e.skill = skill
	e.Refresh()
}

// 添加节点到画布（鼠标点击位置）
func (e *SkillEditor) AddNode(n *node.Node, pos fyne.Position) {
	if e.skill == nil {
		return
	}

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
func (e *SkillEditor) ShowNodeRightClickMenu(ev *fyne.PointEvent) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	contextMenu := fyne.NewMenu("")
	for _, nodeType := range node.NodeTypes {
		mainMenu := fyne.NewMenuItem(nodeType, nil)
		mainMenu.ChildMenu = e.getChildMenu(nodeType, ev.AbsolutePosition)
		contextMenu.Items = append(contextMenu.Items, mainMenu)
	}

	widget.ShowPopUpMenuAtPosition(
		contextMenu,
		win.Canvas(),
		ev.AbsolutePosition,
	)
}
func (e *SkillEditor) getChildMenu(nodeType string, pos fyne.Position) *fyne.Menu {
	menus := fyne.NewMenu("")
	for i := range e.nodeManager.Nodes {
		n := e.nodeManager.Nodes[i]
		if n.Type != nodeType {
			continue
		}
		nodeCopy := n
		menu := fyne.NewMenuItem(nodeCopy.Name, func() {
			e.AddNode(&nodeCopy, pos)
		})
		menus.Items = append(menus.Items, menu)
	}
	return menus
}
