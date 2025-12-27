package skill

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSkillUI() fyne.CanvasObject {
	manager := NewSkillManager()

	list := widget.NewList(
		func() int {
			return len(manager.Skills)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(manager.Skills[i].Name)
		},
	)

	editor := container.NewCenter(
		widget.NewLabel("技能编辑区域（待实现）"),
	)

	split := container.NewHSplit(list, editor)
	split.Offset = 0.3

	return split
}
