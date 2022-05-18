package web

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/Global"
	"teamC/db/inMemory"
)

func NonProductionConfiguration(cfg *Global.Configuration) {
	cfg.UserRepo = inMemory.NewInMemoryRepository()
	app := cfg.WebApp
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config

	app.Post("/login", SimulatedLoginHandler(cfg))
}
