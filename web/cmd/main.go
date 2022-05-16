package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"teamC/models"
	"teamC/web"
)

func main() {
	cfg := models.Configuration{}

	err := web.LoadConfiguration(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.WebApp = fiber.New()
	app := cfg.WebApp

	app.Use(logger.New())

	if cfg.Production {
		productionConfiguration(&cfg)
	} else {
		nonProductionConfiguration(&cfg)
		// Setup Routes
		// static page to test out back and forth websocket connection
		app.Static("/", "./static/home.html")
	}
	// Authentication middleware
	app.Use(web.GetJwtMiddleware(&cfg))

	// Start the communication hub
	go web.RunHub() // on a separate goroutine|thread

	// Websocket setup
	app.Use(web.SetupWebsocketUpgrade())
	app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}
