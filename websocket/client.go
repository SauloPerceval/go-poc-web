package main

import (
	//"fmt"
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
	open bool
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

func (c *client) receive() {
	defer func (){
		c.socket.Close()
		c.wpclientConn.whatsappWeb.Close()
	}()
	for {
		_, byte_msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		s := string(byte_msg)
		log.Printf("client: %d; message received: %s", c.id, s)
		contain := strings.Contains(s, "backend-connectWhatsApp")
		if contain {
			c.wpclientConn.connectWebsocket()
			c.wpclientConn.open = true
		}

		err = c.wpclientConn.whatsappWeb.WriteMessage(
			websocket.TextMessage, []byte(strings.Replace(string(byte_msg), "´", `"`, -1)))
		if err != nil {
			return
		}

		log.Printf("client: %d; Message sent: %s", c.id, string(byte_msg))
	}
}

func (c *client) send() {
	defer func (){
		c.socket.Close()
		c.wpclientConn.whatsappWeb.Close()
	}()
	for {
		if c.wpclientConn.open {
			_, recv_msg, err := c.wpclientConn.whatsappWeb.ReadMessage()
			if err != nil {
				return
			}
			err = c.socket.WriteMessage(websocket.TextMessage, recv_msg)
			if err != nil {
				return
			}
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
	wpClient := whatsappClient{}
	client_count++
	client := &client{
		id:     client_count,
		socket: conn,
		wpclientConn: wpClient,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.receive()
	go client.send()
}
