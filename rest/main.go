package main

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WebHook struct {
	Url string `json:"url"`
}

type Message struct {
	Text string `json:"text"`
}

type Hub struct {
	wh_cliets_map map[int]*WebHookClient
}

var wh_count int

func (h *Hub) RegisterWebhook(c *gin.Context) {
	var webhook WebHook

	if err := c.BindJSON(&webhook); err == nil {
		wh_count++

		new_client := WebHookClient{
			url: webhook.Url,
		}
		go new_client.StartReadingSocket()

		h.wh_cliets_map[wh_count] = &new_client

		c.JSON(http.StatusOK, gin.H{
			"id": strconv.Itoa(wh_count),
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
	}
	r := gin.Default()

	r.POST("/webhook", h.RegisterWebhook)
	r.POST("/send/:webhook_id", h.SendMessage)
	r.Use((cors.Default()))
	r.Run(":8000")
}
