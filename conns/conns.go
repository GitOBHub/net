package conns

import (
	"bufio"
	"io"
	//	"log"
	"net"
	"strconv"
	"strings"
)

type ConnInterface interface {
	net.Conn
	New(net.Conn) ConnInterface
	Recv() ([]byte, error)
	Send([]byte) (int, error)
	IsConnected() bool
	SetConnected(bool)
}

type Conn struct {
	net.Conn
	Connected bool
	reader    *bufio.Reader
}

func NewConn(c net.Conn) *Conn {
	conn := new(Conn)
	conn.Conn = c
	conn.Connected = true
	conn.reader = bufio.NewReader(c)
	return conn
}

func (conn *Conn) New(c net.Conn) ConnInterface {
	return NewConn(c)
}

func (conn *Conn) Recv() ([]byte, error) {
	dataLenStr, err := conn.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	dataLenStr = strings.TrimSuffix(dataLenStr, "\r\n")
	dataLen, _ := strconv.Atoi(dataLenStr)
	data := make([]byte, dataLen)
	conn.reader.Read(data)
	//	log.Printf("read: %q", string(data))
	return data, nil
}

func (conn *Conn) Send(data []byte) (int, error) {
	dataLen := strconv.Itoa(len(data))
	toSend := dataLen + "\r\n" + string(data)
	//log.Printf("send: %q", toSend)
	return io.WriteString(conn, toSend)
}

func (conn *Conn) IsConnected() bool {
	return conn.Connected
}

func (conn *Conn) SetConnected(status bool) {
	conn.Connected = status
}
