package websocket

import (
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

// Add client
func AddClient(conn *websocket.Conn) {
	clients[conn] = true
}

// Remove client
func RemoveClient(conn *websocket.Conn) {
	delete(clients, conn)
}

// Trigger notification
func SendOrderNotification(message string) {
	broadcast <- message
}

func CancelOrderNotification(message string) {
	broadcast <- message
}

func ReturnOrderNotification(message string) {
	broadcast <- message
}
// Start broadcaster
func StartBroadcast() {
	for {
		msg := <-broadcast
		for conn := range clients {
			conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}
