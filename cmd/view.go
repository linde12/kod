package main

import (
	"log"

	"github.com/gdamore/tcell"
)

type View struct {
	buf    *Buffer
	Width  int
	Height int

	// Topmost line in the text buffer, used for vertical scrolling
	Topline int
}

func NewView(buf *Buffer) *View {
	// fullscreen view
	w, h := screen.Size()
	return &View{
		buf:     buf,
		Width:   w,
		Height:  h,
		Topline: 0,
	}
}

func (v *View) Draw() {
	log.Printf("data:% x\n", v.buf.CurLine().data)
	for y, line := range v.buf.lines[v.Topline:] {
		visualX := 0
		for _, char := range []rune(string(line.data)) {
			// TODO: Highlight line and use line.Style as style
			if char == '\t' {
				ts := tabSize - (visualX % tabSize)
				for i := 0; i < ts; i++ {
					screen.SetCell(visualX+i, y, defaultStyle, ' ')
				}
				visualX += ts
			} else {
				screen.SetCell(visualX, y, defaultStyle, char)
				visualX++
			}
		}
	}
	screen.ShowCursor(v.buf.Cursor.GetVisualX(), v.buf.Cursor.Y-v.Topline)
}

func (v *View) Relocate() {
	y := v.buf.Cursor.Y
	if y > v.Topline+v.Height-1 {
		v.Topline++
	} else if y < v.Topline {
		v.Topline--
	}
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
						v.buf.RemoveRune(Pos{v.buf.Cursor.X - 1, v.buf.Cursor.Y})
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
			case tcell.KeyTAB:
				v.buf.Insert([]byte("\t"), v.buf.CursorPos())
				v.buf.CursorRight()
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

		// relocate view
		v.Relocate()
	}
}
