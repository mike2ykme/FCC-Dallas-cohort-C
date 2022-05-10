package models

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type Configuration struct {
	Port          string
	Production    bool
	WebApp        *fiber.App
	LimiterConfig LimiterConfig
	JwtSecret     string
}

type LimiterConfig struct {
	Max               int
	ExpirationSeconds time.Duration
}
