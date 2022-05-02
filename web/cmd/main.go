package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
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
	app := fiber.New()

	// Default middleware config
	app.Use(logger.New())

	// Routes
	if !cfg.Production {
		log.Println("using non-prod configurations")
		app.Static("/", "./static/home.html")
	}

	// All other calls should be handled by websockets
	app.Use(func(c *fiber.Ctx) error {
		// Returns true if the client requested upgrade to the WebSocket protocol
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go web.RunHub()

	app.Get("/ws/:id", web.WebsocketRoom())

	log.Fatalln(app.Listen(cfg.Port))
}
