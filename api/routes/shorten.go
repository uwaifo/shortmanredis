package routes

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/uwaifo/shortmanredis/api/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiration  time.Duration `json:"expiration"`
}

type resposne struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiration      time.Duration `json:"expiration"`
	XRateRemaining  int           `json:"x_rate_remaining"`
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

// ShortenURL is a function that takes a request and returns a response.
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO Impliment the rate limiting here

	// First chech if the input is an actual URL
	if !valid.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// Allow the user to call the api 10 times in 30 minutes only
	// Check for domain error (eg in localhost to avoid getting into an infinit loop)

	if !helpers.RemoteDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// Enforce https , SSL

	body.URL = helpers.EnforceHTTP(body.URL)
	return nil

}
