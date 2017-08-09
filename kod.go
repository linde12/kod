package main

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/linde12/kod/editor"
	"github.com/linde12/kod/mode"
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
	p.Stderr = os.Stderr
	stdout, _ := p.StdoutPipe()
	stdin, _ := p.StdinPipe()
	p.Start()

	rw := readwriter{stdout, stdin}
	e := editor.NewEditor(rw)
	e.SetMode(mode.NewNormalMode(e))
	e.Start()
}
