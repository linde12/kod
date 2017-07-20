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
	screen.ShowCursor(v.buf.Cursor.X, v.buf.Cursor.Y)
}

func (v *View) HandleEvent(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyRune {
			v.buf.Insert([]byte(string(e.Rune())), v.buf.CursorPos())
			v.buf.CursorRight()
		} else {
			switch e.Key() {
			case tcell.KeyBackspace2, tcell.KeyBackspace:
				line := v.buf.CurLine()
				if len(line.data) > 0 {
					if v.buf.Cursor.X > 0 {
						line.RemoveRune(v.buf.Cursor.X - 1)
						v.buf.CursorLeft()
					} else if v.buf.Cursor.Y != 0 {
						lineAbove := v.buf.lines[v.buf.Cursor.Y-1]
						end := len(lineAbove.data)
						v.buf.JoinLines(v.buf.Cursor.Y-1, v.buf.Cursor.Y)
						v.buf.CursorUp()
						v.buf.Cursor.X = end
					}
				} else if v.buf.Cursor.Y != 0 {
					v.buf.RemoveLine(v.buf.Cursor.Y)
					v.buf.CursorUp()
					v.buf.CursorEnd()
				}
			case tcell.KeyEnter:
				v.buf.Insert([]byte("\n"), v.buf.CursorPos())
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
