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
		Id:          0,
		Name:        "testName",
		Value:       "testValue",
		IsCorrect:   false,
		FlashCardId: 1,
	}
	_, err := repo.SaveAnswer(&a)

	if a.Id != 1 || err != nil {
		t.Fatalf("Expected Id to be 1 but ID was %d", a.Id)
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
		Id:          0,
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

	repo.SaveAnswer(&models.Answer{FlashCardId: 1})
	repo.SaveAnswer(&models.Answer{FlashCardId: 1})
	repo.GetAnswersByFlashcardId(&all, 1)

	if len(all) != 2 {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v", all)
	}

	a := all[0]
	if (a.Id < 1 || a.Id > 2) || a.Name != "" || a.Value != "" || a.IsCorrect != false || a.FlashCardId != 1 {
		t.Fatalf("there was a problem with the value being loaded. It does not match expected results")
	}

	a = all[1]
	if (a.Id < 1 || a.Id > 2) || a.Name != "" || a.Value != "" || a.IsCorrect != false || a.FlashCardId != 1 {
		t.Fatalf("there was a problem with the value being loaded. It does not match expected results")
	}
}

func TestRepository_GetAllAnswers(t *testing.T) {
	repo := NewInMemoryRepository()
	all := make([]models.Answer, 0)

	repo.SaveAnswer(&models.Answer{FlashCardId: 1})
	repo.SaveAnswer(&models.Answer{FlashCardId: 1})
	repo.GetAllAnswers(&all)

	if len(all) != 2 || (all[0].Id < 1 || all[0].Id > 2) || (all[1].Id > 2 || all[1].Id < 2) {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v", all)
	}
	a := all[0]
	if (a.Id < 1 || a.Id > 2) || a.Name != "" || a.Value != "" || a.IsCorrect != false || a.FlashCardId != 1 {
		t.Fatalf("there was a problem with the value being loaded. It does not match expected results")
	}
}
