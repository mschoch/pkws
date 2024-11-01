package pkws

import (
	"fmt"
	"log/slog"
	"sync"
)

// Service keeps track of the known rooms and allows
// for registration of specific Server instances
// to handle specific parties.  Rooms without
// a server registered will use the default server
// which simply logs requests.
type Service struct {
	m     sync.RWMutex
	rooms map[string]*Room
	log   *slog.Logger

	defaultServer ServerConstructor
	servers       map[string]ServerConstructor
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		rooms:         make(map[string]*Room),
		defaultServer: NewLogger,
		servers:       make(map[string]ServerConstructor),
		log:           log,
	}
}

func (s *Service) Room(party string, room string) *Room {
	roomKey := s.roomKey(party, room)

	s.m.RLock()
	if room, exists := s.rooms[roomKey]; exists {
		s.m.RUnlock()
		return room
	}
	s.m.RUnlock()

	s.m.Lock()
	defer s.m.Unlock()

	serverCons := s.defaultServer
	if sc, ok := s.servers[party]; ok {
		serverCons = sc
	}
	r := NewRoom(party, room, serverCons, s.log)
	s.rooms[roomKey] = r

	r.s.OnStart()
	return r
}

func (s *Service) roomKey(party string, room string) string {
	return fmt.Sprintf("%s/%s", party, room)
}

func (s *Service) RegisterServer(party string, cons ServerConstructor) {
	s.servers[party] = cons
}
