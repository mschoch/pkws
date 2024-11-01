package pkws

import (
	"net/http"
	"sync"
)

// Single is an implementation of the Server interface
// which uses a mutex to enforce processing of a single
// operation at a time.
//
// To use it, wrap your actual server:
// single := NewServer(actualServer)
type Single struct {
	m sync.Mutex

	s Server
}

func NewSingle(s Server) *Single {
	return &Single{
		s: s,
	}
}

func (s *Single) OnStart() {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnStart()
}

func (s *Single) OnConnect(conn *Connection) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnConnect(conn)
}

func (s *Single) OnMessage(msg []byte, sender *Connection) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnMessage(msg, sender)
}

func (s *Single) OnClose(conn *Connection) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnClose(conn)
}

func (s *Single) OnError(conn *Connection, err error) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnError(conn, err)
}

func (s *Single) OnRequest(w http.ResponseWriter, r *http.Request) {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnRequest(w, r)
}

func (s *Single) OnAlarm() {
	s.m.Lock()
	defer s.m.Unlock()
	s.s.OnAlarm()
}
