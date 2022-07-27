package rdbms

import (
	"os"
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
func getURL() string {
	return os.Getenv("DB_URL")
}

func getDBType() string {
	return os.Getenv("DB_TYPE")
}

func getRepo(t *testing.T) *repository {
	r, err := NewRdbmsRepository(getURL(), getDBType())
	if err != nil {
		t.Fatalf("unable to get a db connection %#v \n", err)
		return nil
	}
	r.DB.AutoMigrate(&models.User{})
	r.DB.AutoMigrate(&models.Deck{})
	r.DB.AutoMigrate(&models.FlashCard{})
	r.DB.AutoMigrate(&models.Answer{})

	return r
}

func TestRepository_SaveAnswer(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Answer{}, "1=1")
	a := models.Answer{
		//ID:          0,
		Name:        "testName",
		Value:       "testValue",
		IsCorrect:   false,
		FlashCardId: 1,
	}
	_, err := repo.SaveAnswer(&a)

	if a.ID == 0 || err != nil {
		t.Fatalf("Expected ID to be 1 but ID was %d", a.ID)
	}
}
func TestRepository_SaveAnswerWithInvalidFlashcardId(t *testing.T) {
	_, err := getRepo(t).SaveAnswer(&models.Answer{})
	if err == nil {
		t.Fatalf("We should have had an error when saving a flashcard id of 0")
	}
}

func TestRepository_GetAnswerById(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Answer{}, "1=1")
	//repo := NewInMemoryRepository()

	a := models.Answer{
		//ID:          0,
		Name:        "testName",
		Value:       "testValue",
		IsCorrect:   false,
		FlashCardId: 1,
	}
	repo.SaveAnswer(&a)
	var b models.Answer

	repo.GetAnswerById(&b, a.ID)

	if !a.IsEqual(&b) {
		t.Fatalf("expected a and b to be the same but they aren't \n\n a-> %#v \n\n b-> %#v", a, b)
	}
}

func TestRepository_GetAnswersByFlashcardId(t *testing.T) {
	//repo := NewInMemoryRepository()
	repo := getRepo(t)
	repo.DB.Delete(&models.Answer{}, "1=1")
	all := make([]models.Answer, 0)

	a := models.Answer{FlashCardId: 3}
	b := models.Answer{FlashCardId: 3}

	idA, err := repo.SaveAnswer(&a)
	if err != nil {
		t.Fatalf("there was an error: %#v", err)
	}
	idB, err := repo.SaveAnswer(&b)

	if err != nil {
		t.Fatalf("there was an error: %#v", err)
	}

	if err := repo.GetAnswersByFlashcardId(&all, 3); err != nil {
		t.Fatalf("there was an error: %#v", err)
	}

	if len(all) != 2 || a.ID != idA || b.ID != idB {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v", all)
	}

}

func TestRepository_GetAllAnswers(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Answer{}, "1=1")
	//repo := NewInMemoryRepository()
	all := make([]models.Answer, 0)

	a := models.Answer{FlashCardId: 1}
	b := models.Answer{FlashCardId: 2}

	idA, err := repo.SaveAnswer(&a)
	if err != nil {
		t.Fatalf("there was an error: %#v", err)
	}
	idB, err := repo.SaveAnswer(&b)
	if err != nil {
		t.Fatalf("there was an error: %#v", err)
	}
	repo.GetAllAnswers(&all)

	if len(all) != 2 || a.ID != idA || b.ID != idB {
		t.Fatalf("expected to have two answers in the slice with ids of 1 and 2, but have \n\n %#v \n\n count->%d", all, len(all))
	}
}
