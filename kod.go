package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

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
	p.Stderr = os.Stderr
	stdout, _ := p.StdoutPipe()
	stdin, _ := p.StdinPipe()
	p.Start()

	rw := readwriter{stdout, stdin}
	c := rpc.NewConnection(rw, func(msg *rpc.Message) {
		if msg.Method == "update" {
			update := msg.Params["update"].(map[string]interface{})
			ops := update["ops"].([]interface{})

			for _, op := range ops {
				opMap := op.(map[string]interface{})
				opType := opMap["op"].(string)

				switch opType {
				case "invalidate":
					fmt.Println("invalidate")
				case "ins":
					fmt.Println("insert")
				default:
					fmt.Println("other: " + opType)
				}
			}
		}
	})

	c.Send("new_view", &rpc.Params{"file_path": "README.md"})

	p.Wait()
}
