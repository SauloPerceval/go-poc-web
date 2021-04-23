package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

const ws_addres = "ws://localhost:8080/ws"

type WebHookClient struct {
	url  string
	c    *websocket.Conn
	open bool
}

func create_qrcode_file(file string){
	fmt.Println(file)
	message1 := strings.Split(file, `"image":`)
	message2 := strings.Split(message1[1], ":")
	message3 := strings.Split(message2[1], "base64,")
	qrContent := strings.Split(message3[1], `",`)
	dec, err := base64.StdEncoding.DecodeString(qrContent[0])
	f, err := os.Create("qrcode.svg")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if _, err := f.Write([]byte(dec)); err != nil {
        panic(err)
    }
    if err := f.Sync(); err != nil {
        panic(err)
    }
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
		if strings.Contains(string(message), "generated_qr_code"){
			qrcode_content := string(message)
			create_qrcode_file(qrcode_content)
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
