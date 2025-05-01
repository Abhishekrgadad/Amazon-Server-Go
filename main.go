package amz

import (
	"errors"
	"fmt"
	"log"
	"os"
	"server/config"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil{
		log.Fatal("Failed to load .env file")
	}

	config.ConnectDB()

	app := fiber.New()

	// routes.SetupRoutes(app)
	err = app.Listen(":3000")
	if err != nil {
		log.Fatal("Failed to connect to server")
	}

}