package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/joho/godotenv"
)

// setupRoutes holds a listing of all the routes that are available to the API.
func setupRoutes(app *fiber.App) {
	//app.Get("/:url", routes.ResolveURL)

}
func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New()

	app.Use(logger.New)

	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
