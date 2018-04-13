package server

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	Address  string
	Handler  MessageCloseHandler
	ConnType ConnInterface
	mu       sync.Mutex
	numConn  int
}

type MessageHandlerFunc func(ConnInterface, []byte)
type ConnectionHandlerFunc func(ConnInterface)

type MessageHandler interface {
	HandleMessage(ConnInterface, []byte)
}

type CloseHandler interface {
	HandleClose(ConnInterface)
}

type MessageCloseHandler interface {
	MessageHandler
	CloseHandler
}

//NopCloser
type NopCloser struct {
	MessageHandler
}

func (h NopCloser) HandleClose(ConnInterface) {
	//null
}

func NopCloseHandler(h MessageHandler) MessageCloseHandler {
	return NopCloser{h}
}

//TODO:
var defaultHandler MessageCloseHandler

func NewServer(addr string, handler MessageHandler) *Server {
	mc, ok := handler.(MessageCloseHandler)
	if !ok && handler != nil {
		mc = NopCloseHandler(handler)
	}
	srv := &Server{Address: addr, Handler: mc}
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
		s.Handler.HandleClose(conn)
	}

}
