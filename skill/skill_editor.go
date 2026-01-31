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
	nodeManager *node.NodeManager //可使用的子节点类型

	viewport          *Viewport
	canvas            *fyne.Container
	scale             float32
	listener          *uihelp.AreaListener
	connectionManager *ConnectionManager // 连接管理器
}

// 节点布局常量
const (
	nodeStartX          = float32(200) // 第一个节点的 X 坐标
	nodeHorizontalSpace = float32(220) // 节点间的水平间隔
	nodeVerticalPos     = float32(300) // 所有节点的 Y 坐标
)

// getNodePosition 根据节点在 Nodes 中的索引计算其位置
func (e *SkillEditor) getNodePosition(nodeIndex int) fyne.Position {
	xPos := nodeStartX + float32(nodeIndex)*nodeHorizontalSpace
	yPos := nodeVerticalPos
	return fyne.NewPos(xPos, yPos)
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

	// 技能节点 - 固定间隔水平排列
	const nodeWidth = float32(180)

	for i := range e.skill.Nodes {
		n := &e.skill.Nodes[i]
		// 根据索引计算位置
		pos := e.getNodePosition(i)

		w := NewNodeWidget(n, e, i)
		// 计算动态高度（基于输入/输出端口数量）
		maxPorts := len(n.NodeData.Inputs)
		if len(n.NodeData.Outputs) > maxPorts {
			maxPorts = len(n.NodeData.Outputs)
		}
		nodeHeight := float32(50 + maxPorts*25)

		w.Resize(fyne.NewSize(nodeWidth, nodeHeight))
		w.Move(pos)
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
	nodeWidth := e.skill.MainNodeWidth
	maxHeight := e.skill.MainNodeHeight

	scroll := container.NewVScroll(form)
	scroll.SetMinSize(fyne.NewSize(nodeWidth, 0)) // 高度由内容决定

	// 计算 Form 的最小高度
	formMin := form.MinSize()
	height := maxHeight
	if (formMin.Height + 80) < height {
		height = formMin.Height + 80
	}
	scroll.Resize(fyne.NewSize(nodeWidth, height))

	// 用 Card 包裹
	mainCard := widget.NewCard("主节点", "", scroll)
	mainCard.Resize(scroll.Size())
	mainCard.Move(fyne.NewPos(e.skill.MainNodeX, e.skill.MainNodeY))

	return mainCard
}

func (e *SkillEditor) SetSkill(skill *Skill) {
	e.skill = skill
	e.Refresh()
}

// 添加节点到画布（位置由排列算法自动计算）
func (e *SkillEditor) AddNode(n *node.Node, pos fyne.Position) {
	if e.skill == nil {
		return
	}

	sn := &SkillNode{
		ID:       tool.GenerateUniqueID(),
		NodeID:   n.ID,
		NodeData: *n,
		Inputs:   map[string]string{},
		Outputs:  map[string][]string{},
	}
	e.skill.Nodes = append(e.skill.Nodes, *sn)
	e.Refresh()
}

// ---------------- 右键菜单 ----------------
// ShowNodeRightClickMenu 右键点击显示菜单
// 只有在以下情况才显示菜单：
// 1. Nodes 为空时，点击主节点区域
// 2. Nodes 不为空时，只能点击最后一个子节点
func (e *SkillEditor) ShowNodeRightClickMenu(ev *fyne.PointEvent) {
	if e.skill == nil {
		return
	}

	// 获取相对于 AreaListener 的点击位置
	clickPos := ev.Position

	// 主节点的位置和大小
	mainNodeX := e.skill.MainNodeX
	mainNodeY := e.skill.MainNodeY
	mainNodeWidth := e.skill.MainNodeWidth
	mainNodeHeight := e.skill.MainNodeHeight

	// 如果 Nodes 为空，检查是否点击在主节点上
	if len(e.skill.Nodes) == 0 {
		inX := clickPos.X >= mainNodeX-5 && clickPos.X <= mainNodeX+mainNodeWidth+5
		inY := clickPos.Y >= mainNodeY-5 && clickPos.Y <= mainNodeY+mainNodeHeight+5
		if inX && inY {
			e.showContextMenu(ev)
		}
		return
	}

	// Nodes 不为空，检查是否点击在最后一个节点上
	if len(e.skill.Nodes) > 0 {
		lastNodePos := e.getNodePosition(len(e.skill.Nodes) - 1)
		lastNode := &e.skill.Nodes[len(e.skill.Nodes)-1]

		// 节点大小（需要动态计算）
		nodeWidth := float32(180)
		maxPorts := len(lastNode.NodeData.Inputs)
		if len(lastNode.NodeData.Outputs) > maxPorts {
			maxPorts = len(lastNode.NodeData.Outputs)
		}
		nodeHeight := float32(50 + maxPorts*25)

		// 检查点击位置是否在最后一个节点范围内
		// 使用更宽松的边界以处理浮点数精度问题
		inX := clickPos.X >= lastNodePos.X-5 && clickPos.X <= lastNodePos.X+nodeWidth+5
		inY := clickPos.Y >= lastNodePos.Y-5 && clickPos.Y <= lastNodePos.Y+nodeHeight+5

		if inX && inY {
			e.showContextMenu(ev)
		}
	}
}

func (e *SkillEditor) showContextMenu(ev *fyne.PointEvent) {
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
