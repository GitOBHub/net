package conns

import (
	"bufio"
	"io"
	//	"log"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	net.Conn
	Number    int
	Connected bool
}

func (conn *Connection) Recv() ([]byte, error) {
	rd := bufio.NewReader(conn.Conn)
	dataLenStr, err := rd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	dataLenStr = strings.TrimSuffix(dataLenStr, "\r\n")
	dataLen, _ := strconv.Atoi(dataLenStr)
	data := make([]byte, dataLen)
	rd.Read(data)
	//	log.Printf("read: %q", string(data))
	return data, nil
}

func (conn *Connection) Send(data []byte) (int, error) {
	dataLen := strconv.Itoa(len(data))
	toSend := dataLen + "\r\n" + string(data)
	//	log.Printf("send: %q", toSend)
	return io.WriteString(conn, toSend)
}
