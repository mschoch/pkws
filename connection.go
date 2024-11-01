package pkws

import "github.com/gorilla/websocket"

// Connection represents a websocket connection to a room
type Connection struct {
	id     string
	ws     *websocket.Conn
	send   chan []byte
	parent *Room
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) Send(msg []byte) {
	c.send <- msg
}

func (c *Connection) Writer() {
	for msg := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func (c *Connection) Reader() {
	for {
		if _, msg, err := c.ws.ReadMessage(); err == nil {
			c.parent.s.OnMessage(msg, c)
		} else {
			break
		}
	}
	c.ws.Close()
}
