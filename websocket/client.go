package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	id     int
	socket *websocket.Conn
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, byte_msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		s := string(byte_msg)
		log.Printf("client: %d; message: %s", c.id, s)
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for {
		err := c.socket.WriteMessage(websocket.TextMessage, []byte("Teste"))
		if err != nil {
			return
		}
		<-time.After(2 * time.Second)
	}
}

var client_count int

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	client_count++
	client := &client{
		id:     client_count,
		socket: conn,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.read()
	go client.write()
}
