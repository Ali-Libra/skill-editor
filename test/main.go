package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 1. 创建应用和窗口
	a := app.New()
	w := a.NewWindow("Fyne 右键多级菜单示例")
	w.Resize(fyne.NewSize(600, 400))

	// 2. 创建 Label 内容（纯展示）
	labelContent := widget.NewLabel("在窗口内右键点击，体验多级菜单\n悬停在「操作」「设置」上可展开二级菜单")
	labelContent.Alignment = fyne.TextAlignCenter

	// 3. 用 Container 包裹 Label（通用容器，支持所有事件）
	mainContent := container.NewCenter(labelContent)
	mainContent.Resize(fyne.NewSize(600, 400))

	// 4. 构建右键多级菜单
	contextMenu := buildContextMenu()

	// 创建一个透明按钮覆盖整个容器，捕获右键事件
	clickItem := NewEmptyButton(func(ev *fyne.PointEvent) {
		// 在鼠标右键点击位置弹出菜单
		widget.ShowPopUpMenuAtPosition(contextMenu, w.Canvas(), ev.AbsolutePosition)
	})
	clickItem.Resize(fyne.NewSize(600, 400))

	content := container.NewWithoutLayout(
		clickItem, // 覆盖在最上层
		mainContent,
	)

	// 5. 设置窗口内容并运行
	w.SetContent(content)
	w.ShowAndRun()
}

// buildContextMenu 构建右键多级菜单结构
func buildContextMenu() *fyne.Menu {
	// 二级子菜单：操作 -> 新建/删除/复制
	operateSubMenu := fyne.NewMenu("操作子菜单",
		fyne.NewMenuItem("新建文件", func() {
			fmt.Println("右键菜单：操作 -> 新建文件")
		}),
		fyne.NewMenuItem("删除选中", func() {
			fmt.Println("右键菜单：操作 -> 删除选中")
		}),
		fyne.NewMenuItemSeparator(), // 菜单分隔线
		fyne.NewMenuItem("复制内容", func() {
			fmt.Println("右键菜单：操作 -> 复制内容")
		}),
	)

	// 二级子菜单：设置 -> 主题/字号
	settingSubMenu := fyne.NewMenu("设置子菜单",
		fyne.NewMenuItem("切换深色主题", func() {
			fmt.Println("右键菜单：设置 -> 切换深色主题")
		}),
		fyne.NewMenuItem("增大字号", func() {
			fmt.Println("右键菜单：设置 -> 增大字号")
		}),
		fyne.NewMenuItem("减小字号", func() {
			fmt.Println("右键菜单：设置 -> 减小字号")
		}),
	)

	operateMenu := fyne.NewMenuItem("操作", nil)
	operateMenu.ChildMenu = operateSubMenu
	settingMenu := fyne.NewMenuItem("设置", nil)
	settingMenu.ChildMenu = settingSubMenu
	// 一级右键菜单（包含带二级子菜单的项）
	contextMenu := fyne.NewMenu("", // 主菜单名称为空（右键菜单无需标题）
		operateMenu,
		settingMenu,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("刷新页面", func() {
			fmt.Println("右键菜单：刷新页面")
		}),
		fyne.NewMenuItem("关闭窗口", func() {
			fmt.Println("右键菜单：关闭窗口")
		}),
	)

	return contextMenu
}

type EmptyButton struct {
	widget.BaseWidget
	OnRightClick func(ev *fyne.PointEvent)
}

func NewEmptyButton(onRightClick func(ev *fyne.PointEvent)) *EmptyButton {
	b := &EmptyButton{
		OnRightClick: onRightClick,
	}
	b.ExtendBaseWidget(b)
	// ⚠️ 此时 CreateRenderer 已“存在”
	return b
}

func (b *EmptyButton) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	return &emptyButtonRenderer{rect: rect}
}
func (b *EmptyButton) TappedSecondary(ev *fyne.PointEvent) {
	if b.OnRightClick != nil {
		b.OnRightClick(ev)
	}
}

type emptyButtonRenderer struct {
	rect *canvas.Rectangle
}

func (r *emptyButtonRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *emptyButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (r *emptyButtonRenderer) Refresh() {}

func (r *emptyButtonRenderer) Destroy() {}

func (r *emptyButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}
