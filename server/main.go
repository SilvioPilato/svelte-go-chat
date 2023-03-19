package main

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

func main() {
	fmt.Println("Server starting...")
	s := WSServer{conns: make(map[*websocket.Conn]WSUser), logged: make(map[string]bool)}
	http.HandleFunc("/ws", s.ServeWS)
	http.Handle("/", http.FileServer(http.Dir("dist/")))
	fmt.Println("Server started")
	http.ListenAndServe(":8080", nil)
}