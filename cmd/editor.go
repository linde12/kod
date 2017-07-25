package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
)

type Editor struct {
	screen tcell.Screen
	Views  []*View

	defaultStyle tcell.Style

	events chan tcell.Event
}

func (e *Editor) CurView() *View {
	return e.Views[0]
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

// TODO: Add support for multiple files & stdin pipe
func (e *Editor) loadInput() (b *Buffer) {
	if len(os.Args) > 1 {
		// TODO: Check if file exists, if it's really a file, error handling
		f, _ := os.Open(os.Args[1])
		defer f.Close()
		b = NewBuffer(f, os.Args[1])
	}
	return b
}

func NewEditor() *Editor {
	e := &Editor{}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	e.initScreen()
	defer e.screen.Fini()

	// screen event channel
	e.events = make(chan tcell.Event, 100)
	quit := make(chan bool, 1)

	go func() {
		for {
			if e.screen != nil {
				// feed events into channel
				e.events <- e.screen.PollEvent()
			}
		}
	}()

	buf := e.loadInput()
	e.Views = append(e.Views, NewView(e, buf))

	// main loop
	for {
		e.screen.Clear()
		e.CurView().Draw()
		e.screen.Show()

		var event tcell.Event
		select {
		case event = <-e.events:
		case <-quit:
			e.screen.Fini()
			log.Println("bye")
			os.Exit(0)
		}

		for event != nil {
			switch ev := event.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
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
