package inMemory

import (
	"teamC/models"
	"testing"
)

/*

type AnswerRepository interface {
	SaveAnswer(answer *models.Answer) (uint, error)
	GetAnswerById(answer *models.Answer) error
	GetAnswersByFlashcardId(answers *[]models.Answer) error
	GetAllAnswers(answers *[]models.Answer) error
}

*/

func TestRepository_SaveAnswer(t *testing.T) {
	repo := NewInMemoryRepository()
	a := models.Answer{
		//ID:          0,
		Name:        "testName",
		Value:       "testValue",
		IsCorrect:   false,
		FlashCardId: 1,
	}
	_, err := repo.SaveAnswer(&a)

	if a.ID != 1 || err != nil {
		t.Fatalf("Expected ID to be 1 but ID was %d", a.ID)
	}
}
func TestRepository_SaveAnswerWithInvalidFlashcardId(t *testing.T) {
	_, err := NewInMemoryRepository().SaveAnswer(&models.Answer{})
	if err == nil {
		t.Fatalf("We should have had an error when saving a flashcard id of 0")
	}
}
func TestRepository_GetAnswerById(t *testing.T) {
	repo := NewInMemoryRepository()

	a := models.Answer{
		//ID:          0,
		Name:        "testName",
		Value:       "testValue",
		IsCorrect:   false,
		FlashCardId: 1,
	}
	repo.SaveAnswer(&a)
	var b models.Answer

	repo.GetAnswerById(&b, 1)

	if !a.IsEqual(&b) {
		t.Fatalf("expected a and b to be the same but they aren't \n\n a-> %#v \n\n b-> %#v", a, b)
	}
}

func TestRepository_GetAnswersByFlashcardId(t *testing.T) {
	repo := NewInMemoryRepository()
	all := make([]models.Answer, 0)

	a := models.Answer{FlashCardId: 1}
	b := models.Answer{FlashCardId: 1}
	repo.SaveAnswer(&a)
	repo.SaveAnswer(&b)
	repo.GetAnswersByFlashcardId(&all, 1)

	if len(all) != 2 || a.ID != 1 || b.ID != 2 {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v", all)
	}

}

func TestRepository_GetAllAnswers(t *testing.T) {
	repo := NewInMemoryRepository()
	all := make([]models.Answer, 0)

	a := models.Answer{FlashCardId: 1}
	b := models.Answer{FlashCardId: 1}

	repo.SaveAnswer(&a)
	repo.SaveAnswer(&b)
	repo.GetAllAnswers(&all)

	if len(all) != 2 || a.ID != 1 || b.ID != 2 {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v", all)
	}
}
