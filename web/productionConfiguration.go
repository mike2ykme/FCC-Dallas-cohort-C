package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"teamC/Global"
	"teamC/db/inMemory"
	"time"
)

func ProductionConfiguration(cfg *Global.Configuration) {
	{
		// TODO this needs to be removed when a production DB is setup
		// in-memory db setup
		repo := inMemory.NewInMemoryRepository()
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

}
