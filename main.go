package main

import (
	"log"
	"server/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()
	
	app := fiber.New()

	app.Get("/",func (c *fiber.Ctx) error {
		return c.SendString("Amazon clone api")
	})

	app.Listen(":3000")
}