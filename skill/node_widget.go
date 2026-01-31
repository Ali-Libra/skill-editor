package skill

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// NodeWidget 是画布上的节点控件
type NodeWidget struct {
	widget.BaseWidget
	node      *SkillNode
	editor    *SkillEditor
	nodeIndex int // 节点在 Nodes 中的索引
	pos       fyne.Position

	// 端口信息
	inputPorts  map[string]fyne.Position // 输入端口名称 -> 位置
	outputPorts map[string]fyne.Position // 输出端口名称 -> 位置
}

// NewNodeWidget 创建一个节点控件
func NewNodeWidget(node *SkillNode, editor *SkillEditor, nodeIndex int) *NodeWidget {
	pos := editor.getNodePosition(nodeIndex)
	w := &NodeWidget{
		node:        node,
		editor:      editor,
		nodeIndex:   nodeIndex,
		pos:         pos,
		inputPorts:  make(map[string]fyne.Position),
		outputPorts: make(map[string]fyne.Position),
	}
	w.ExtendBaseWidget(w)
	return w
}

// CreateRenderer 实现 widget.Widget 接口
func (w *NodeWidget) CreateRenderer() fyne.WidgetRenderer {
	return &nodeWidgetRenderer{
		widget: w,
	}
}

type nodeWidgetRenderer struct {
	widget *NodeWidget
}

func (r *nodeWidgetRenderer) Layout(size fyne.Size) {
	// 计算端口位置
	const portHeight = 20
	const portRadius = 5
	const leftMargin = 10
	const rightMargin = 10

	// 左侧输入端口
	inputY := float32(25)
	for _, input := range r.widget.node.NodeData.Inputs {
		r.widget.inputPorts[input.Name] = fyne.NewPos(leftMargin, inputY)
		inputY += portHeight + 5
	}

	// 右侧输出端口
	outputY := float32(25)
	for _, output := range r.widget.node.NodeData.Outputs {
		r.widget.outputPorts[output.Name] = fyne.NewPos(size.Width-rightMargin-2*portRadius, outputY)
		outputY += portHeight + 5
	}
}

func (r *nodeWidgetRenderer) MinSize() fyne.Size {
	// 计算最小高度（根据端口数量）
	maxPorts := len(r.widget.node.NodeData.Inputs)
	if len(r.widget.node.NodeData.Outputs) > maxPorts {
		maxPorts = len(r.widget.node.NodeData.Outputs)
	}

	height := float32(50 + maxPorts*25)
	return fyne.NewSize(180, height)
}

func (r *nodeWidgetRenderer) Objects() []fyne.CanvasObject {
	objs := []fyne.CanvasObject{}

	// 节点背景
	bg := canvas.NewRectangle(color.NRGBA{R: 230, G: 230, B: 230, A: 255})
	objs = append(objs, bg)

	// 节点标题
	title := canvas.NewText(r.widget.node.NodeData.Name, color.White)
	title.Move(fyne.NewPos(10, 5))
	objs = append(objs, title)

	// 输入端口
	for name, pos := range r.widget.inputPorts {
		portCircle := canvas.NewCircle(color.NRGBA{R: 100, G: 200, B: 100, A: 255})
		portCircle.Resize(fyne.NewSize(10, 10))
		portCircle.Move(pos)
		objs = append(objs, portCircle)

		portLabel := canvas.NewText(name, color.White)
		portLabel.Move(pos.Add(fyne.NewPos(15, 0)))
		objs = append(objs, portLabel)
	}

	// 输出端口
	for name, pos := range r.widget.outputPorts {
		portCircle := canvas.NewCircle(color.NRGBA{R: 200, G: 100, B: 100, A: 255})
		portCircle.Resize(fyne.NewSize(10, 10))
		portCircle.Move(pos)
		objs = append(objs, portCircle)

		portLabel := canvas.NewText(name, color.White)
		textSize := portLabel.MinSize()
		portLabel.Move(pos.Add(fyne.NewPos(-textSize.Width-5, 0)))
		objs = append(objs, portLabel)
	}

	return objs
}

func (r *nodeWidgetRenderer) Refresh() {}
func (r *nodeWidgetRenderer) Destroy() {}

// Dragged 处理端口拖动（用于连线）
func (w *NodeWidget) Dragged(e *fyne.DragEvent) {
	// 检查是否点击在输出端口上
	for portName, portPos := range w.outputPorts {
		distance := fyne.NewPos(
			e.Position.X-portPos.X,
			e.Position.Y-portPos.Y,
		)
		if distance.X*distance.X+distance.Y*distance.Y < 100 { // 半径 5 的圆
			// 开始拖动连线
			if w.editor.connectionManager == nil {
				w.editor.connectionManager = NewConnectionManager(w.editor)
			}
			port := &PortLocation{
				SkillNodeID: w.node.ID,
				PortName:    portName,
				PortType:    "output",
				Position:    w.Position().Add(portPos),
			}
			w.editor.connectionManager.StartDragFromPort(port)
			return
		}
	}

	// 检查是否点击在输入端口上
	for portName, portPos := range w.inputPorts {
		distance := fyne.NewPos(
			e.Position.X-portPos.X,
			e.Position.Y-portPos.Y,
		)
		if distance.X*distance.X+distance.Y*distance.Y < 100 { // 半径 5 的圆
			// 开始拖动连线
			if w.editor.connectionManager == nil {
				w.editor.connectionManager = NewConnectionManager(w.editor)
			}
			port := &PortLocation{
				SkillNodeID: w.node.ID,
				PortName:    portName,
				PortType:    "input",
				Position:    w.Position().Add(portPos),
			}
			w.editor.connectionManager.StartDragFromPort(port)
			return
		}
	}
}

func (w *NodeWidget) DragEnd() {
	if w.editor.connectionManager != nil {
		w.editor.connectionManager.EndDrag()
	}
}

// Tapped 处理端口点击以完成连接
func (w *NodeWidget) Tapped(e *fyne.PointEvent) {
	// 检查是否在输出端口上
	for portName, portPos := range w.outputPorts {
		distance := fyne.NewPos(
			e.Position.X-portPos.X,
			e.Position.Y-portPos.Y,
		)
		if distance.X*distance.X+distance.Y*distance.Y < 100 {
			if w.editor.connectionManager != nil && w.editor.connectionManager.isDrawingLine {
				port := &PortLocation{
					SkillNodeID: w.node.ID,
					PortName:    portName,
					PortType:    "output",
				}
				w.editor.connectionManager.ConnectPorts(port)
				w.editor.Refresh()
			}
			return
		}
	}

	// 检查是否在输入端口上
	for portName, portPos := range w.inputPorts {
		distance := fyne.NewPos(
			e.Position.X-portPos.X,
			e.Position.Y-portPos.Y,
		)
		if distance.X*distance.X+distance.Y*distance.Y < 100 {
			if w.editor.connectionManager != nil && w.editor.connectionManager.isDrawingLine {
				port := &PortLocation{
					SkillNodeID: w.node.ID,
					PortName:    portName,
					PortType:    "input",
				}
				w.editor.connectionManager.ConnectPorts(port)
				w.editor.Refresh()
			}
			return
		}
	}
}

// 注：节点不可拖动移动，位置由排列算法固定
