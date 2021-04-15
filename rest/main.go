package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebHook struct {
	Url string `json:"url"`
}

type Message struct {
	Text string `json:"text"`
}

type Hub struct {
	wh_clients_count int
	wh_cliets_map    map[int]*WebHookClient
	add_client       chan *WebHookClient
	remove_client    chan *WebHookClient
}

func (h *Hub) MonitorClients() {
	for {
		select {
		case new_client := <-h.add_client:
			h.wh_clients_count++
			new_client.id = h.wh_clients_count
			h.wh_cliets_map[h.wh_clients_count] = new_client
			go new_client.StartReadingSocket(h.remove_client)
			log.Println("Starting redirecting messages to client ", new_client.id)
		case exiting_client := <-h.remove_client:
			delete(h.wh_cliets_map, exiting_client.id)
			log.Println("Stop redirecting messages to client ", exiting_client.id)
		}
	}
}

func (h *Hub) RegisterWebhook(c *gin.Context) {
	var webhook WebHook

	if err := c.BindJSON(&webhook); err == nil {
		new_client := WebHookClient{
			url: webhook.Url,
		}
		h.add_client <- &new_client

		c.JSON(http.StatusOK, gin.H{
			"id": strconv.Itoa(new_client.id),
		})
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func (h *Hub) SendMessage(c *gin.Context) {
	var message Message

	id, err := strconv.Atoi(c.Param("webhook_id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Must pass a integer number as id")
	}

	if err := c.BindJSON(&message); err == nil {
		h.wh_cliets_map[id].SendMessageOnSocket(message.Text)
		c.String(http.StatusOK, "Message sent")
	}
}

func main() {
	h := &Hub{
		wh_cliets_map: make(map[int]*WebHookClient),
		add_client:    make(chan *WebHookClient),
		remove_client: make(chan *WebHookClient),
	}
	go h.MonitorClients()

	r := gin.Default()

	r.POST("/webhook", h.RegisterWebhook)
	r.POST("/send/:webhook_id", h.SendMessage)

	r.Run(":8000")
}
