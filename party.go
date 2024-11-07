package pkws

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// PartyServer is an HTTP Handler to implement partykit-like parties and rooms
// with the Server logic written in Go
type PartyServer[Env any] struct {
	m     sync.RWMutex
	rooms map[string]*Room

	servers map[string]ServerConstructor[Env]

	log    *slog.Logger
	env    Env
	router *mux.Router
}

func NewPartyServer[Env any](env Env, prefix string, log *slog.Logger) *PartyServer[Env] {
	ps := &PartyServer[Env]{
		log:     log,
		env:     env,
		rooms:   make(map[string]*Room),
		servers: make(map[string]ServerConstructor[Env]),
	}

	path := "/{party}/{room}"
	if len(prefix) > 0 {
		path = prefix + path
	}

	router := mux.NewRouter()
	router.Handle(path, http.HandlerFunc(ps.handleHTTP)).Methods("GET", "PUT", "POST", "DELETE", "OPTIONS")

	ps.router = router

	return ps
}

func (p *PartyServer[Env]) RegisterServer(namespace string, cons ServerConstructor[Env]) {
	p.servers[namespace] = cons
}

func (p *PartyServer[Env]) RoomFromRequest(r *http.Request) (*Room, error) {
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

	return p.RoomFromName(partyName, roomName)
}

func (p *PartyServer[Env]) RoomFromName(namespace string, name string) (*Room, error) {
	roomKey := p.roomKey(namespace, name)

	p.m.RLock()
	if room, exists := p.rooms[roomKey]; exists {
		p.m.RUnlock()
		return room, nil
	}
	p.m.RUnlock()

	p.m.Lock()
	defer p.m.Unlock()

	serverCons, ok := p.servers[namespace]
	if !ok {
		return nil, fmt.Errorf("no server registered for namespace %q", namespace)
	}

	r := NewRoom(namespace, name, p.log)
	s := serverCons(p.env, r, p.log)
	r.AttachServer(s)

	p.rooms[roomKey] = r

	r.s.OnStart()
	return r, nil
}

func (p *PartyServer[Env]) roomKey(party string, room string) string {
	return fmt.Sprintf("%s/%s", party, room)
}

func (p *PartyServer[Env]) handleHTTP(w http.ResponseWriter, r *http.Request) {
	room, err := p.RoomFromRequest(r)
	if err != nil {
		p.log.Error(err.Error())
		return
	}

	rh := NewRoomHandler(room)
	rh.ServeHTTP(w, r)
}

func (p *PartyServer[Env]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.log.Info("in serve http", "url", r.URL.String())
	p.router.ServeHTTP(w, r)
}
