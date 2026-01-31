package uihelp

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type AreaListener struct {
	widget.BaseWidget
	OnRightClick func(ev *fyne.PointEvent)
}

func NewAreaListener() *AreaListener {
	a := &AreaListener{}
	a.ExtendBaseWidget(a)
	return a
}

func (a *AreaListener) TappedSecondary(ev *fyne.PointEvent) {
	if a.OnRightClick != nil {
		// 创建一个新的 PointEvent，其中 Position 是相对于 AreaListener 的坐标
		newEv := &fyne.PointEvent{
			Position:         ev.Position,         // 相对于 AreaListener 的位置
			AbsolutePosition: ev.AbsolutePosition, // 屏幕绝对位置
		}
		a.OnRightClick(newEv)
	}
}

func (a *AreaListener) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	return &emptyAreaRenderer{rect: rect}
}

type emptyAreaRenderer struct {
	rect *canvas.Rectangle
}

func (r *emptyAreaRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *emptyAreaRenderer) MinSize() fyne.Size {
	return fyne.NewSize(1, 1)
}

func (r *emptyAreaRenderer) Refresh() {}

func (r *emptyAreaRenderer) Destroy() {}

func (r *emptyAreaRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}
