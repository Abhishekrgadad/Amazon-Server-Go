package main

import (
	"log"
	"os"
	"server/config"
	"server/modules/websocket"
	"server/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDatabase()
	config.ConnectRedis()
	app := fiber.New()

	go websocket.StartBroadcast()

	router.SetupRoutes(app)
	websocket.WebSocketRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
