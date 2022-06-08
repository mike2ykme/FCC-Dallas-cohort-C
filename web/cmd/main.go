package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v4"
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
		if err := web.ProductionConfiguration(&cfg); err != nil {
			cfg.Logger.Fatal(err)
		}
	} else {
		web.NonProductionConfiguration(&cfg)
		// static page to test out back and forth websocket connection
		app.Static("/", "./static/home.html")
	}

	web.Configs = &cfg

	// Since websockets don't support headers read it from the url and update the request header
	// to allow all requests to follow standard JWT middleware
	app.Use(func(c *fiber.Ctx) error {
		if token := c.Query("token", ""); token != "" {
			c.Request().Header.Set("Authorization", "Bearer "+token)
		}
		return c.Next()
	})

	// Authentication middleware
	app.Use(web.GetJwtMiddleware(&cfg))

	// Start the communication hub
	go web.RunHub() // on a separate goroutine|thread

	web.SetupAPIRoutes(&cfg)

	websockets := app.Group("/ws", func(c *fiber.Ctx) error {
		c.Accepts("json", "text")     // "json"
		c.Accepts("application/json") // "application/json"
		if user, ok := c.Locals("user").(*jwt.Token); ok {
			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				if userId, ok := claims["id"].(float64); ok {
					c.Locals("userId", uint(userId))
					return c.Next()
				}
			}
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	})
	websockets.Use(web.SetupWebsocketUpgrade())
	websockets.Get("/:id", web.WebsocketRoom(&cfg))

	//app.Use(web.SetupWebsocketUpgrade())
	//app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}
