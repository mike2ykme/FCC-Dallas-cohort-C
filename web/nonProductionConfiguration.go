package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/Global"
	"teamC/db/inMemory"
)

func NonProductionConfiguration(cfg *Global.Configuration) {
	// in-memory db setup
	repo := inMemory.NewInMemoryRepository()
	cfg.UserRepo = repo
	cfg.DeckRepo = repo
	cfg.FlashcardRepo = repo
	cfg.AnswerRepo = repo

	app := cfg.WebApp
	app.Use(cors.New())
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config

	app.Post("/login/:id?", SimulatedLoginHandler(cfg))
}
