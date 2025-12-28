package skill

import (
	"fmt"
	"skill-editor/model"
	"skill-editor/node"
	"skill-editor/skill/uihelp"
	"skill-editor/tool"
	"strconv"
	"strings"

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
	main := e.CreateMainNodeWidget()
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

func (e *SkillEditor) CreateMainNodeWidget() fyne.CanvasObject {
	if e.skill == nil {
		return widget.NewLabel("未加载技能")
	}

	// ID 和 Name 用 Label 显示（不可编辑）
	idLabel := widget.NewLabel(e.skill.ID)
	nameLabel := widget.NewLabel(e.skill.Name)

	// 可编辑字段
	timesEntry := widget.NewEntry()
	timesEntry.SetText(fmt.Sprintf("%d", e.skill.Times))
	timesEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			e.skill.Times = v
		}
	}

	times5Entry := widget.NewEntry()
	times5Entry.SetText(fmt.Sprintf("%d", e.skill.Times5))
	times5Entry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			e.skill.Times5 = v
		}
	}

	heroEntry := widget.NewEntry()
	heroEntry.SetText(fmt.Sprintf("%d", e.skill.Hero))
	heroEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			e.skill.Hero = v
		}
	}

	masterSkillEntry := widget.NewEntry()
	masterSkillEntry.SetText(fmt.Sprintf("%d", e.skill.MasterSkill))
	masterSkillEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			e.skill.MasterSkill = v
		}
	}

	activeEntry := widget.NewEntry()
	activeEntry.SetText(fmt.Sprintf("%d", e.skill.Active))
	activeEntry.OnChanged = func(s string) {
		if v, err := strconv.Atoi(s); err == nil {
			e.skill.Active = v
		}
	}

	stagesEntry := widget.NewEntry()
	stagesEntry.SetText(fmt.Sprintf("%v", e.skill.Stages))
	stagesEntry.OnChanged = func(s string) {
		parts := strings.Split(s, ",")
		var stages []int
		for _, p := range parts {
			if v, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				stages = append(stages, v)
			}
		}
		e.skill.Stages = stages
	}

	// Form 显示主节点参数
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idLabel},
			{Text: "Name", Widget: nameLabel},
			{Text: "Times", Widget: timesEntry},
			{Text: "Times5", Widget: times5Entry},
			{Text: "Hero", Widget: heroEntry},
			{Text: "MasterSkill", Widget: masterSkillEntry},
			{Text: "Active", Widget: activeEntry},
			{Text: "Stages", Widget: stagesEntry},
		},
	}

	// 固定宽度，高度自适应，超过最大高度滚动
	const nodeWidth = 220
	const maxHeight = 500 // 最大显示高度，超过出现滚动

	scroll := container.NewVScroll(form)
	scroll.SetMinSize(fyne.NewSize(nodeWidth, 0)) // 高度由内容决定

	// 计算 Form 的最小高度
	formMin := form.MinSize()
	height := float32(maxHeight)
	if (formMin.Height + 80) < float32(height) {
		height = formMin.Height + 80
	}
	scroll.Resize(fyne.NewSize(nodeWidth, height))

	// 用 Card 包裹
	mainCard := widget.NewCard("主节点", "", scroll)
	mainCard.Resize(scroll.Size())
	mainCard.Move(fyne.NewPos(100, 50))

	return mainCard
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
	for _, nodeType := range model.NodeTypes {
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
