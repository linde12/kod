package main

import (
	"github.com/gdamore/tcell"
)

type View struct {
	buf *Buffer
}

func NewView(buf *Buffer) *View {
	return &View{
		buf: buf,
	}
}

func (v *View) Draw() {
	for y, line := range v.buf.lines {
		for x, char := range []rune(string(line.data)) {
			// TODO: Highlight line and use line.Style as style
			screen.SetCell(x, y, defaultStyle, rune(char))
		}
	}
	screen.ShowCursor(v.buf.Cursor.x, v.buf.Cursor.y)
}

func (v *View) HandleEvent(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyRune {
			line := v.buf.CurLine()
			line.Insert([]byte(string(e.Rune())), v.buf.Cursor.x)
			v.buf.CursorRight()
		} else {
			switch e.Key() {
			case tcell.KeyBackspace2, tcell.KeyBackspace:
				line := v.buf.CurLine()
				if len(line.data) > 0 {
					if v.buf.Cursor.x > 0 {
						line.RemoveRune(v.buf.Cursor.x - 1)
						v.buf.CursorLeft()
					} else if v.buf.Cursor.y != 0 {
						lineAbove := v.buf.lines[v.buf.Cursor.y-1]
						end := len(lineAbove.data)
						v.buf.JoinLines(v.buf.Cursor.y-1, v.buf.Cursor.y)
						v.buf.CursorUp()
						v.buf.Cursor.SetX(end)
					}
				} else if v.buf.Cursor.y != 0 {
					v.buf.RemoveLine(v.buf.Cursor.y)
					v.buf.CursorUp()
					v.buf.CursorEnd()
				}
			case tcell.KeyEnter:
				v.buf.Split(v.buf.Cursor.x, v.buf.Cursor.y)
				v.buf.CursorDown()
				v.buf.CursorBegin()
			case tcell.KeyLeft:
				v.buf.CursorLeft()
			case tcell.KeyUp:
				v.buf.CursorUp()
			case tcell.KeyRight:
				v.buf.CursorRight()
			case tcell.KeyDown:
				v.buf.CursorDown()
			}
		}
	}
}
