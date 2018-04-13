package server

import (
	"bufio"
	//	"log"
	"net"
)

type ConnInterface interface {
	net.Conn
	New(net.Conn) ConnInterface
	IsConnected() bool
	SetConnected(bool)
	Recv() ([]byte, error)
}

type Conn struct {
	net.Conn
	Connected bool
	Reader    *bufio.Reader
}

func NewConn(c net.Conn) *Conn {
	conn := new(Conn)
	conn.Conn = c
	conn.Connected = true
	conn.Reader = bufio.NewReader(c)
	return conn
}

func (conn *Conn) New(c net.Conn) ConnInterface {
	return NewConn(c)
}

func (conn *Conn) IsConnected() bool {
	return conn.Connected
}

func (conn *Conn) SetConnected(status bool) {
	conn.Connected = status
}

func (conn *Conn) Recv() ([]byte, error) {
	p := make([]byte, 4096)
	n, err := conn.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return nil, err
}
