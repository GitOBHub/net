package server

import (
	"log"
	"net"
	"sync"

	"net/conns"
)

type Server struct {
	Address        string
	Listener       net.Listener
	NumConn        int
	mu             sync.Mutex
	messageHandler HandlerFunc
}

type HandlerFunc func(*conns.Connection, []byte)

var defaultMessageHandler HandlerFunc

func (f HandlerFunc) handle(conn *conns.Connection, data []byte) {
	f(conn, data)
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv := &Server{Address: addr, Listener: ln}
	return srv, nil
}

func (s *Server) HandleFunc(handler func(*conns.Connection, []byte)) {
	s.messageHandler = HandlerFunc(handler)
}

func (s *Server) Serve() {
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		s.mu.Lock()
		s.NumConn++
		s.mu.Unlock()
		conn := &conns.Connection{Conn: c, Number: s.NumConn}
		log.Printf("connection#%d is up", conn.Number)
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *conns.Connection) {
	defer func() {
		conn.Close()
		log.Printf("connection#%d is down", conn.Number)
	}()
	for {
		data, _ := conn.Read()
		if data == nil {
			break
		}
		s.messageHandler.handle(conn, data)
	}
}
