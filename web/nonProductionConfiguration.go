package web

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"teamC/Global"
	"teamC/db/rdbms"
)

func NonProductionConfiguration(cfg *Global.Configuration) error {
	repo, err := rdbms.NewRdbmsRepository(cfg.DatabaseURL, "sqlite")
	if cfg.AutoMigrate {
		repo.AutoMigrate()
	}
	if err != nil {
		return err
	}

	//deck := models.Deck{
	//	Description: "Test Deck 1",
	//	FlashCards: []models.FlashCard{
	//		{
	//			Question: "Test Question 1",
	//			Answers: []models.Answer{
	//				{
	//					Name:        "Test Answer 1",
	//					Value:       "",
	//					IsCorrect:   true,
	//					FlashCardId: 0,
	//				},
	//				{
	//					Name:        "Test Answer 2",
	//					Value:       "",
	//					IsCorrect:   false,
	//					FlashCardId: 0,
	//				},
	//			},
	//		},
	//	},
	//	OwnerId: 1,
	//}
	//deck2 := models.Deck{
	//	Description: "Test Deck 2",
	//	FlashCards: []models.FlashCard{
	//		{
	//			Question: "Test Question 1",
	//			Answers: []models.Answer{
	//				{
	//					Name:        "Test Answer 1",
	//					Value:       "",
	//					IsCorrect:   true,
	//					FlashCardId: 0,
	//				},
	//				{
	//					Name:        "Test Answer 2",
	//					Value:       "",
	//					IsCorrect:   false,
	//					FlashCardId: 0,
	//				},
	//			},
	//		},
	//	},
	//	OwnerId: 1,
	//}
	//if _, err := repo.SaveDeck(&deck); err != nil {
	//	cfg.Logger.Fatalf(err.Error())
	//
	//}
	//if _, err := repo.SaveDeck(&deck2); err != nil {
	//	cfg.Logger.Fatalf(err.Error())
	//}

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
