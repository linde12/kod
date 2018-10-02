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

// Object represents a JSON object.
type Object map[string]interface{}

// Array represents a JSON array.
type Array []interface{}

type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ViewID string      `json:"view_id,omitempty"`
}

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
	*Request
	ID int `json:"id,omitempty"`
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
			switch msg.Method {
			case "update":
				var update Update
				json.Unmarshal(msg.Params, &update)
				c.Messages <- &Message{msg.Method, &update}
			case "def_style":
				var defStyle DefineStyle
				json.Unmarshal(msg.Params, &defStyle)
				c.Messages <- &Message{msg.Method, &defStyle}
			case "theme_changed":
				var themeChanged ThemeChanged
				json.Unmarshal(msg.Params, &themeChanged)
				c.Messages <- &Message{msg.Method, &themeChanged}
			case "scroll_to":
				var scrollTo ScrollTo
				json.Unmarshal(msg.Params, &scrollTo)
				c.Messages <- &Message{msg.Method, &scrollTo}
			default:
				log.Println("unhandled request: " + msg.Method)
			}
		}
	}

	if in.Err() != nil {
		log.Printf("error: %s", in.Err())
	}
}

// TODO: notify function
func (c *Connection) send(msg *outgoingMessage) int {
	b, _ := json.Marshal(msg)
	log.Printf(">>> %s\n", b)
	c.rw.Write(b)
	c.rw.Write([]byte{LF})
	return c.rpcIndex
}

func (c *Connection) Request(r *Request) (*Message, error) {
	ch := make(chan *Message, 1)
	id := c.RequestAsync(r, func(m *Message) {
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

func (c *Connection) RequestAsync(r *Request, callback func(*Message)) int {
	c.rpcIndex++
	c.pending[c.rpcIndex] = callback

	c.send(&outgoingMessage{
		ID:      c.rpcIndex,
		Request: r,
	})
	return c.rpcIndex
}

func (c *Connection) Notify(r *Request) {
	c.send(&outgoingMessage{
		Request: r,
	})
}
