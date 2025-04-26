package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebSocketHandler(c *fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		AddClient(conn)
		defer RemoveClient(conn)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})(c)
}
