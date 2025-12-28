package skill

import (
	"skill-editor/node"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SkillUI 是技能编辑页面主界面
type SkillUI struct {
	manager       *SkillManager
	selectedSkill *Skill

	editor *SkillEditor
	list   *widget.List
	Split  *container.Split
}

// NewSkillUI 创建 SkillUI
func NewSkillUI(nodeManager *node.NodeManager) *SkillUI {
	manager := NewSkillManager()
	editor := NewSkillEditor(nodeManager)

	skillUI := &SkillUI{
		manager: manager,
		editor:  editor,
	}

	// 技能列表
	list := widget.NewList(
		func() int { return len(manager.Skills) },
		func() fyne.CanvasObject {
			return NewSkillListItem(nil, skillUI, 0)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := o.(*SkillListItem)
			item.skill = &manager.Skills[i]
			item.index = i
			item.ui = skillUI
			item.SetText(item.skill.Name)
		},
	)
	skillUI.list = list

	list.OnSelected = func(id widget.ListItemID) {
		editor.SetSkill(&manager.Skills[id])
	}

	// 创建新技能按钮
	createBtn := widget.NewButton("创建新技能", func() {
		skillUI.showCreateSkillDialog()
	})

	// 左侧列表 + 创建按钮
	left := container.NewBorder(nil, createBtn, nil, nil, list)

	// 右侧编辑器
	right := editor.Root()

	// 水平分割
	skillUI.Split = container.NewHSplit(left, right)
	skillUI.Split.Offset = 0.2

	// ✅ 默认选中第一个技能（如果有）
	if len(manager.Skills) > 0 {
		list.Select(0)                      // 选中列表第一个
		editor.SetSkill(&manager.Skills[0]) // 刷新右侧编辑器
	}
	return skillUI
}

func (ui *SkillUI) showCreateSkillDialog() {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("ID", idEntry),
		widget.NewFormItem("名称", nameEntry),
	)

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("创建新技能", "确认", "取消", form.Items, func(ok bool) {
		if !ok || idEntry.Text == "" || nameEntry.Text == "" {
			return
		}

		ui.manager.AddSkill(idEntry.Text, nameEntry.Text)
		ui.list.Refresh()

		// 自动选中新技能
		index := len(ui.manager.Skills) - 1
		ui.list.Select(index)
	}, win)
}

func (ui *SkillUI) showRenameSkillDialog(s *Skill) {
	entry := widget.NewEntry()
	entry.SetText(s.Name)

	form := widget.NewForm(
		widget.NewFormItem("新名称", entry),
	)

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("重命名技能", "确认", "取消", form.Items, func(ok bool) {
		if !ok || entry.Text == "" {
			return
		}
		s.Name = entry.Text
		ui.list.Refresh()   // 刷新技能列表
		ui.editor.Refresh() // 刷新右侧编辑区
	}, win)
}

func (ui *SkillUI) showDeleteSkillDialog(s *Skill) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("删除技能", "确定要删除该技能吗？", func(ok bool) {
		if !ok {
			return
		}
		ui.manager.RemoveSkill(s)
		ui.list.Refresh()
		ui.editor.SetSkill(nil)
	}, win)
}

func (ui *SkillUI) Save() error {
	if ui.manager == nil {
		return nil
	}
	return ui.manager.Save() // 调用 SkillManager 保存
}

type SkillListItem struct {
	widget.Label
	skill *Skill
	ui    *SkillUI
	index int
}

func NewSkillListItem(skill *Skill, ui *SkillUI, index int) *SkillListItem {
	item := &SkillListItem{
		Label: *widget.NewLabel(""),
		skill: skill,
		ui:    ui,
		index: index,
	}
	item.ExtendBaseWidget(item)
	return item
}

// 普通点击
func (i *SkillListItem) Tapped(_ *fyne.PointEvent) {
	i.ui.list.Select(i.index)
}

// 右键菜单
func (i *SkillListItem) TappedSecondary(e *fyne.PointEvent) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("重命名", func() { i.ui.showRenameSkillDialog(i.skill) }),
		fyne.NewMenuItem("删除", func() { i.ui.showDeleteSkillDialog(i.skill) }),
	)

	win := fyne.CurrentApp().Driver().AllWindows()[0]
	widget.ShowPopUpMenuAtPosition(menu, win.Canvas(), e.AbsolutePosition)
}
