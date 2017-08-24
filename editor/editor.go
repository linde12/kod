package editor

import (
	"io"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/linde12/kod/app"
	"github.com/linde12/kod/rpc"
)

type Params map[string]interface{}

type Editor struct {
	views.Panel

	topbar *views.SimpleStyledTextBar
	status *views.SimpleStyledTextBar

	app *app.Application

	Views     map[string]*View
	curViewID string
	rpc       *rpc.Connection

	styleMap *StyleMap

	// ui events
	events chan tcell.Event
	// user events
	RedrawEvents chan struct{}
}

func (e *Editor) CurView() *View {
	return e.Views[e.curViewID]
}

func (e *Editor) HandleEvent(ev tcell.Event) bool {
	switch kev := ev.(type) {
	case *tcell.EventKey:
		if kev.Key() == tcell.KeyCtrlQ {
			e.app.Quit()
			return true
		}
		return e.CurView().HandleEvent(ev)
	}
	return e.Panel.HandleEvent(ev)
}

func NewEditor(rw io.ReadWriter, app *app.Application) *Editor {
	e := &Editor{}
	e.app = app

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	// screen event channel
	e.events = make(chan tcell.Event, 50)
	e.RedrawEvents = make(chan struct{}, 50)

	e.styleMap = NewStyleMap()
	e.Views = make(map[string]*View)

	e.rpc = rpc.NewConnection(rw)

	// Set theme, this might be removed when xi-editor has a config file
	e.rpc.Notify(&rpc.Request{
		Method: "set_theme",
		// TODO: Ability to change this would be nice...
		// Try "InspiredGitHub" or "Solarized (dark)"
		Params: rpc.Object{"theme_name": "base16-eighties.dark"},
	})

	return e
}

func (e *Editor) CloseView(v *View) {
	delete(e.Views, v.ID)
}

func (e *Editor) handleRequests() {
	for {
		msg := <-e.rpc.Messages

		switch msg.Value.(type) {
		case *rpc.Update:
			update := msg.Value.(*rpc.Update)
			view := e.Views[update.ViewID]
			view.ApplyUpdate(msg.Value.(*rpc.Update))
		case *rpc.DefineStyle:
			e.styleMap.DefineStyle(msg.Value.(*rpc.DefineStyle))
			e.SetStyle(defaultStyle)
		case *rpc.ThemeChanged:
			themeChanged := msg.Value.(*rpc.ThemeChanged)
			theme := themeChanged.Theme

			bg := tcell.NewRGBColor(theme.Bg.ToRGB())
			defaultStyle = defaultStyle.Background(bg)
			fg := tcell.NewRGBColor(theme.Fg.ToRGB())
			defaultStyle = defaultStyle.Foreground(fg)

			e.SetStyle(defaultStyle)

			log.Printf("Theme:%v", theme)
		}

		v := e.CurView()
		// TODO: Research if this is actually needed
		v.PostEventWidgetContent(v)
		e.app.Update()
	}
}

func (e *Editor) Start() {
	path := os.Args[1]
	view, _ := NewView(path, e)
	e.Views[view.ID] = view
	e.curViewID = view.ID

	e.topbar = views.NewSimpleStyledTextBar()
	e.topbar.SetCenter("kod")

	e.status = views.NewSimpleStyledTextBar()
	e.status.SetLeft(e.CurView().FilePath)

	e.Panel.SetTitle(e.topbar)
	e.Panel.SetContent(view)
	e.Panel.SetStatus(e.status)

	go e.handleRequests()
}
