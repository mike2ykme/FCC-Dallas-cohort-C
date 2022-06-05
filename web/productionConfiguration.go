package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"teamC/db/rdbms"

	"teamC/Global"
	"time"
)

func ProductionConfiguration(cfg *Global.Configuration) error {
	{
		repo, err := rdbms.NewRdbmsRepository(cfg.DatabaseURL)
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
	}

	app := cfg.WebApp
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.LimiterConfig.Max,
		Expiration: cfg.LimiterConfig.ExpirationSeconds * time.Second,
	}))

	app.Post("/login", ProductionLoginHandler(cfg))
	return nil
}
