package editor

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type TextView struct {
	x, y int
}

func (t *TextView) Draw() {
	panic("not implemented")
}

func (t *TextView) Resize() {
	panic("not implemented")
}

func (t *TextView) HandleEvent(ev tcell.Event) bool {
	panic("not implemented")
}

func (t *TextView) SetView(view views.View) {
	panic("not implemented")
}

func (t *TextView) Size() (int, int) {
	return t.x, t.y
}

func (t *TextView) Watch(handler tcell.EventHandler) {
	panic("not implemented")
}

func (t *TextView) Unwatch(handler tcell.EventHandler) {
	panic("not implemented")
}

type TextModel struct {
	lc *LineCache
}

//func (v *View) Draw() {
//// TODO: tcell setcontent
//if len(v.lines) == 0 {
//return
//}

//// TODO: Line numbers
//// TODO: Fix choppy scrolling
//for y, line := range v.lines {
//visualX := 0
//for x, char := range []rune(line.Text) {
//if char == '\t' {
//ts := tabSize - (visualX % tabSize)
//for i := 0; i < ts; i++ {
//v.Editor.screen.SetContent(visualX+i, y, ' ', nil, stylemap[line.StyleIds[x]])
//}
//visualX += ts
//} else if char != '\n' {
//// TODO: Trim newline in a better way?
//v.Editor.screen.SetContent(x, y, char, nil, stylemap[line.StyleIds[x]])
//visualX++
//}

//if len(line.Cursors) != 0 {
//// TODO: Verify if xi-core will take care of tabs for us
//cX := GetCursorVisualX(line.Cursors[0], line.Text)
//// TODO: Multiple cursor support
//v.Editor.screen.ShowCursor(cX, y)
//}
//}
//}
//}

func (t *TextModel) GetCell(x int, y int) (rune, tcell.Style, []rune, int) {
	lines := t.lc.lines
	if len(lines) == 0 {
		return ' ', 0, nil, 1
	}

	if y < len(lines) {
		line := lines[y]
		if x < len(line.Text) {
			return rune(line.Text[x]), styles[line.StyleIds[x]], nil, 1
		}
	}

	return ' ', 0, nil, 1
}

func (t *TextModel) GetBounds() (int, int) {
	return 20, 20
	//panic("not implemented")
}

func (t *TextModel) SetCursor(int, int) {
	//panic("not implemented")
}

func (t *TextModel) GetCursor() (int, int, bool, bool) {
	return 0, 0, true, true
	//panic("not implemented")
}

func (t *TextModel) MoveCursor(offx int, offy int) {
	//panic("not implemented")
}
