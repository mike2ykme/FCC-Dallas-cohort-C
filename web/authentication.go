package web

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"teamC/Global"
	"time"
)

func GetJwtMiddleware(cfg *Global.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.JwtSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println("jwt error handler called, returning 404 to user-- ", err)
			return c.SendStatus(fiber.StatusNotFound)
		},
	})
}

func ProductionLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		return nil
	}
}

const hoursInWeek = 168

func SimulatedLoginHandler(cfg *Global.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create the Claims
		claims := jwt.MapClaims{
			"username":  "John Doe",
			"firstName": "John",
			"lastName":  "Doe",
			"id":        uint(1),
			"admin":     true,
			"exp":       time.Now().Add(time.Hour * hoursInWeek).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(cfg.JwtSecret))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
