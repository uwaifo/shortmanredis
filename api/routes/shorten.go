package routes

import (
	"os"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/uwaifo/shortmanredis/api/database"
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

	// Impliment the rate limiting here

	redisClient := database.CreateClient(1)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = redisClient.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()

	} else {
		value, _ = redisClient.Get(database.Ctx, c.IP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := redisClient.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":           "Rate limit exceeded",
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}

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

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// NOTICE Start , why are we using another redis client (rc) here ?
	rc := database.CreateClient(0)
	defer rc.Close()

	value, _ = rc.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Custom short already exists"})
	}

	if body.Expiration == 0 {
		body.Expiration = time.Second * 24 * 3600
	}

	err = rc.Set(database.Ctx, id, body.URL, body.Expiration).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//NOTICE End

	redisClient.Decr(database.Ctx, c.IP())
	return nil

}
