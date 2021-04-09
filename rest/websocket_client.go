package main

import (
	"log"

	"github.com/gorilla/websocket"
)

const ws_addres = "ws://localhost:8080/ws"

type WebHookClient struct {
	url  string
	c    *websocket.Conn
	open bool
}

func connectWebsocket() (c *websocket.Conn) {
	c, _, err := websocket.DefaultDialer.Dial(ws_addres, nil)
	if err != nil {
		log.Fatal("dial error:", err)
		c.Close()
	}
	return
}

func (whc *WebHookClient) StartReadingSocket() {
	whc.c = connectWebsocket()
	whc.open = true
	defer func() {
		whc.c.Close()
		whc.open = false
	}()

	for {
		_, message, err := whc.c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		log.Printf("recv: %s; url: %s", message, whc.url)
	}
}

func (whc *WebHookClient) SendMessageOnSocket(msg string) {
	err := whc.c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write error:", err)
		return
	}
}
