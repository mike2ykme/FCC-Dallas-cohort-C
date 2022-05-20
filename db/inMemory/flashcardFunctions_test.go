package inMemory

import (
	"teamC/models"
	"testing"
)

/*
type FlashcardRepository interface {
	SaveFlashcard(*models.FlashCard) (uint, error)
	GetFlashcardById(*models.FlashCard, uint) error
	GetAllFlashcardByDeckId(*[]models.FlashCard, uint) error
	GetAllFlashcards(*[]models.FlashCard) error
}

*/

func TestRepository_SaveFlashcard(t *testing.T) {
	newCard := models.FlashCard{
		Id:       0,
		Question: "",
		DeckId:   1,
		Answers:  []models.Answer{},
	}

	id, err := NewInMemoryRepository().SaveFlashcard(&newCard)

	if id != 1 || err != nil {
		t.Fatalf("expected the ID to be 1 and err to be nil, instead got %d & %#v", id, err)
	}
}

func TestRepository_SaveFlashcardInvalidDeckId(t *testing.T) {
	newCard := models.FlashCard{
		DeckId: 0,
	}

	id, err := NewInMemoryRepository().SaveFlashcard(&newCard)

	if id == 1 || err == nil {
		t.Fatalf("expected call to return id of 0 and error, instead got %d & %#v", id, err)
	}
}

func TestRepository_GetFlashcardById(t *testing.T) {
	repo := NewInMemoryRepository()
	repo.SaveFlashcard(&models.FlashCard{DeckId: 1})
	var newCard models.FlashCard
	err := repo.GetFlashcardById(&newCard, 1)

	if newCard.Id != 1 || err != nil {
		t.Fatalf("Expected id of 1 and erro to be nil, instead got %d and %#v", newCard.Id, err)
	}

}

/*
func TestRepository_GetFlashcardById(t *testing.T) {

}
*/