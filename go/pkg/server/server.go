package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	sockets   map[uint]*Socket
	currentID uint
	lock      sync.Mutex

	In          chan *Message
	Out        chan *Message
	fromSocket chan *Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow any origin
} // use default options

func NewServer() (*Server, error) {
	out := make(chan *Message, 10000)
	fromSocket := make(chan *Message, 10000)
	server := Server{
		currentID:  0,
		sockets:    make(map[uint]*Socket),
		In:         make(chan *Message),
		fromSocket: fromSocket,
		Out:        out,
		lock:       sync.Mutex{},
	}

	go func() {
		for {
			select {
			case msg := <-out:
				server.lock.Lock()
				server.sockets[msg.ID].Out <- msg
				server.lock.Unlock()

			case msg := <-fromSocket:
				if msg.Type == websocket.CloseMessage {
					server.lock.Lock()
					delete(server.sockets, msg.ID)
					server.lock.Unlock()
				}
				server.In <- msg
			}
		}
	}()

	return &server, nil
}

func (s *Server) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	s.lock.Lock()
	id := s.currentID
	s.currentID += 1

	socket, err := NewSocket(id, s.fromSocket, w, r)

	s.sockets[id] = socket
	s.lock.Unlock()

	if err != nil {
		log.Print("couldn't upgrade socket.", err)
		return
	}
}
