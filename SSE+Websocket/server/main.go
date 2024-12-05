package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go broadcastServerInfo()

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()


	clients[conn] = true
	fmt.Println("New connection added!")


	defer func() {
		delete(clients, conn)
		fmt.Println("Connection removed!")
	}()


	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
	}
}

func broadcastServerInfo() {
	for {
		serverTime := time.Now().Format("15:04:05")
		connectionCount := len(clients)
		message := fmt.Sprintf(`{"server_time":"%s", "connections":%d}`, serverTime, connectionCount)

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
