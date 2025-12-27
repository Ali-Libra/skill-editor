package skill

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type SkillUI struct {
	manager       *SkillManager
	selectedSkill *Skill
	editor        *fyne.Container
	list          *widget.List
	Split         *container.Split
}

func NewSkillUI() *SkillUI {
	ui := &SkillUI{
		manager: NewSkillManager(),
		editor:  container.NewVBox(widget.NewLabel("请选择或创建一个技能")),
	}

	ui.list = widget.NewList(
		func() int { return len(ui.manager.Skills) },
		func() fyne.CanvasObject { return NewSkillListItem(nil, ui, 0) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := o.(*SkillListItem)
			item.skill = &ui.manager.Skills[i]
			item.index = i
			item.ui = ui
			item.SetText(item.skill.Name)
		},
	)

	ui.list.OnSelected = func(id widget.ListItemID) {
		ui.selectedSkill = &ui.manager.Skills[id]
		ui.rebuildEditor()
	}

	createBtn := widget.NewButton("创建新技能", func() {
		ui.showCreateSkillDialog()
	})

	left := container.NewBorder(nil, createBtn, nil, nil, ui.list)
	split := container.NewHSplit(left, ui.editor)
	split.Offset = 0.3

	ui.Split = split
	return ui
}

// ---------------- SkillListItem ----------------
type SkillListItem struct {
	widget.Label
	skill *Skill
	ui    *SkillUI
	index int
}

func NewSkillListItem(s *Skill, ui *SkillUI, index int) *SkillListItem {
	item := &SkillListItem{
		Label: *widget.NewLabel(""),
		skill: s,
		ui:    ui,
		index: index,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (i *SkillListItem) Tapped(_ *fyne.PointEvent) {
	if i.skill == nil || i.ui == nil {
		return
	}
	i.ui.list.Select(i.index)
}

func (i *SkillListItem) TappedSecondary(e *fyne.PointEvent) {
	if i.skill == nil || i.ui == nil {
		return
	}

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("重命名", func() { i.ui.showRenameSkillDialog(i.skill) }),
		fyne.NewMenuItem("删除", func() { i.ui.showDeleteSkillDialog(i.skill) }),
	)

	widget.ShowPopUpMenuAtPosition(
		menu,
		fyne.CurrentApp().Driver().CanvasForObject(i),
		e.AbsolutePosition,
	)
}

// ---------------- Skill Editor ----------------
func (ui *SkillUI) rebuildEditor() {
	if ui.selectedSkill == nil {
		ui.editor.Objects = []fyne.CanvasObject{widget.NewLabel("请选择或创建一个技能")}
		ui.editor.Refresh()
		return
	}

	// 这里编辑区先留空，只显示技能名称
	ui.editor.Objects = []fyne.CanvasObject{
		widget.NewLabel(ui.selectedSkill.Name),
	}
	ui.editor.Refresh()
}

// ---------------- Dialogs ----------------
func (ui *SkillUI) showCreateSkillDialog() {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("ID", idEntry),
		widget.NewFormItem("名称", nameEntry),
	)

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("创建技能", "确认", "取消", form.Items, func(ok bool) {
		if !ok || idEntry.Text == "" || nameEntry.Text == "" {
			return
		}

		ui.manager.AddSkill(idEntry.Text, nameEntry.Text)
		ui.list.Refresh()
		index := len(ui.manager.Skills) - 1
		ui.list.Select(index)
	}, win)
}

func (ui *SkillUI) showRenameSkillDialog(s *Skill) {
	entry := widget.NewEntry()
	entry.SetText(s.Name)
	form := widget.NewForm(widget.NewFormItem("新名称", entry))
	win := fyne.CurrentApp().Driver().AllWindows()[0]

	dialog.ShowForm("重命名技能", "确认", "取消", form.Items, func(ok bool) {
		if !ok || entry.Text == "" {
			return
		}
		s.Name = entry.Text
		ui.rebuildEditor()
		ui.list.Refresh()
	}, win)
}

func (ui *SkillUI) showDeleteSkillDialog(s *Skill) {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("删除技能", "确定要删除该技能吗？", func(ok bool) {
		if !ok {
			return
		}
		ui.manager.RemoveSkill(s)
		if ui.selectedSkill == s {
			ui.selectedSkill = nil
		}
		ui.rebuildEditor()
		ui.list.Refresh()
	}, win)
}

// ---------------- Save 方法 ----------------
func (ui *SkillUI) Save() error {
	return ui.manager.Save()
}
