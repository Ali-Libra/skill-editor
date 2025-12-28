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

	editorWidget *SkillEditorCanvas
}

func NewSkillEditor(nodeMgr *node.NodeManager) *SkillEditor {
	canvas := container.NewWithoutLayout()
	editor := &SkillEditor{
		nodeManager: nodeMgr,
		canvas:      canvas,
		scale:       1.0,
	}

	editorWidget := NewSkillEditorCanvas(editor)
	editor.editorWidget = editorWidget
	editorWidget.Resize(fyne.NewSize(2000, 2000)) // 或者根据 viewport 大小动态调整
	editor.canvas.Add(editorWidget)

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
	newObjects := []fyne.CanvasObject{e.editorWidget}

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
func (e *SkillEditor) ShowNodeRightClickMenu(pos fyne.Position) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	// 一级菜单按钮
	triggerBtn := uihelp.NewHoverButton("创建触发器", nil)
	effectBtn := uihelp.NewHoverButton("创建效果器", nil)

	// 一级菜单容器
	menuContainer := container.NewVBox(triggerBtn, effectBtn)
	popup := widget.NewPopUp(menuContainer, win.Canvas())
	popup.Move(pos)
	popup.Show()

	// 一级悬停显示二级
	triggerBtn.OnMouseIn = func() {
		e.showSubMenu("trigger", pos, triggerBtn, popup)
	}
	effectBtn.OnMouseIn = func() {
		e.showSubMenu("effect", pos, effectBtn, popup)
	}
}

// 显示二级菜单
func (e *SkillEditor) showSubMenu(nodeType string, pos fyne.Position, parentBtn *uihelp.HoverButton, parentPopup *widget.PopUp) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	buttons := e.buildNodeMenuButtons(nodeType, pos)
	canvasObjs := make([]fyne.CanvasObject, len(buttons))
	for i, b := range buttons {
		canvasObjs[i] = b
	}

	subMenu := container.NewVBox(canvasObjs...)
	subPopup := widget.NewPopUp(subMenu, win.Canvas())

	// 位置在父按钮右侧
	subPopup.Move(parentBtn.Position().Add(fyne.NewPos(parentBtn.Size().Width, 0)))
	subPopup.Show()
}

func (e *SkillEditor) buildNodeMenuButtons(nodeType string, pos fyne.Position) []*widget.Button {
	var buttons []*widget.Button
	for i := range e.nodeManager.Nodes {
		n := e.nodeManager.Nodes[i]
		if n.Type != nodeType {
			continue
		}
		// 捕获 n
		nodeCopy := n
		btn := widget.NewButton(nodeCopy.Name, func() {
			e.AddNode(&nodeCopy, pos)
		})
		buttons = append(buttons, btn)
	}
	return buttons
}

// ---------------- 自定义 Canvas Widget ----------------
type SkillEditorCanvas struct {
	widget.BaseWidget
	editor *SkillEditor
}

func NewSkillEditorCanvas(editor *SkillEditor) *SkillEditorCanvas {
	c := &SkillEditorCanvas{editor: editor}
	c.ExtendBaseWidget(c)
	return c
}

// 右键
func (c *SkillEditorCanvas) TappedSecondary(ev *fyne.PointEvent) {
	c.editor.ShowNodeRightClickMenu(ev.AbsolutePosition)
}

// 拖动画布
func (c *SkillEditorCanvas) Dragged(e *fyne.DragEvent) {
	c.editor.viewport.Dragged(e)
}

func (c *SkillEditorCanvas) DragEnd() {
	c.editor.viewport.DragEnd()
}

// 缩放
func (c *SkillEditorCanvas) Scrolled(e *fyne.ScrollEvent) {
	c.editor.viewport.Scrolled(e)
}

// 必须实现 CreateRenderer
func (c *SkillEditorCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &skillEditorCanvasRenderer{}
}

type skillEditorCanvasRenderer struct{}

func (r *skillEditorCanvasRenderer) Layout(size fyne.Size)        {}
func (r *skillEditorCanvasRenderer) MinSize() fyne.Size           { return fyne.NewSize(1, 1) }
func (r *skillEditorCanvasRenderer) Refresh()                     {}
func (r *skillEditorCanvasRenderer) Objects() []fyne.CanvasObject { return nil }
func (r *skillEditorCanvasRenderer) Destroy()                     {}
