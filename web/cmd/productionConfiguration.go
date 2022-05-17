package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"log"
	"os"
	"teamC/Global"
	"time"
)

func productionConfiguration(cfg *Global.Configuration) {
	goth.UseProviders(
		google.New(os.Getenv("OAUTH_KEY"), os.Getenv("OAUTH_SECRET"), "http://127.0.0.1:8088/auth/callback/google"),
	)

	app := cfg.WebApp

	app.Use(cors.New())

	// rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.LimiterConfig.Max,
		Expiration: cfg.LimiterConfig.ExpirationSeconds * time.Second,
	}))

	app.Get("/login", func(c *fiber.Ctx) error {
		c.Set("provider", "google")
		return c.Redirect("/login/google", fiber.StatusTemporaryRedirect)
	})
	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)

	app.Get("/auth/callback/:provider", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Fatal(err)
		}

		return ctx.SendString(user.Email)
	})

	app.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		return ctx.SendString("logout")
	})
}
