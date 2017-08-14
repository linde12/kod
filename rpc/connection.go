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

// Params is the JSON field `params` in the RPC message.
type Params map[string]interface{}

// incomingMessage is a representation of an incoming RPC message.
type incomingMessage struct {
	Id     int             `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Error  string          `json:"error,omitempty"`
	Result interface{}     `json:"result,omitempty"`
}

// outgoingMessage is a representation of an outgoing RPC message.
type outgoingMessage struct {
	Id     int     `json:"id"`
	Method string  `json:"method"`
	Params *Params `json:"params"`
}

// Message is a deserialized message which will be passed to the `Messages` channel.
type Message struct {
	Method string
	Value  interface{}
}

// Connection represents the connection to the backend.
// The underlying protocol doesn't matter as long as it is writeable and readable.
type Connection struct {
	Messages chan *Message
	rw       io.ReadWriter
	rpcIndex int
	pending  map[int]func(*Message)
}

func NewConnection(rw io.ReadWriter) *Connection {
	c := &Connection{
		Messages: make(chan *Message, 1),
		rw:       rw,
		pending:  make(map[int]func(*Message)),
	}

	go c.recv()

	return c
}

func (c *Connection) recv() {
	in := bufio.NewScanner(c.rw)

	for in.Scan() {
		log.Printf("<<< %s\n", in.Text())
		var msg incomingMessage
		json.Unmarshal([]byte(in.Text()), &msg)

		if msg.Id != 0 {
			if msg.Result != nil {
				// response
				if fn, ok := c.pending[msg.Id]; ok {
					fn(&Message{msg.Method, msg.Result})
				} else {
					log.Println("unhandled response: ", msg)
				}
			}
		} else {
			// request
			if msg.Method == "update" {
				var update Update
				json.Unmarshal(msg.Params, &update)
				c.Messages <- &Message{msg.Method, &update}
			} else {
				log.Println("unhandled request: " + msg.Method)
			}
		}
	}

	if in.Err() != nil {
		log.Printf("error: %s", in.Err())
	}
}

func (c *Connection) send(id int, method string, params *Params) int {
	msg := outgoingMessage{
		Id:     c.rpcIndex,
		Method: method,
		Params: params,
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
	c.rpcIndex++
	c.pending[c.rpcIndex] = callback

	c.send(c.rpcIndex, method, params)
	return c.rpcIndex
}
