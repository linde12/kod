package editor

import (
	"github.com/gdamore/tcell"
)

type View struct {
	*Buffer
	Editor *Editor
	Width  int
	Height int

	// Topmost line in the text buffer, used for vertical scrolling
	Topline int
}

func NewView(e *Editor, buf *Buffer) *View {
	// fullscreen view
	w, h := e.screen.Size()
	return &View{
		Buffer:  buf,
		Editor:  e,
		Width:   w,
		Height:  h,
		Topline: 0,
	}
}

func (v *View) Draw() {
	for y, line := range v.lines[v.Topline:] {
		visualX := 0
		for _, char := range []rune(string(line.data)) {
			// TODO: Highlight line and use line.Style as style
			if char == '\t' {
				ts := tabSize - (visualX % tabSize)
				for i := 0; i < ts; i++ {
					v.Editor.screen.SetCell(visualX+i, y, v.Editor.defaultStyle, ' ')
				}
				visualX += ts
			} else {
				v.Editor.screen.SetCell(visualX, y, v.Editor.defaultStyle, char)
				visualX++
			}
		}
	}
	v.Editor.screen.ShowCursor(v.Cursor.GetVisualX(), v.Cursor.Y-v.Topline)
}

func (v *View) Relocate() {
	y := v.Cursor.Y
	if y > v.Topline+v.Height-1 {
		v.Topline++
	} else if y < v.Topline {
		v.Topline--
	}
}

func (v *View) HandleEvent(ev tcell.Event) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		v.Editor.Mode.OnKey(e)
		//if e.Key() == tcell.KeyRune {
		//v.Insert([]byte(string(e.Rune())), v.CursorPos())
		//v.CursorRight()
		//} else {
		//switch e.Key() {
		//case tcell.KeyBackspace2, tcell.KeyBackspace:
		//line := v.CurLine()
		//if len(line.data) > 0 {
		//if v.Cursor.X > 0 {
		//v.RemoveRune(Pos{v.Cursor.X - 1, v.Cursor.Y})
		//v.CursorLeft()
		//} else if v.Cursor.Y != 0 {
		//lineAbove := v.lines[v.Cursor.Y-1]
		//end := len(lineAbove.data)
		//v.JoinLines(v.Cursor.Y-1, v.Cursor.Y)
		//v.CursorUp()
		//v.Cursor.X = end
		//}
		//} else if v.Cursor.Y != 0 {
		//v.RemoveLine(v.Cursor.Y)
		//v.CursorUp()
		//v.CursorEnd()
		//}
		//case tcell.KeyTAB:
		//v.Insert([]byte("\t"), v.CursorPos())
		//v.CursorRight()
		//case tcell.KeyEnter:
		//v.Insert([]byte("\n"), v.CursorPos())
		//v.CursorDown()
		//v.CursorBegin()
		//case tcell.KeyLeft:
		//v.CursorLeft()
		//case tcell.KeyUp:
		//v.CursorUp()
		//case tcell.KeyRight:
		//v.CursorRight()
		//case tcell.KeyDown:
		//v.CursorDown()
		//}
		//}

		// relocate view
		v.Relocate()
	}
}
