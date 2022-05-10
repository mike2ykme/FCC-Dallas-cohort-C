package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupWebsocketUpgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// if upgrading to websocket connection then continue
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		// else we're returning to the client saying we're expecting to have to upgrade connection
		return c.SendStatus(fiber.StatusUpgradeRequired)
	}
}
