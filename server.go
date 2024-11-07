package pkws

import (
	"log/slog"
	"net/http"
)

type Env any

// ServerConstructor is the function signature used to register
// Server implementations for specific parties
type ServerConstructor[Env any] func(env Env, room *Room, log *slog.Logger) Server

// Server is the primary interface by which one implements the logic
// backing the behavior of the room.
type Server interface {
	OnStart()
	OnConnect(conn *Connection)
	OnMessage(msg []byte, sender *Connection)
	OnClose(conn *Connection)
	OnError(conn *Connection, err error)
	OnRequest(w http.ResponseWriter, r *http.Request)
	OnAlarm()
}
