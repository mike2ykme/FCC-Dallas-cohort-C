package rdbms

import (
	"teamC/models"
	"testing"
)

/*
type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, errorr)
	GetFlashcardById(*models.FlashCard, uint) errorr
	GetAllFlashcardByDeckId(*[]models.FlashCard, uint) errorr
	GetAllFlashcards(*[]models.FlashCard) errorr
}

*/

func TestRepository_SaveFlashcard(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.FlashCard{}, "1=1")
	newCard := models.FlashCard{
		//ID:       0,
		Question: "",
		DeckId:   1,
		Answers:  []models.Answer{},
	}

	id, err := repo.SaveFlashcard(&newCard)

	if id != newCard.ID || err != nil {
		t.Fatalf("expected the ID to be the same and err to be nil, instead got %d & %#v", id, err)
	}
}

func TestRepository_SaveFlashcardInvalidDeckId(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.FlashCard{}, "1=1")

	newCard := models.FlashCard{
		DeckId: 0,
	}

	id, err := repo.SaveFlashcard(&newCard)

	if id != 0 || err == nil {
		t.Fatalf("expected call to return id of 0 and errorr, instead got %d & %#v", id, err)
	}
}

func TestRepository_GetFlashcardById(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.FlashCard{}, "1=1")

	oldCard := &models.FlashCard{DeckId: 1}

	repo.SaveFlashcard(oldCard)
	var newCard models.FlashCard
	err := repo.GetFlashcardById(&newCard, oldCard.ID)

	if newCard.ID != oldCard.ID || err != nil {
		t.Fatalf("Expected id of 1 and error to be nil, instead got %d and %#v", newCard.ID, err)
	}

}

func TestRepository_GetAllFlashcardByDeckId(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.FlashCard{}, "1=1")
	oldCard := &models.FlashCard{DeckId: 1}
	repo.SaveFlashcard(oldCard)

	all := make([]models.FlashCard, 0)
	err := repo.GetAllFlashcardByDeckId(&all, oldCard.DeckId)

	if len(all) < 1 {
		t.Fatal("length of all is less than 1")
	}
	if len(all) != 1 || all[0].ID != oldCard.ID || err != nil {
		t.Fatalf("Expected id of 1 and error to be nil, instead got %d and %#v", all[0].ID, err)
	}
}

func TestRepository_GetAllFlashcards(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.FlashCard{}, "1=1")
	repo.SaveFlashcard(&models.FlashCard{DeckId: 1})
	repo.SaveFlashcard(&models.FlashCard{DeckId: 1})
	repo.SaveFlashcard(&models.FlashCard{DeckId: 1})
	var all []models.FlashCard
	err := repo.GetAllFlashcards(&all)

	if len(all) != 3 || err != nil {
		t.Fatalf("Expected length to be 3 and error to be nil, instead got %d and %#v", all[0].ID, err)
	}
}
