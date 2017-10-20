package editor

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gdamore/tcell/views"
	"github.com/linde12/kod/rpc"
)

type Update func(*Editor)

type Params map[string]interface{}

type Editor struct {
	views.Application

	Views     map[string]*View
	curViewID string
	xi        *rpc.Connection

	updates chan Update
	draws   chan Update // TODO: Not update
}

func (e *Editor) CurView() *View {
	return e.Views[e.curViewID]
}

func NewEditor(conn *rpc.Connection) *Editor {
	e := &Editor{
		Views:   make(map[string]*View),
		xi:      conn,
		updates: make(chan Update, 1),
	}

	// Set theme, this might be removed when xi-editor has a config file
	e.xi.Notify(&rpc.Request{
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
		msg := <-e.xi.Messages

		switch msg.Value.(type) {
		case *rpc.Update:
			update := msg.Value.(*rpc.Update)
			e.updates <- func(e *Editor) {
				view := e.Views[update.ViewID]
				view.ApplyUpdate(msg.Value.(*rpc.Update))
			}
		case *rpc.DefineStyle:
			styles.defineStyle(msg.Value.(*rpc.DefineStyle))
		case *rpc.ThemeChanged:
			themeChanged := msg.Value.(*rpc.ThemeChanged)
			theme := themeChanged.Theme
			e.updates <- func(e *Editor) {
				log.Println("THEME CHANGED :D")
			}
			//bg := tcell.NewRGBColor(theme.Bg.ToRGB())
			//defaultStyle = defaultStyle.Background(bg)
			//fg := tcell.NewRGBColor(theme.Fg.ToRGB())
			//defaultStyle = defaultStyle.Foreground(fg)

			//e.SetStyle(defaultStyle)

			log.Printf("Theme:%v", theme)
		}
	}
}

func (e *Editor) Start() {
	path := os.Args[1]
	view, _ := NewView(path, e.xi)
	e.Views[view.ID] = view
	e.curViewID = view.ID

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	panel := views.NewPanel()
	title := views.NewTextBar()
	title.SetCenter("kod", 55)
	panel.SetTitle(title)
	panel.SetContent(e.CurView())

	e.SetRootWidget(panel)

	go func() {
		if err := e.Run(); e != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}()

	go e.handleRequests()

	for {
		select {
		case update := <-e.updates:
			log.Println("APPL")
			e.PostFunc(func() {
				update(e)
				log.Println("DONE")
			})
			update(e)
		case <-e.draws:
			log.Println("DRAWWW")
			e.Update()
			//e.Draw()
		case <-sig:
			return
		}
	}
}
