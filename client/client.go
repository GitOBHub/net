package client

import (
	"net"

	"github.com/gitobhub/net/conns"
)

type Client struct {
	RemoteAddr string
	Conn       conns.Connection
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client := &Client{RemoteAddr: addr}
	client.Conn = conns.Connection{Conn: conn}
	return client, nil
}

func (c *Client) Read() ([]byte, error) {
	return c.Conn.Read()
}

func (c *Client) Send(data []byte) {
	c.Conn.Send(data)
}
