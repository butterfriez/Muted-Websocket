// https://stackoverflow.com/questions/57783858/sending-a-websocket-message-to-a-specific-channel-in-go-using-gorilla
package main

import (
	"log"
	"muted/util"
	"net/http"
)

func main() {
	hub := util.NewHub()

	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		util.ServeWS(hub, w, r)
	})

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
