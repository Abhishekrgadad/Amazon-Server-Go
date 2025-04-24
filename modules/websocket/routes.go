package websocket

import (
	"github.com/gofiber/fiber/v2"
)

func WebSocketRoutes(router fiber.Router) {
	router.Get("/ws", WebSocketHandler)
}
