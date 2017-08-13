package rpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

const LF = 0xA // Line-feed

type RequestHandler func(*Message)

// Params is the JSON field `params` in the RPC message
type Params map[string]interface{}

// Message is a representation of a RPC message
type Message struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params Params      `json:"params"`
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

// Connection represents the connection to the backend.
// The underlying protocol doesn't matter as long as it is writeable and readable.
type Connection struct {
	rw            io.ReadWriter
	handleRequest RequestHandler
	rpcIndex      int
	pending       map[int]func(*Message)
}

func NewConnection(rw io.ReadWriter, rh RequestHandler) *Connection {
	c := &Connection{
		rw:            rw,
		handleRequest: rh,
		pending:       make(map[int]func(*Message)),
	}

	go c.recv()

	return c
}

func (c *Connection) recv() {
	in := bufio.NewScanner(c.rw)

	for in.Scan() {
		log.Printf("<<< %s\n", in.Text())
		var msg Message
		json.Unmarshal([]byte(in.Text()), &msg)

		if msg.Id != 0 {
			if msg.Result != nil {
				// response
				if fn, ok := c.pending[msg.Id]; ok {
					fn(&msg)
				}
			} else {
				// request
				c.handleRequest(&msg)
				// TODO: Handle requests
			}
		} else {
			// notification
			c.handleRequest(&msg)
			// TODO: Handle notifications
		}
	}

	if in.Err() != nil {
		log.Printf("error: %s", in.Err())
	}
}

func (c *Connection) send(method string, params *Params) int {
	msg := Message{
		Id:     c.rpcIndex,
		Method: method,
		Params: *params,
	}
	b, _ := json.Marshal(&msg)
	log.Printf(">>> %s\n", b)
	c.rw.Write(b)
	c.rw.Write([]byte{LF})
	return c.rpcIndex
}

func (c *Connection) Send(method string, params *Params) (*Message, error) {
	ch := make(chan *Message, 1)
	id := c.SendAsync(method, params, func(m *Message) {
		ch <- m
	})

	select {
	case m := <-ch:
		return m, nil
		// TODO: const values
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("request %d timed out", id)
	}
}

func (c *Connection) SendAsync(method string, params *Params, callback func(*Message)) int {
	id := c.send(method, params)
	c.pending[id] = callback
	return id
}
