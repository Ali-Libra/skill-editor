package uihelp

import (
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type HoverButton struct {
	widget.Button
	OnMouseIn  func()
	OnMouseOut func()
}

func NewHoverButton(label string, tapped func()) *HoverButton {
	h := &HoverButton{}
	h.ExtendBaseWidget(h)
	h.SetText(label)
	h.OnTapped = tapped
	return h
}

// 实现 desktop.Hoverable
func (h *HoverButton) MouseIn(*desktop.MouseEvent) {
	if h.OnMouseIn != nil {
		h.OnMouseIn()
	}
}
func (h *HoverButton) MouseOut(*desktop.MouseEvent) {
	if h.OnMouseOut != nil {
		h.OnMouseOut()
	}
}
func (h *HoverButton) MouseMoved(*desktop.MouseEvent) {}
