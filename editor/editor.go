package editor

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/rpc"
)

type Params map[string]interface{}

type Command interface {
	Apply(e *Editor)
}

type Mode interface {
	OnKey(ev *tcell.EventKey)
}

type Editor struct {
	screen tcell.Screen
	Views  []*View
	Mode   Mode
	rpc    *rpc.Connection

	defaultStyle tcell.Style

	// ui events
	events chan tcell.Event
	// user events
	Commands chan Command
}

func (e *Editor) SetMode(m Mode) {
	e.Mode = m
}

func (e *Editor) CurView() *View {
	if len(e.Views) > 0 {
		return e.Views[0]
	}
	return nil
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
		e.CurView().HandleEvent(ev)
	}
}

func NewEditor(rw io.ReadWriter) *Editor {
	e := &Editor{}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	// screen event channel
	e.events = make(chan tcell.Event, 50)
	e.Commands = make(chan Command, 50)

	e.rpc = rpc.NewConnection(rw)
	return e
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
	_, err := e.rpc.Send("new_view", &rpc.Params{"file_path": path})
	if err != nil {
		log.Println(err)
	}
	//buf := NewBuffer(strings.NewReader(res.Params["text"].(string)), path)
	//e.Views = append(e.Views, NewView(e, buf))

	// editor loop
	for {
		e.screen.Clear()
		if e.CurView() != nil {
			e.CurView().Draw()
		}
		e.screen.Show()

		var event tcell.Event
		select {
		case event = <-e.events:
		case cmd := <-e.Commands:
			cmd.Apply(e)
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
