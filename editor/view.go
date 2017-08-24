package editor

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/linde12/kod/rpc"
)

const tabSize = 4

// TODO: Move this to RPC
type requestLines struct {
	Method string `json:"method"`
	Params []int  `json:"params"`
	ViewID string `json:"view_id"`
}

type View struct {
	views.View
	views.WidgetWatchers

	*LineCache
	*InputHandler
	ID     string
	Editor *Editor
	ViewID string
	Width  int
	Height int

	// Topmost line in the text buffer, used for vertical scrolling
	Topline int
}

var tmpStyle tcell.Style

func NewView(path string, e *Editor) (*View, error) {
	view := &View{}
	view.Editor = e
	view.LineCache = NewLineCache()

	msg, err := e.rpc.Request(&rpc.Request{
		Method: "new_view",
		Params: &rpc.Object{"file_path": path},
	})
	if err != nil {
		return view, err
	}

	view.ID = msg.Value.(string)
	view.InputHandler = &InputHandler{view.ID, path, e.rpc}

	return view, nil
}

func (v *View) Draw() {
	if len(v.lines) == 0 {
		return
	}

	// TODO: Line numbers
	// TODO: Fix choppy scrolling
	for y, line := range v.lines {
		visualX := 0
		for x, char := range []rune(line.Text) {
			// TODO: Do this somewhere else
			var style tcell.Style = defaultStyle
			// TODO: Reserved??
			if line.StyleIds[x] >= 2 {
				fg, _, _ := v.Editor.styleMap.Get(line.StyleIds[x]).Decompose()
				style = defaultStyle.Foreground(fg)
			}

			if char == '\t' {
				ts := tabSize - (visualX % tabSize)
				for i := 0; i < ts; i++ {
					v.SetContent(visualX+i, y, ' ', nil, style)
				}
				visualX += ts
			} else if char != '\n' {
				// TODO: Trim newline in a better way?
				v.SetContent(visualX, y, char, nil, style)
				visualX++
			}

			if len(line.Cursors) != 0 {
				// TODO: Verify if xi-core will take care of tabs for us
				cX := GetCursorVisualX(line.Cursors[0], line.Text)
				// TODO: Multiple cursor support
				v.SetContent(cX, y, char, nil, style.Reverse(true))
				//v.ShowCursor(cX, y)
			}
		}
	}
}

func (v *View) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		ctrl := e.Modifiers()&tcell.ModCtrl != 0

		if e.Key() == tcell.KeyRune && !ctrl {
			v.Insert(string(e.Rune()))
		} else {
			if !ctrl {
				switch e.Key() {
				case tcell.KeyBackspace2, tcell.KeyBackspace:
					v.DeleteBackward()
				case tcell.KeyTAB:
					// TODO: Use v.Tab() when it's ready
					v.Insert("\t")
				case tcell.KeyEnter:
					v.Newline()
				case tcell.KeyLeft:
					v.MoveLeft()
				case tcell.KeyUp:
					v.MoveUp()
				case tcell.KeyRight:
					v.MoveRight()
				case tcell.KeyDown:
					v.MoveDown()
				case tcell.KeyDelete:
					v.DeleteForward()
				}
			} else {
				// Ctrl
				switch e.Key() {
				case tcell.KeyLeft:
					v.MoveWordLeft()
				case tcell.KeyRight:
					v.MoveWordRight()
				case tcell.KeyCtrlQ:
					v.Editor.CloseView(v)
				case tcell.KeyCtrlS:
					v.Save()
				case tcell.KeyCtrlU:
					v.Undo()
				case tcell.KeyCtrlR:
					v.Redo()
				}
			}
		}
	}
	return true
}

func (v *View) Resize() {
	v.Width, v.Height = v.Size()
	// Set scroll window size
	v.Editor.rpc.Notify(&rpc.Request{
		Method: "edit",
		Params: &rpc.Object{
			"method":  "scroll",
			"params":  &rpc.Array{0, v.Height - 2},
			"view_id": v.ID,
		},
	})
	v.PostEventWidgetResize(v)
}

func (v *View) SetView(view views.View) {
	v.View = view
}
