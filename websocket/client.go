package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const ws_addres = "ws://localhost:2020"

type whatsappClient struct {
	whatsappWeb *websocket.Conn
}

func (wpclient *whatsappClient) connectWebsocket() (c *websocket.Conn) {
	c, _, err := websocket.DefaultDialer.Dial(ws_addres, nil)
	wpclient.whatsappWeb = c
	if err != nil {
		log.Fatal("dial error:", err)
		c.Close()
	}
	return
}

type client struct {
	id     int
	socket *websocket.Conn
	wpclientConn whatsappClient
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, byte_msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		s := string(byte_msg)
		contain := strings.Contains(s, "backend-connectWhatsApp")
		if contain {
			c.wpclientConn.connectWebsocket()
		}

		log.Printf("client: %d; message: %s", c.id, s)
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for {
		err := c.socket.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Test channel %d", c.id)))
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
	wpclient := whatsappClient{}
	client_count++
	client := &client{
		id:     client_count,
		socket: conn,
		wpclientConn: wpclient,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.read()
	go client.write()
}
