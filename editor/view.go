package editor

import (
	"os"

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
	views.BoxLayout
	*LineCache
	*InputHandler

	TextView   *views.CellView
	GutterView *views.CellView

	ID       string
	ViewID   string
	FilePath string
	Width    int
	Height   int
}

var tmpStyle tcell.Style

func NewView(path string, xi *rpc.Connection) (*View, error) {
	view := &View{}
	view.FilePath = path
	view.LineCache = NewLineCache()
	view.TextView = views.NewCellView()
	view.TextView.SetModel(&TextModel{lc: view.LineCache})

	view.SetOrientation(views.Horizontal)
	view.AddWidget(view.TextView, 1)

	msg, err := xi.Request(&rpc.Request{
		Method: "new_view",
		Params: &rpc.Object{"file_path": path},
	})
	if err != nil {
		return view, err
	}

	view.ID = msg.Value.(string)
	view.InputHandler = &InputHandler{ViewID: view.ID, xi: xi}

	return view, nil
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
					os.Exit(0)
				case tcell.KeyCtrlS:
					v.Save(v.FilePath)
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
