package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/Global"
	"teamC/db/inMemory"
	"teamC/models"
)

func NonProductionConfiguration(cfg *Global.Configuration) {
	// in-memory db setup
	repo := inMemory.NewInMemoryRepository()
	cfg.UserRepo = repo
	cfg.DeckRepo = repo
	cfg.FlashcardRepo = repo
	cfg.AnswerRepo = repo

	deck := models.Deck{
		Description: "Test Deck 1",
		FlashCards: []models.FlashCard{
			{
				Question: "Test Question 1",
				Answers: []models.Answer{
					{
						Name:        "Test Answer 1",
						Value:       "",
						IsCorrect:   true,
						FlashCardId: 0,
					},
					{
						Name:        "Test Answer 2",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
		OwnerId: 1,
	}
	_, err := repo.SaveDeck(&deck)
	if err != nil {
		cfg.Logger.Fatalf(err.Error())
	}
	app := cfg.WebApp
	app.Use(cors.New())
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config

	app.Post("/login/:id?", SimulatedLoginHandler(cfg))
}
