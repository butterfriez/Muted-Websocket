// https://stackoverflow.com/questions/57783858/sending-a-websocket-message-to-a-specific-channel-in-go-using-gorilla
package main

import (
	"flag"
	"log"
	"muted/util"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	hub := util.NewHub()
	flag.Parse()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		util.ServeWS(hub, w, r)
	})
	log.Println("Serving on port :8080")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
