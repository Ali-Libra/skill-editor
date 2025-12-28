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
	nodeManager := node.NewNodeManager()
	skillUI := skill.NewSkillUI(nodeManager)
	nodeUI := node.NewNodeUI()

	// 默认显示技能界面
	content.Objects = []fyne.CanvasObject{skillUI.Split}

	// 顶部按钮：技能 / 节点
	btnSkill := widget.NewButton("技能", func() {
		if curUI == "技能" {
			return
		}
		curUI = "技能"
		content.Objects = []fyne.CanvasObject{skillUI.Split}
		content.Refresh()
	})

	btnNode := widget.NewButton("节点", func() {
		if curUI == "节点" {
			return
		}
		curUI = "节点"
		content.Objects = []fyne.CanvasObject{nodeUI.Split}
		content.Refresh()
	})

	// 保存按钮
	btnSave := widget.NewButton("保存", func() {
		switch curUI {
		case "技能":
			skillUI.Save()
		case "节点":
			nodeUI.Save()
		}
	})

	// 顶部栏布局：技能 / 节点 / 右对齐保存按钮
	topBar := container.NewHBox(
		btnSkill,
		btnNode,
		btnSave,
		layout.NewSpacer(),
	)

	w.SetContent(container.NewBorder(topBar, nil, nil, nil, content))
	w.ShowAndRun()
}
