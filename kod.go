package main

import (
	"log"
	"os"

	"github.com/linde12/kod/editor"
	"github.com/linde12/kod/mode"
)

func main() {
	f, _ := os.OpenFile("out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)
	e := editor.NewEditor()
	e.SetMode(mode.NewNormalMode(e))
	e.Start()
}
