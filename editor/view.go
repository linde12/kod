package editor

import (
	"strconv"

	"github.com/gdamore/tcell"
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
	*LineCache
	*InputHandler

	ID        string
	view      *Viewport
	gutter    *Viewport
	xi        *rpc.Connection
	ViewID    string
	lineStart int
	lineEnd   int
}

var tmpStyle tcell.Style

func ralign(str string, width int) string {
	diff := width - len(str)
	res := ""
	for i := 0; i < diff; i++ {
		res += " "
	}

	res += str
	return res
}

func NewView(path string, vp *Viewport, xi *rpc.Connection) (*View, error) {
	view := &View{}
	view.view = NewViewport(vp, 3, 0)
	view.gutter = NewViewport(vp, 0, 0)
	//view.editview = NewViewport(view.view, 3, 0)
	view.xi = xi
	view.LineCache = NewLineCache()

	msg, err := xi.Request(&rpc.Request{
		Method: "new_view",
		Params: &rpc.Object{"file_path": path},
	})
	if err != nil {
		return view, err
	}

	view.ID = msg.Value.(string)
	view.InputHandler = &InputHandler{view.ID, path, xi}

	// Set scroll window size
	_, height := vp.Size()
	xi.Notify(&rpc.Request{
		Method: "edit",
		Params: &rpc.Object{
			"method":  "scroll",
			"params":  &rpc.Array{0, height - 2},
			"view_id": view.ID,
		},
	})

	return view, nil
}

func (v *View) Draw() {
	if len(v.lines) == 0 {
		return
	}

	// render line numbers
	// TODO: improve, alot... :-)
	style := defaultStyle.Foreground(tcell.ColorBlue)
	width := len(strconv.Itoa(len(v.lines) + v.LineCache.invalidBefore))

	v.gutter.SetWidth(width + 1)
	v.view.SetOffsetX(width + 2)
	for i := 0; i < len(v.lines); i++ {
		nLine := i + v.invalidBefore
		txt := ralign(strconv.Itoa(nLine), width)
		width := len(txt)
		for x := 0; x < width; x++ {
			v.gutter.SetContent(1+x, nLine, rune(txt[x]), nil, style)
		}
	}

	// TODO: Line numbers in a separate viewport
	for y, line := range v.lines {
		if line == nil {
			continue
		}
		nLine := y + v.invalidBefore
		visualX := 0
		for x, char := range []rune(line.Text) {
			// TODO: Do this somewhere else
			var style = defaultStyle
			// TODO: Reserved??
			if line.StyleIds[x] >= 2 {
				fg, _, _ := styles[line.StyleIds[x]].Decompose()
				style = style.Foreground(fg)
			}

			if char == '\t' {
				ts := tabSize - (visualX % tabSize)
				for i := 0; i < ts; i++ {
					v.view.SetContent(visualX+i, nLine, ' ', nil, style)
				}
				visualX += ts
			} else if char != '\n' {
				// TODO: Trim newline in a better way?
				v.view.SetContent(visualX, nLine, char, nil, style)
				visualX++
			}
		}
		if len(line.Cursors) != 0 {
			// TODO: Verify if xi-core will take care of tabs for us
			cX := GetCursorVisualX(line.Cursors[0], line.Text)
			// TODO: Multiple cursor support
			v.view.ShowCursor(cX, nLine)
		}
	}
	lineStart, lineEnd := v.view.GetViewport()
	if lineStart != v.lineStart || lineEnd != v.lineEnd {
		v.Scroll(lineStart, lineEnd)
		v.lineStart = lineStart
		v.lineEnd = lineEnd
	}
}

func (v *View) MakeLineVisible(line int) {
	v.view.MakeVisibleY(line)
	v.gutter.MakeVisibleY(line)
}

func (v *View) HandleEvent(ev tcell.Event) {
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
					v.Tab()
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
}
