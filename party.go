package pkws

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// PartyServer is an HTTP Handler to implement partykit-like parties and rooms
// with the Server logic written in Go
type PartyServer struct {
	svc    *Service
	log    *slog.Logger
	router *mux.Router
}

func NewPartyServer(svc *Service, log *slog.Logger) *PartyServer {
	ps := &PartyServer{
		svc: svc,
		log: log,
	}

	router := mux.NewRouter()
	router.Handle("/parties/{party}/{room}", http.HandlerFunc(ps.handleHTTP)).Methods("GET", "PUT", "POST", "DELETE", "OPTIONS")

	ps.router = router

	return ps
}

func (p *PartyServer) getRoom(r *http.Request) (*Room, error) {
	vars := mux.Vars(r)
	partyName, ok := vars["party"]
	if !ok {
		p.log.Error("missing party")
		return nil, fmt.Errorf("missing party")
	}
	roomName, ok := vars["room"]
	if !ok {
		p.log.Error("missing room")
		return nil, fmt.Errorf("missing room")
	}

	return p.svc.Room(partyName, roomName), nil
}

func (p *PartyServer) handleWS(w http.ResponseWriter, r *http.Request) {

	room, err := p.getRoom(r)
	if err != nil {
		p.log.Error(err.Error())
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		p.log.Info("error upgrading websocket", "err", err)
		return
	}

	connectionId := r.FormValue("_pk")
	if len(connectionId) == 0 {
		p.log.Error("missing _pk")
		return
	}

	wsConn := &Connection{
		id:     connectionId,
		send:   make(chan []byte, 256),
		ws:     conn,
		parent: room,
	}

	room.Add(wsConn)
	defer room.Remove(wsConn)

	go wsConn.Writer()
	wsConn.Reader()
}

func (p *PartyServer) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		p.handleWS(w, r)
		return
	}
	p.log.Info("in handleHTTP")
	room, err := p.getRoom(r)
	if err != nil {
		p.log.Error(err.Error())
		return
	}
	room.s.OnRequest(w, r)
}

func (p *PartyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.log.Info("in serve http", "url", r.URL.String())
	p.router.ServeHTTP(w, r)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
