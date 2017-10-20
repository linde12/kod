package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/linde12/kod/editor"
	"github.com/linde12/kod/rpc"
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
	xi := rpc.NewConnection(rw)
	e := editor.NewEditor(xi)
	e.Start()

	fmt.Println("\nInterrupted.")
}
