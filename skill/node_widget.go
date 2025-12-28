package skill

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NodeWidget 是画布上的节点控件
type NodeWidget struct {
	widget.BaseWidget
	node   *SkillNode
	editor *SkillEditor
	pos    fyne.Position
}

// NewNodeWidget 创建一个节点控件
func NewNodeWidget(node *SkillNode, editor *SkillEditor) *NodeWidget {
	w := &NodeWidget{
		node:   node,
		editor: editor,
		pos:    fyne.NewPos(node.X, node.Y),
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer 实现 widget.Widget 接口
func (w *NodeWidget) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabel(w.node.NodeData.Name)
	rect := container.NewMax(widget.NewLabel(""), label) // 简单容器
	return &nodeWidgetRenderer{
		widget: w,
		object: rect,
	}
}

type nodeWidgetRenderer struct {
	widget *NodeWidget
	object *fyne.Container
}

func (r *nodeWidgetRenderer) Layout(size fyne.Size)        {}
func (r *nodeWidgetRenderer) MinSize() fyne.Size           { return fyne.NewSize(180, 100) }
func (r *nodeWidgetRenderer) Objects() []fyne.CanvasObject { return r.object.Objects }
func (r *nodeWidgetRenderer) Refresh()                     {}
func (r *nodeWidgetRenderer) Destroy()                     {}

// ----------------- 拖动 -----------------
func (w *NodeWidget) Dragged(e *fyne.DragEvent) {
	if w.pos == fyne.NewPos(0, 0) { // 初始化位置
		w.pos = w.Position()
	}

	// 计算增量
	dx := e.Position.X - w.pos.X
	dy := e.Position.Y - w.pos.Y

	w.pos = w.pos.Add(fyne.NewPos(dx, dy))
	w.Move(w.pos)

	w.node.X = w.pos.X
	w.node.Y = w.pos.Y
}

func (w *NodeWidget) DragEnd() {}
