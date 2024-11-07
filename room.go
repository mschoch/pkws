package pkws

import (
	"log/slog"
)

// Room is the unit of group communication
// Participants in the room can communicate via a websocket
// as well as invoke HTTP methods on a common server
type Room struct {
	namespace string
	name      string
	storage   Storage

	connections map[string]*Connection
	s           Server
	log         *slog.Logger
}

func NewRoom(namespace string, name string, log *slog.Logger) *Room {
	return &Room{
		namespace:   namespace,
		name:        name,
		connections: make(map[string]*Connection),
		storage:     NewEphemeralStorage(),
		log:         log,
	}
}

func (r *Room) AttachServer(s Server) {
	single := NewSingle(s)
	r.s = single
}

func (r *Room) Add(c *Connection) {
	r.connections[c.id] = c
}

func (r *Room) Remove(c *Connection) {
	delete(r.connections, c.id)
}

func (r *Room) Storage() Storage {
	return r.storage
}

func (r *Room) Broadcast(msg []byte) {
	r.BroadcastExcept(msg, nil)
}

func (r *Room) BroadcastExcept(msg []byte, except []string) {
OUTER:
	for _, conn := range r.connections {
		for _, ex := range except {
			if conn.id == ex {
				continue OUTER
			}
		}
		conn.Send(msg)
	}
}
