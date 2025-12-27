package main

import (
	"skill-editor/node"
	"skill-editor/skill"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var curUI string = "技能"

func main() {
	a := app.New()
	w := a.NewWindow("百将牌技能编辑器")
	w.Resize(fyne.NewSize(900, 600))

	content := container.NewStack()

	// 创建技能和节点 UI
	skillUI := skill.NewSkillUI()
	nodeUI := node.NewNodeUI()

	// 默认显示技能界面
	content.Objects = []fyne.CanvasObject{skillUI}

	// 顶部按钮：技能 / 节点
	btnSkill := widget.NewButton("技能", func() {
		curUI = "技能"
		content.Objects = []fyne.CanvasObject{skillUI}
		content.Refresh()
	})

	btnNode := widget.NewButton("节点", func() {
		curUI = "节点"
		content.Objects = []fyne.CanvasObject{nodeUI.Split}
		content.Refresh()
	})

	// 保存按钮
	btnSave := widget.NewButton("保存", func() {
		if curUI == "技能" {
			// skillUI.Save()
		} else if curUI == "节点" {
			nodeUI.Save()
		}
	})

	// 顶部栏布局：技能 / 节点 / 右对齐保存按钮
	topBar := container.NewHBox(
		btnSkill,
		btnNode,
		layout.NewSpacer(),
		btnSave,
	)

	w.SetContent(container.NewBorder(topBar, nil, nil, nil, content))
	w.ShowAndRun()
}
