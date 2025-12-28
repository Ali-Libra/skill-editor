package skill

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Viewport 包装一个 container，用于支持拖动和缩放
type Viewport struct {
	widget.BaseWidget
	content *fyne.Container
	offset  fyne.Position
	scale   float32
	lastPos *fyne.Position
}

// NewViewport 创建 Viewport
func NewViewport(content *fyne.Container) *Viewport {
	v := &Viewport{
		content: content,
		scale:   1.0,
	}
	v.ExtendBaseWidget(v)
	return v
}

// CreateRenderer 渲染内容
func (v *Viewport) CreateRenderer() fyne.WidgetRenderer {
	return &viewportRenderer{v: v}
}

type viewportRenderer struct {
	v *Viewport
}

func (r *viewportRenderer) Layout(size fyne.Size) {
	c := r.v.content
	scaled := fyne.NewSize(c.Size().Width*r.v.scale, c.Size().Height*r.v.scale)
	c.Resize(scaled)
	c.Move(r.v.offset)
}

func (r *viewportRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (r *viewportRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.v.content}
}

func (r *viewportRenderer) Refresh() {
	r.Layout(r.v.Size())
}

func (r *viewportRenderer) Destroy() {}

// ----------------- 拖动 -----------------
func (v *Viewport) Dragged(e *fyne.DragEvent) {
	if v.lastPos == nil {
		p := e.Position
		v.lastPos = &p
		return
	}

	dx := e.Position.X - v.lastPos.X
	dy := e.Position.Y - v.lastPos.Y

	v.offset = v.offset.Add(fyne.NewPos(dx, dy))
	v.lastPos = &e.Position
	v.Refresh()
}

func (v *Viewport) DragEnd() {
	v.lastPos = nil
}

// ----------------- 缩放 -----------------
func (v *Viewport) Scrolled(e *fyne.ScrollEvent) {
	if e.Scrolled.DY == 0 {
		return
	}

	scale := v.scale
	if e.Scrolled.DY > 0 {
		scale *= 1.1
	} else {
		scale /= 1.1
	}

	if scale < 0.3 {
		scale = 0.3
	} else if scale > 2.5 {
		scale = 2.5
	}

	v.scale = scale
	v.Refresh()
}

// 获取缩放比例
func (v *Viewport) Scale() float32 {
	return v.scale
}

// 获取偏移
func (v *Viewport) Offset() fyne.Position {
	return v.offset
}
