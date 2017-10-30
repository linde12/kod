package editor

import "github.com/gdamore/tcell"

type Painter interface {
	SetContent(x int, y int, ch rune, comb []rune, style tcell.Style)
	ShowCursor(x int, y int)
	Size() (int, int)
}

type Viewport struct {
	offx, offy    int
	width, height int
	view          Painter
}

func (v *Viewport) SetContent(x int, y int, ch rune, comb []rune, style tcell.Style) {
	phyx := v.offx + x
	phyy := v.offy + y

	if phyx < v.width && phyy < v.height {
		v.view.SetContent(v.offx+x, v.offy+y, ch, comb, style)
	}
}

func (v *Viewport) ShowCursor(x int, y int) {
	v.view.ShowCursor(v.offx+x, v.offy+y)
}

func (v *Viewport) Size() (int, int) {
	return v.width, v.height
}

func (v *Viewport) FillParent() {
	width, height := v.view.Size()
	v.width = width
	v.height = height
}

func (v *Viewport) SetOffsetX(x int) {
	v.offx = x
}

func (v *Viewport) SetOffsetY(y int) {
	v.offy = y
}

func (v *Viewport) SetWidth(w int) {
	v.width = w
}

func (v *Viewport) SetHeight(h int) {
	v.height = h
}

func NewViewport(parent Painter, offx, offy int) *Viewport {
	width, height := parent.Size()
	return &Viewport{
		view:   parent,
		offx:   offx,
		offy:   offy,
		width:  width - offx,
		height: height - offy,
	}
}
