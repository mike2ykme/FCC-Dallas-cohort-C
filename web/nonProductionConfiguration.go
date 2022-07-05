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
				Question: "Test Question 1-1",
				Answers: []models.Answer{
					{
						Name:        "Test Answer 1-1-1",
						Value:       "This would be what the answer would contain 111",
						IsCorrect:   true,
						FlashCardId: 0,
					},
					{
						Name:        "Test Answer 1-1-2",
						Value:       "This would be what the answer would contain 112",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
			{
				Question: "Test Question 1-2",
				Answers: []models.Answer{
					{
						Name:        "Test Answer 1-2-1",
						Value:       "This would be what the answer would contain 121",
						IsCorrect:   true,
						FlashCardId: 0,
					},
					{
						Name:        "Test Answer 1-2-2",
						Value:       "This would be what the answer would contain 122",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
		OwnerId: 1,
	}
	deck2 := models.Deck{
		Description: "Test Deck 2",
		FlashCards: []models.FlashCard{
			{
				Question: "Test Question 2-1",
				Answers: []models.Answer{
					{
						Name:        "Test Answer 2-1",
						Value:       "This would be what the answer would contain 21",
						IsCorrect:   true,
						FlashCardId: 0,
					},
					{
						Name:        "Test Answer 2-2",
						Value:       "This would be what the answer would contain 22",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
		OwnerId: 1,
	}
	if _, err := repo.SaveDeck(&deck); err != nil {
		cfg.Logger.Fatalf(err.Error())

	}
	if _, err := repo.SaveDeck(&deck2); err != nil {
		cfg.Logger.Fatalf(err.Error())
	}

	app := cfg.WebApp
	app.Use(cors.New())
	// performance monitoring w/ page
	app.Get("/monitor", monitor.New()) // monitor.Config{APIOnly: true} // optional config

	app.Post("/login/:id?", SimulatedLoginHandler(cfg))
}
