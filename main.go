package main

import (
	"skill-editor/node"
	"skill-editor/skill"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("技能编辑器")
	w.Resize(fyne.NewSize(900, 600))

	content := container.NewStack()

	// 默认显示技能界面
	skillUI := skill.NewSkillUI()
	nodeUI := node.NewNodeUI()

	content.Objects = []fyne.CanvasObject{skillUI}

	btnSkill := widget.NewButton("技能", func() {
		content.Objects = []fyne.CanvasObject{skillUI}
		content.Refresh()
	})

	btnNode := widget.NewButton("节点", func() {
		content.Objects = []fyne.CanvasObject{nodeUI}
		content.Refresh()
	})

	topBar := container.NewHBox(btnSkill, btnNode)

	w.SetContent(container.NewBorder(topBar, nil, nil, nil, content))
	w.ShowAndRun()
}
