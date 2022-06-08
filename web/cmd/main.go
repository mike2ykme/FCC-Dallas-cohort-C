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

	const UserId = "userId"
	const FirstName = "firstName"
	websockets := app.Group("/ws", func(c *fiber.Ctx) error {
		if user, ok := c.Locals("user").(*jwt.Token); ok {
			claims, OK := user.Claims.(jwt.MapClaims)
			if OK {
				if userId, exists := claims["id"].(float64); exists {
					cfg.Logger.Println("found id")
					c.Locals(UserId, uint(userId))
				} else {
					cfg.Logger.Println("unable to get a user's ID")
					return c.SendStatus(fiber.StatusInternalServerError)
				}

				if username, exists := claims[FirstName].(string); exists {
					cfg.Logger.Println("found username")
					c.Locals(FirstName, username)
				} else {
					cfg.Logger.Println("unable to get a user's first name")
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				return c.Next()
			}
		}
		cfg.Logger.Println("unable to get the jwt token so handing back a 400")
		return c.SendStatus(fiber.StatusBadRequest)
	})
	websockets.Use(web.SetupWebsocketUpgrade())
	websockets.Get("/:id", web.WebsocketRoom(&cfg))

	//app.Use(web.SetupWebsocketUpgrade())
	//app.Get("/ws/:id", web.WebsocketRoom())

	// Start the web server
	log.Fatalln(app.Listen(cfg.Port))
}
