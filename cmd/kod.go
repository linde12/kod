package main

import (
	"log"
	"os"
)

func main() {
	f, _ := os.OpenFile("out.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)
	NewEditor()
}
