# pkws

a [PartyKit](https://www.partykit.io/)-like web serving library for Go

## Status

this project, like this README, is a work in progress

## Concepts

Small group communication, coordinated at a URL, we call a Room.

Each room can be programmed with custom handling of websocket messages and HTTP requests to child URLs, we call this a Server.  The programming model is simple because a single instance of the Server processes messages and requests one at a time.

Each room has a private key/value store.  Currently, we only have an ephemeral implementation of this storage layer.


## Usage

Build simple servers with the following interface:

```go
type Server interface {
	OnStart()
	OnConnect(conn *Connection)
	OnMessage(msg []byte, sender *Connection)
	OnClose(conn *Connection)
	OnError(conn *Connection, err error)
	OnRequest(w http.ResponseWriter, r *http.Request)
	OnAlarm()
}
```

If you want to manage a bunch of rooms for several different types of server behavior:

```go
// a custom environment struct your server will be passed
type Env struct{}
env := &Env{}

ps := pkws.NewPartyServer(env, "/parties", log)

// build a custom server, and define a constructor that takes our custom environment
func NewServer(env *Env, room *pkws.Room, log *slog.Logger) pkws.Server

// register the custom server for the fireproof party
ps.RegisterServer("fireproof", NewServer)
```