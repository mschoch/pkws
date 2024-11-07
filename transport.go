package pkws

import "net/http"

type RoomHandler struct {
	room *Room
}

func NewRoomHandler(room *Room) *RoomHandler {
	return &RoomHandler{
		room: room,
	}
}

func (r *RoomHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		r.handleWS(w, req)
		return
	}
	r.room.log.Info("in handleHTTP")
	r.room.s.OnRequest(w, req)
}

func (r *RoomHandler) handleWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.room.log.Info("error upgrading websocket", "err", err)
		return
	}

	connectionId := req.FormValue("_pk")
	if len(connectionId) == 0 {
		r.room.log.Error("missing _pk")
		return
	}

	wsConn := &Connection{
		id:     connectionId,
		send:   make(chan []byte, 256),
		ws:     conn,
		parent: r.room,
	}

	r.room.Add(wsConn)
	defer r.room.Remove(wsConn)

	go wsConn.Writer()
	wsConn.Reader()
}
