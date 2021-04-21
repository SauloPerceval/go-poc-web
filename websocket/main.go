package main

import (
	"log"
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("Webscoket Service Running on 8080... ")
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
