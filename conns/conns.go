package conns

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	net.Conn
	Number int
}

func (conn *Connection) Read() ([]byte, error) {
	rd := bufio.NewReader(conn.Conn)
	dataLenStr, err := rd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	dataLenStr = strings.TrimSuffix(dataLenStr, "\r\n")
	dataLen, _ := strconv.Atoi(dataLenStr)
	data := make([]byte, dataLen)
	rd.Read(data)
	//	io.ReadFull(conn.Conn, data)
	return data, nil
}

func (conn *Connection) Send(data []byte) (int, error) {
	dataLen := strconv.Itoa(len(data))
	toSend := dataLen + "\r\n" + string(data)
	return io.WriteString(conn, toSend)
}
