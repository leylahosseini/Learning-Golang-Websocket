package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	// آدرس سرور وب‌سوکت
	url := "ws://localhost:8080/ws"

	// اتصال به سرور وب‌سوکت
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// ارسال پیام به سرور
	go func() {
		for {
			var msg string
			fmt.Print("Enter message: ")
			_, err := fmt.Scanln(&msg)
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}()

	// دریافت پیام از سرور
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Printf("Received: %s\n", msg)
	}
}
