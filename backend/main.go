// backend/main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (dev only)
	},
}

var doc = &CRDT{}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// Add client
	clients[ws] = true

	// Send current document
	ws.WriteMessage(websocket.TextMessage, []byte(doc.GetText()))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			delete(clients, ws)
			break
		}

		// Apply operation
		var op Operation
		err = json.Unmarshal(msg, &op)
		if err != nil {
			fmt.Println("Invalid operation:", err)
			continue
		}
		doc.ApplyOp(op)

		// Send updated document to all clients
		broadcast <- doc.GetText()
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleBroadcast()
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
