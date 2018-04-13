package server

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	Address  string
	Handler  Handler
	ConnType ConnInterface
	mu       sync.Mutex
	numConn  int
}

type MessageHandlerFunc func(ConnInterface, []byte)
type ConnectionHandlerFunc func(ConnInterface)

type Handler interface {
	HandleMessage(ConnInterface, []byte)
	HandleConn(ConnInterface)
}

//TODO:
var defaultHandler Handler

func NewServer(addr string, handler Handler) *Server {
	srv := &Server{Address: addr, Handler: handler}
	return srv
}

func (s *Server) SetConnType(c ConnInterface) {
	s.ConnType = c
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
		s.numConn++

		var conn ConnInterface
		if s.ConnType == nil {
			conn = NewConn(c)
		} else {
			conn = s.ConnType.New(c)
		}
		log.Printf("connection#%d %s -> %s is up", s.numConn, conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn ConnInterface) {
	defer func() {
		conn.Close()
		log.Printf("connection %s -> %s is down", conn.RemoteAddr(), conn.LocalAddr())
	}()
	for {
		data, _ := conn.Recv()
		if data == nil {
			break
		}
		//TODO
		if s.Handler != nil {
			s.Handler.HandleMessage(conn, data)
		}
	}
	conn.SetConnected(false)
	if s.Handler != nil {
		s.Handler.HandleConn(conn)
	}

}
