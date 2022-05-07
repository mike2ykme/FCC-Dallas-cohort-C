package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
	web.MiddlewareSetup(&cfg)

	// All other calls should be handled by websockets
	cfg.WebApp.Use(func(c *fiber.Ctx) error {
		// if upgrading to websocket connection then continue
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		// else we're returning to the client saying we're expecting to have to upgrade connection
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	go web.RunHub()

	cfg.WebApp.Get("/ws/:id", web.WebsocketRoom())

	log.Fatalln(cfg.WebApp.Listen(cfg.Port))
}
