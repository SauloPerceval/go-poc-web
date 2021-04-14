package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHookMessage struct {
	Text string `json:"text"`
}

func receiveMessage(c *gin.Context) {
	var message WebHookMessage

	if err := c.BindJSON(&message); err == nil {
		log.Printf("message received: %s", message.Text)
		c.String(http.StatusOK, "Ok")
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}

}

func main() {
	var addr string
	flag.StringVar(&addr, "p", "5000", "http service port")
	flag.Parse()
	r := gin.Default()

	r.POST("/receive", receiveMessage)

	r.Run(":" + addr)
}
