package server

import (
	"log"
	"net"
	"sync"

	"github.com/GitOBHub/net/conns"
)

type Server struct {
	Address           string
	NumConn           int
	mu                sync.Mutex
	messageHandler    MessageHandlerFunc
	connectionHandler ConnectionHandlerFunc
}

type MessageHandlerFunc func(*conns.Connection, []byte)
type ConnectionHandlerFunc func(*conns.Connection)

var defaultMessageHandler MessageHandlerFunc

func (f MessageHandlerFunc) handle(conn *conns.Connection, data []byte) {
	f(conn, data)
}

func (f ConnectionHandlerFunc) handle(conn *conns.Connection) {
	f(conn)
}

func NewServer(addr string) *Server {
	srv := &Server{Address: addr}
	return srv
}

func (s *Server) MessageHandleFunc(handler func(*conns.Connection, []byte)) {
	s.messageHandler = MessageHandlerFunc(handler)
}

func (s *Server) ConnectionHandleFunc(handler func(*conns.Connection)) {
	s.connectionHandler = ConnectionHandlerFunc(handler)
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		s.mu.Lock()
		s.NumConn++
		s.mu.Unlock()
		conn := &conns.Connection{Conn: c, Number: s.NumConn, Connected: true}
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
		data, _ := conn.Recv()
		if data == nil {
			break
		}
		if s.messageHandler != nil {
			s.messageHandler.handle(conn, data)
		}
	}
	conn.Connected = false
	if s.connectionHandler != nil {
		s.connectionHandler.handle(conn)
	}
}
