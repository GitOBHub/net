package main

import (
	"github.com/gitobhub/net/server"
	"log"
)

type EchoHandler struct{}

func main() {
	handler := EchoHandler{}
	srv := server.NewServer(":5000", handler)
	log.Fatal(srv.ListenAndServe())
}

func (h EchoHandler) HandleMessage(c server.ConnInterface, b []byte) {
	c.Write(b)
}

func (h EchoHandler) HandleConn(c server.ConnInterface) {
}
