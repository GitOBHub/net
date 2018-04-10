package server

import (
	"log"
	"net"
	"sync"

	"github.com/gitobhub/net/conns"
)

type Server struct {
	Address string
	NumConn int
	mu      sync.Mutex
	Handler Handler
}

type MessageHandlerFunc func(*conns.Conn, []byte)
type ConnectionHandlerFunc func(*conns.Conn)

type Handler interface {
	HandleMessage(*conns.Conn, []byte)
	HandleConn(*conns.Conn)
}

var defaultMessageHandler MessageHandlerFunc

func NewServer(addr string, handler Handler) *Server {
	srv := &Server{Address: addr, Handler: handler}
	return srv
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
		s.Handler.HandleMessage(conn, data)
	}
	conn.Connected = false
	s.Handler.HandleConn(conn)
}
