package main

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

func main() {
	fmt.Println("Server starting...")
	s := WSServer{conns: make(map[*websocket.Conn]WSUser)}
	http.HandleFunc("/ws", s.ServeWS)
	http.HandleFunc("/users", s.ServeUsers)
	http.Handle("/", http.FileServer(http.Dir("dist/")))
	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}

