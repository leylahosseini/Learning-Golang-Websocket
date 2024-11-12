package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	Room string `json:"room,omitempty"`
}

type Client struct {
	Name string
	Conn *websocket.Conn
	Room string
}

var clients = make(map[*Client]bool)                         // لیست کلاینت‌ها
var broadcast = make(chan Message)                           // کانال برای پیام‌ها
var rooms = []string{"room1", "General", "DevOps", "Golang"} // روم‌های از پیش تعریف‌شده
var mu sync.Mutex

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("server is running :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("server is error ....", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // مجاز کردن هر-Origin
	},
}

// این تابع به کلاینت‌ها امکان می‌دهد به سرور متصل شوند
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error in upgrade", err)
		return
	}
	defer conn.Close()

	var clientName string
	var roomName string

	// ارسال لیست روم‌ها به کلاینت
	if err := conn.WriteJSON(rooms); err != nil {
		log.Println("error is show rooms", err)
		return
	}

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("error in reading message  ", err)
			break
		}

		if msg.Type == "join" {
			clientName = msg.Name
			roomName = msg.Room

			mu.Lock()
			clients[&Client{Name: clientName, Conn: conn, Room: roomName}] = true
			mu.Unlock()
			log.Printf("%s room %s join\n", clientName, roomName)

			continue
		}

		broadcast <- msg
	}
}

// این تابع پیام‌های دریافتی را به همه کلاینت‌ها ارسال می‌کند
func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			if client.Room == msg.Room { // ارسال فقط به روم مربوطه
				if err := client.Conn.WriteJSON(msg); err != nil {
					log.Println("Error in start server", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}
