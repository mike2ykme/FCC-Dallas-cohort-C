package web

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"log"
	"teamC/models"
)

func httpMethodBasedFilter(ctx *fiber.Ctx) bool {
	m := ctx.Method()
	if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
		return true
	}
	return false
}

func getJwtFilter(cfg *models.Configuration) func(*fiber.Ctx) bool {
	if cfg.Production {
		return nil

	} else {
		return httpMethodBasedFilter
	}
}

func GetJwtMiddleware(cfg *models.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println("jwt error handler called, returning 404 to user-- ", err)
			return c.SendStatus(fiber.StatusNotFound)
		},
		//Filter: getJwtFilter(cfg),
	})
}
