package editor

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/rpc"
)

type Params map[string]interface{}

type Editor struct {
	screen    tcell.Screen
	Views     map[string]*View
	curViewID string
	rpc       *rpc.Connection

	defaultStyle tcell.Style

	// ui events
	events chan tcell.Event
	// user events
	RedrawEvents chan struct{}
}

func (e *Editor) CurView() (*View, error) {
	return e.ViewByID(e.curViewID)
}

func (e *Editor) initScreen() {
	var err error

	e.screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = e.screen.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	e.screen.SetStyle(e.defaultStyle)

	e.screen.Clear()
}

func (e *Editor) handleEvent(ev tcell.Event) {
	switch ev.(type) {
	case *tcell.EventKey:
		// TODO: Check if normal mode, if so check for
		// "global" keybindings which aren't bound to the buffer
		// and pass on buffer-specific keybindings
		v, err := e.CurView()
		if err != nil {
			log.Printf("can't find view: %s", err)
		}

		v.HandleEvent(ev)
	}
}

func NewEditor(rw io.ReadWriter) *Editor {
	e := &Editor{}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	// screen event channel
	e.events = make(chan tcell.Event, 50)
	e.RedrawEvents = make(chan struct{}, 50)

	e.Views = make(map[string]*View)

	e.rpc = rpc.NewConnection(rw)

	return e
}

func (e *Editor) ViewByID(viewID string) (*View, error) {
	view, ok := e.Views[viewID]
	if ok {
		return view, nil
	} else {
		return nil, errors.New("view not found:" + viewID)
	}
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

			if view, err := e.ViewByID(update.ViewID); err == nil {
				view.ApplyUpdate(msg)
				// TODO: Better way to signal redraw?
				e.RedrawEvents <- struct{}{}
			} else {
				log.Printf("can't update view: %s", err)
			}
		}
	}
}

func (e *Editor) Start() {
	e.initScreen()
	defer e.screen.Fini()

	quit := make(chan bool, 1)

	go func() {
		for {
			if e.screen != nil {
				// feed events into channel
				e.events <- e.screen.PollEvent()
			}
		}
	}()

	path := os.Args[1]
	view, _ := NewView(path, e)
	e.Views[view.ID] = view
	e.curViewID = view.ID

	go e.handleRequests()

	// editor loop
	for {
		if len(e.Views) != 0 {
			curView, err := e.CurView()
			if err != nil {
				log.Printf("can't find view: %s", err)
			}
			e.screen.Clear()
			curView.Draw()
			e.screen.Show()
		} else {
			quit <- true
		}

		var event tcell.Event
		select {
		case event = <-e.events:
		case <-e.RedrawEvents:
		case <-quit:
			e.screen.Fini()
			log.Println("bye")
			os.Exit(0)
		}

		for event != nil {
			switch ev := event.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyF1:
					close(quit)
				}
			case *tcell.EventResize:
				e.screen.Sync()
			}

			e.handleEvent(event)

			// continue handling events
			select {
			case event = <-e.events:
			default:
				event = nil
			}
		}
	}
}
