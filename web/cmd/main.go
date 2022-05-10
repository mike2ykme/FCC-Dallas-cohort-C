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
)

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func main() {
	cfg := models.Configuration{}

	fatalIfErr(web.LoadConfiguration(&cfg))

	cfg.WebApp = fiber.New()
	//web.MiddlewareSetup(&cfg)

	cfg.WebApp.Use(logger.New())
	cfg.WebApp.Use(cors.New())

	cfg.WebApp.Post("/login", web.LoginHandler(&cfg))

	cfg.WebApp.Use(web.GetJwtMiddleware(&cfg))

	if cfg.Production {
		// rate limiting
		cfg.WebApp.Use(limiter.New(limiter.Config{
			Max:        cfg.LimiterConfig.Max,
			Expiration: cfg.LimiterConfig.ExpirationSeconds,
		}))
	} else {
		// performance monitoring w/ page
		cfg.WebApp.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config
	}

	// Setup Routes
	if !cfg.Production {
		// static page to test out back and forth websocket connection
		cfg.WebApp.Static("/", "./static/home.html")
	}
	//cfg.WebApp.Post("/login", web.LoginHandler(&cfg))

	// Start the communication hub
	go web.RunHub()

	// Websocket setup
	cfg.WebApp.Use(web.SetupWebsocketUpgrade())
	cfg.WebApp.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(cfg.WebApp.Listen(cfg.Port))
}
