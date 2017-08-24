package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/app"
	"github.com/linde12/kod/editor"
)

type readwriter struct {
	io.Reader
	io.Writer
}

func main() {
	f, _ := os.OpenFile("out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)

	p := exec.Command("./xi-core")
	stdout, _ := p.StdoutPipe()
	stdin, _ := p.StdinPipe()
	p.Start()

	rw := readwriter{stdout, stdin}

	// TODO: Handle error
	screen, _ := tcell.NewScreen()

	app := app.NewApplication(screen)

	e := editor.NewEditor(rw, app)
	e.Start()

	// TODO: Handle error
	app.SetRootWidget(e)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}
}
