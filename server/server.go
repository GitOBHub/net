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

type MessageHandlerFunc func(*conns.Conn, []byte)
type ConnectionHandlerFunc func(*conns.Conn)

var defaultMessageHandler MessageHandlerFunc

func (f MessageHandlerFunc) handle(conn *conns.Conn, data []byte) {
	f(conn, data)
}

func (f ConnectionHandlerFunc) handle(conn *conns.Conn) {
	f(conn)
}

func NewServer(addr string) *Server {
	srv := &Server{Address: addr}
	return srv
}

func (s *Server) MessageHandleFunc(handler func(*conns.Conn, []byte)) {
	s.messageHandler = MessageHandlerFunc(handler)
}

func (s *Server) ConnectionHandleFunc(handler func(*conns.Conn)) {
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
		conn := conns.NewConn(c)
		log.Printf("connection#%d %s -> %s is up", s.NumConn, conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *conns.Conn) {
	defer func() {
		conn.Close()
		log.Printf("connection# %s -> %s is down", conn.RemoteAddr(), conn.LocalAddr())
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
