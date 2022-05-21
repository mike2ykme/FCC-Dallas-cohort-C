package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"teamC/Global"
	"teamC/web"
)

func main() {
	cfg := Global.Configuration{}

	err := web.LoadConfiguration(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.WebApp = fiber.New()
	app := cfg.WebApp

	app.Use(logger.New())

	if cfg.Production {
		web.ProductionConfiguration(&cfg)
	} else {
		web.NonProductionConfiguration(&cfg)
		// static page to test out back and forth websocket connection
		app.Static("/", "./static/home.html")
	}
	// Authentication middleware
	app.Use(web.GetJwtMiddleware(&cfg))

	// Start the communication hub
	go web.RunHub() // on a separate goroutine|thread

	web.SetupAPIRoutes(&cfg)

	websockets := app.Group("/ws")
	websockets.Use(web.SetupWebsocketUpgrade())
	websockets.Get("/:id", web.WebsocketRoom())

	//app.Use(web.SetupWebsocketUpgrade())
	//app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}
