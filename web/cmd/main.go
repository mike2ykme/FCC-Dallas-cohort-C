package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"log"
	"teamC/models"
	"teamC/web"
	"time"
)

func main() {
	cfg := models.Configuration{}

	err := web.LoadConfiguration(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/login", web.LoginHandler(&cfg))

	app.Use(web.GetJwtMiddleware(&cfg))

	if cfg.Production {
		// rate limiting
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.LimiterConfig.Max,
			Expiration: cfg.LimiterConfig.ExpirationSeconds * time.Second,
		}))
	} else {
		// performance monitoring w/ page
		app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config
	}

	// Setup Routes
	if !cfg.Production {
		// static page to test out back and forth websocket connection
		app.Static("/", "./static/home.html")
	}

	// Start the communication hub
	go web.RunHub() // on a separate goroutine|thread

	// Websocket setup
	app.Use(web.SetupWebsocketUpgrade())
	app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}
