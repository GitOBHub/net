package server

import (
	"log"
	"net"
	"sync"

	"github.com/gitobhub/net/conns"
)

type Server struct {
	Address  string
	Handler  Handler
	ConnType conns.ConnInterface
	mu       sync.Mutex
	numConn  int
}

type MessageHandlerFunc func(conns.ConnInterface, []byte)
type ConnectionHandlerFunc func(conns.ConnInterface)

type Handler interface {
	HandleMessage(conns.ConnInterface, []byte)
	HandleConn(conns.ConnInterface)
}

//TODO:
var defaultHandler Handler

func NewServer(addr string, handler Handler) *Server {
	srv := &Server{Address: addr, Handler: handler}
	return srv
}

func (s *Server) SetConnType(c conns.ConnInterface) {
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

		var conn conns.ConnInterface
		if s.ConnType == nil {
			conn = conns.NewConn(c)
		} else {
			conn = s.ConnType.New(c)
		}
		log.Printf("connection#%d %s -> %s is up", s.numConn, conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn conns.ConnInterface) {
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
