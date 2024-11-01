package pkws

import (
	"log/slog"
	"net/http"
)

// Logger is an implementation of the Server interface
// which simply logs each method invocation
type Logger struct {
	log *slog.Logger
}

func NewLogger(room *Room, log *slog.Logger) Server {
	return &Logger{
		log: log,
	}
}

func (l *Logger) OnStart() {
	l.log.Info("started")
}

func (l *Logger) OnConnect(conn *Connection) {
	l.log.Info("connected", "id", conn.ID())
}

func (l *Logger) OnMessage(msg []byte, sender *Connection) {
	l.log.Info("received", "message", string(msg), "from", sender.ID())
}

func (l *Logger) OnClose(conn *Connection) {
	l.log.Info("closed", "id", conn.ID())
}

func (l *Logger) OnError(conn *Connection, err error) {
	l.log.Error("error", "err", err)
}

func (l *Logger) OnRequest(w http.ResponseWriter, r *http.Request) {
	l.log.Info("request", "url", r.URL.String(), "method", r.Method)
	w.WriteHeader(http.StatusNoContent)
}

func (l *Logger) OnAlarm() {
	l.log.Info("alarm")
}
