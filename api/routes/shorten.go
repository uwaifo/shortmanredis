package routes

import "time"

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
