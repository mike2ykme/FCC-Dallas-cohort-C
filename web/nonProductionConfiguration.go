package web

import (
	"teamC/Global"
	"teamC/db/rdbms"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func NonProductionConfiguration(cfg *Global.Configuration) error {
	repo, err := rdbms.NewRdbmsRepository(cfg.DatabaseURL, "sqlite")
	if cfg.AutoMigrate {
		repo.AutoMigrate()
	}
	if err != nil {
		return err
	}
	cfg.UserRepo = repo
	cfg.DeckRepo = repo
	cfg.FlashcardRepo = repo
	cfg.AnswerRepo = repo
	app := cfg.WebApp
	app.Use(cors.New())
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config
	app.Post("/login/:id?", SimulatedLoginHandler(cfg))
	return nil
}
