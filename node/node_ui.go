package node

import (
	"skill-editor/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type NodeUI struct {
	manager      *NodeManager
	selectedNode *Node
	editor       *fyne.Container
	list         *widget.List
	Split        *container.Split
}

// 创建 NodeUI
func NewNodeUI() *NodeUI {
	ui := &NodeUI{
		manager: NewNodeManager(),
		editor:  container.NewVBox(widget.NewLabel("请选择或创建一个节点")),
	}

	ui.list = widget.NewList(
		func() int { return len(ui.manager.Nodes) },
		func() fyne.CanvasObject { return NewNodeListItem(nil, ui, 0) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := o.(*NodeListItem)
			item.node = &ui.manager.Nodes[i]
			item.index = i
			item.ui = ui
			item.SetText(item.node.Name)
		},
	)

	// 点击节点时触发
	ui.list.OnSelected = func(id widget.ListItemID) {
		ui.selectedNode = &ui.manager.Nodes[id]
		ui.rebuildEditor()
	}

	createBtn := widget.NewButton("创建新节点", func() { ui.showCreateNodeDialog() })
	left := container.NewBorder(nil, createBtn, nil, nil, ui.list)
	split := container.NewHSplit(left, ui.editor)
	split.Offset = 0.2

	ui.Split = split
	return ui
}

// ---------------- NodeListItem ----------------
type NodeListItem struct {
	widget.Label
	node  *Node
	ui    *NodeUI
	index int
}

func NewNodeListItem(n *Node, ui *NodeUI, index int) *NodeListItem {
	item := &NodeListItem{
		Label: *widget.NewLabel(""),
		node:  n,
		ui:    ui,
		index: index,
	}
	item.ExtendBaseWidget(item)
	return item
}

// 点击节点：同步蓝色选中框
func (i *NodeListItem) Tapped(_ *fyne.PointEvent) {
	if i.node == nil || i.ui == nil {
		return
	}
	i.ui.list.Select(i.index) // 会触发 OnSelected
}

// 右键菜单：重命名 / 删除
func (i *NodeListItem) TappedSecondary(e *fyne.PointEvent) {
	if i.node == nil || i.ui == nil {
		return
	}
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("重命名", func() { i.ui.showRenameNodeDialog(i.node) }),
		fyne.NewMenuItem("删除", func() { i.ui.showDeleteNodeDialog(i.node) }),
	)
	widget.ShowPopUpMenuAtPosition(
		menu,
		fyne.CurrentApp().Driver().CanvasForObject(i),
		e.AbsolutePosition,
	)
}

// ---------------- NodeEditor ----------------
func (ui *NodeUI) Save() error {
	if ui.manager != nil {
		return ui.manager.Save()
	}
	return nil
}
func (ui *NodeUI) rebuildEditor() {
	if ui.selectedNode == nil {
		ui.editor.Objects = []fyne.CanvasObject{widget.NewLabel("请选择或创建一个节点")}
		ui.editor.Refresh()
		return
	}
	ui.editor.Objects = []fyne.CanvasObject{ui.buildNodeEditor(ui.selectedNode)}
	ui.editor.Refresh()
}

func (ui *NodeUI) buildNodeEditor(n *Node) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(n.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	buildPortList := func(
		ports *[]Port,
		titleStr string,
		addLabel string,
		showDefault bool, // 是否显示默认值
	) fyne.CanvasObject {
		items := container.NewVBox()
		for i := range *ports {
			index := i
			p := (*ports)[i]
			displayText := p.Name + " : " + p.Type
			if showDefault {
				displayText += " : " + p.DefaultValue
			}
			name := widget.NewLabel(displayText)
			del := widget.NewButton("删除", func() {
				*ports = append((*ports)[:index], (*ports)[index+1:]...)
				ui.rebuildEditor()
			})
			items.Add(container.NewHBox(name, layout.NewSpacer(), del))
		}
		addBtn := widget.NewButton(addLabel, func() {
			ui.showAddPortDialog(ports, showDefault)
		})
		titleRow := container.NewHBox(widget.NewLabel(titleStr), layout.NewSpacer(), addBtn)
		return container.NewVBox(titleRow, items)
	}

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		buildPortList(&n.Params, "节点参数", "+", true), // 显示默认值
		widget.NewSeparator(),
		buildPortList(&n.Inputs, "输入参数", "+", false), // 不显示默认值
		widget.NewSeparator(),
		buildPortList(&n.Outputs, "输出参数", "+", false), // 不显示默认值
	)
}

// showAddPortDialog 修改为接收 showDefault 参数
func (ui *NodeUI) showAddPortDialog(ports *[]Port, showDefault bool) {
	nameEntry := widget.NewEntry()
	typeSelect := widget.NewSelect([]string{"牌", "int", "float", "string", "bool"}, nil)

	var defaultValWidget fyne.CanvasObject
	defaultValContainer := container.NewVBox()

	if showDefault {
		allCards := model.AllCards()
		defaultValWidget = widget.NewSelect(allCards, nil)
		defaultValWidget.(*widget.Select).Selected = allCards[0] // 默认值
		defaultValContainer.Objects = []fyne.CanvasObject{defaultValWidget}
	} else {
		defaultValContainer.Objects = []fyne.CanvasObject{}
	}

	// Form items
	formItems := []*widget.FormItem{
		widget.NewFormItem("名称", nameEntry),
		widget.NewFormItem("类型", typeSelect),
	}
	if showDefault {
		formItems = append(formItems, widget.NewFormItem("默认值", defaultValContainer))
	}

	form := widget.NewForm(formItems...)

	// 设置默认选中类型
	typeSelect.Selected = "牌"

	// 当类型改变时，动态修改默认值控件（只在 showDefault=true 时）
	typeSelect.OnChanged = func(selected string) {
		if !showDefault {
			return
		}
		var newWidget fyne.CanvasObject
		allCards := model.AllCards()
		switch selected {
		case "int", "float", "string":
			newWidget = widget.NewEntry()
		case "bool":
			w := widget.NewSelect([]string{"true", "false"}, nil)
			w.Selected = "false"
			newWidget = w
		case "牌":
			w := widget.NewSelect(allCards, nil)
			w.Selected = allCards[0]
			newWidget = w
		default:
			newWidget = widget.NewLabel("请选择类型")
		}
		defaultValContainer.Objects = []fyne.CanvasObject{newWidget}
		defaultValContainer.Refresh()
	}

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("添加参数", "确认", "取消", form.Items, func(ok bool) {
		if !ok || nameEntry.Text == "" || typeSelect.Selected == "" {
			return
		}

		var defaultVal string
		if showDefault && len(defaultValContainer.Objects) > 0 {
			switch w := defaultValContainer.Objects[0].(type) {
			case *widget.Entry:
				defaultVal = w.Text
			case *widget.Select:
				defaultVal = w.Selected
			}
		}

		*ports = append(*ports, Port{
			Name:         nameEntry.Text,
			Type:         typeSelect.Selected,
			DefaultValue: defaultVal, // 输入/输出参数默认值为空
		})
		ui.rebuildEditor()
	}, win)
}

func (ui *NodeUI) showCreateNodeDialog() {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	typeSelectEntry := widget.NewSelect(model.NodeTypes, nil)
	typeSelectEntry.Selected = model.NODETYPE_TRIGGER

	form := widget.NewForm(
		widget.NewFormItem("ID", idEntry),
		widget.NewFormItem("名称", nameEntry),
		widget.NewFormItem("类型", typeSelectEntry),
	)

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("创建节点", "确认", "取消", form.Items, func(ok bool) {
		if !ok || idEntry.Text == "" || nameEntry.Text == "" || typeSelectEntry.Selected == "" {
			return
		}

		ui.manager.AddNode(idEntry.Text, nameEntry.Text, typeSelectEntry.Selected)
		ui.list.Refresh()

		// 自动选中新节点，蓝色框 + 编辑区同步
		index := len(ui.manager.Nodes) - 1
		ui.list.Select(index)
	}, win)
}

func (ui *NodeUI) showRenameNodeDialog(n *Node) {
	entry := widget.NewEntry()
	entry.SetText(n.Name)
	form := widget.NewForm(widget.NewFormItem("新名称", entry))
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("重命名节点", "确认", "取消", form.Items, func(ok bool) {
		if !ok || entry.Text == "" {
			return
		}
		n.Name = entry.Text
		ui.rebuildEditor()
		ui.list.Refresh()
	}, win)
}

func (ui *NodeUI) showDeleteNodeDialog(n *Node) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("删除节点", "确定要删除该节点吗？", func(ok bool) {
		if !ok {
			return
		}
		ui.manager.RemoveNode(n)
		if ui.selectedNode == n {
			ui.selectedNode = nil
		}
		ui.rebuildEditor()
		ui.list.Refresh()
	}, win)
}
