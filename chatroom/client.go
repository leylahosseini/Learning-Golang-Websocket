package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	Room string `json:"room,omitempty"`
}

func main() {
	// Connecting to WebSocket
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// Receiving the list of rooms from the server
	var rooms []string
	if err := conn.ReadJSON(&rooms); err != nil {
		log.Fatal("Error reading room list:", err)
	}

	// Displaying rooms to the user
	fmt.Println("Available rooms:")
	for i, room := range rooms {
		fmt.Printf("%d: %s\n", i+1, room)
	}

	var roomChoice int
	fmt.Print("Please enter the room number you want to join: ")
	fmt.Scanln(&roomChoice)

	if roomChoice < 1 || roomChoice > len(rooms) {
		fmt.Println("Invalid choice!")
		return
	}

	roomName := rooms[roomChoice-1]
	var name string
	fmt.Print("Please enter your name: ")
	fmt.Scanln(&name)

	// Sending user name and room name to the server
	initialMessage := Message{Type: "join", Name: name, Room: roomName}
	if err := conn.WriteJSON(initialMessage); err != nil {
		log.Fatal("Error sending join message:", err)
	}

	// Processing messages
	go func() {
		for {
			var msg Message
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("%s: %s\n", msg.Name, msg.Text)
		}
	}()

	// Receiving user input and sending messages
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your message: ")
		scanner.Scan()
		text := scanner.Text()

		if text == "exit" {
			log.Println("Exiting chat...")
			break
		}

		message := Message{Text: text, Room: roomName} // sending message along with room name
		fmt.Println(message)
		if err := conn.WriteJSON(message); err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
