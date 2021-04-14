package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

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

		requestBody, err := json.Marshal(map[string]string{
			"text": string(message),
		})

		if err != nil {
			log.Println("read error:", err)
			return
		}

		resp, err := http.Post(whc.url, "application/json", bytes.NewBuffer(requestBody))

		if err != nil {
			log.Printf("Request to %s failed: %s \n", whc.url, err.Error())
			return
		}

		if resp.StatusCode != 200 {
			log.Printf("Request to %s return status %d\n", whc.url, resp.StatusCode)
			return
		}
	}
}

func (whc *WebHookClient) SendMessageOnSocket(msg string) {
	err := whc.c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write error:", err)
		return
	}
}
