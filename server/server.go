package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // مجاز کردن اتصالات از هر منبع
	},
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()
	mu.Lock()
	clients[conn] = true
	mu.Unlock()
	fmt.Println("New client connected")

	for {
		//var msg string
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected:", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}

		// ارسال پیام به همه کلاینت‌ها
		mu.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error sending message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error while upgrading connection:", err)
			return
		}

		handleConnection(conn)
	})

	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
