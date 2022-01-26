package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/uwaifo/shortmanredis/api/routes"

	"github.com/joho/godotenv"
)

// setupRoutes holds a listing of all the routes that are available to the API.
func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)

}
func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)
	//log.Fatal(app.Listen(os.Getenv("APP_PORT")))

	app.Listen(":3000")

}
