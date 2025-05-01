package main

import (
	"log"
	router "server/Router"
	"server/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	config.ConnectDB()

	app := fiber.New()

	router.SetupRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatal("Failed to connect to server")
	}

}
