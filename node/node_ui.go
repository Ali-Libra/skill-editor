package node

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewNodeUI() fyne.CanvasObject {
	manager := NewNodeManager()

	list := widget.NewList(
		func() int {
			return len(manager.Nodes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(manager.Nodes[i].Name)
		},
	)

	editor := container.NewCenter(
		widget.NewLabel("节点编辑区域（待实现）"),
	)

	split := container.NewHSplit(list, editor)
	split.Offset = 0.3

	return split
}
