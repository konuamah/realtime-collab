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
		return true // Dev only
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

	clients[ws] = true

	// Send current document as JSON
	initialData, _ := json.Marshal(map[string]string{
		"type": "update",
		"text": doc.GetText(),
	})
	ws.WriteMessage(websocket.TextMessage, initialData)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			delete(clients, ws)
			break
		}

		var op Operation
		if err := json.Unmarshal(msg, &op); err != nil {
			fmt.Println("Invalid operation:", err)
			continue
		}

		doc.ApplyOp(op)
		broadcast <- doc.GetText() // trigger broadcast
	}
}

func handleBroadcast() {
	for {
		text := <-broadcast
		data, _ := json.Marshal(map[string]string{
			"type": "update",
			"text": text,
		})
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, data)
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
