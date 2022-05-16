package main

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/models"
	"teamC/web"
)

func nonProductionConfiguration(cfg *models.Configuration) {

	app := cfg.WebApp
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config

	app.Post("/login", web.SimulatedLoginHandler(cfg))
}
