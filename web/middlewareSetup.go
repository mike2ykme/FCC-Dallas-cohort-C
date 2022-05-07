package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/models"
)

func MiddlewareSetup(cfg *models.Configuration) {

	cfg.WebApp.Use(logger.New())
	cfg.WebApp.Use(cors.New())

	if cfg.Production {
		// rate limiting
		cfg.WebApp.Use(limiter.New(limiter.Config{
			Max:        cfg.LimiterConfig.Max,
			Expiration: cfg.LimiterConfig.ExpirationSeconds,
		}))
	} else {
		// dev chat page
		cfg.WebApp.Static("/", "./static/home.html")
		// performance monitoring w/ page
		cfg.WebApp.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config
	}

}
