package rdbms

import (
	"teamC/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_SaveDeck(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Deck{}, "1=1")
	deck := models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	_, err := repo.SaveDeck(&deck)

	if err != nil {
		t.Fatalf("There should be no error when saving the deck, but recieved %#v", err)
	}
	for _, card := range deck.FlashCards {
		if card.DeckId == 0 {
			t.Fatalf("There should be no cards with a 0 deckId")
		}
		for _, answer := range card.Answers {
			if answer.FlashCardId != card.ID {
				t.Fatalf("These should be the same as the flashcard, but it is %d instead", answer.FlashCardId)
			}
		}

	}
}

func TestRepository_ModifyDeck(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Deck{}, "1=1")
	repo.DB.Delete(&models.FlashCard{}, "1=1")
	repo.DB.Delete(&models.Answer{}, "1=1")
	repo.DB.Delete(&models.User{}, "1=1")
	deck := models.Deck{
		Description: "old description",
		OwnerId: 1,
		FlashCards: []models.FlashCard{
			{
				Question: "old question",
				DeckId: 1,
				Answers: []models.Answer{
					{
						Name: "old name",
						Value: "old value",
						IsCorrect: false,
						FlashCardId: 1,
					},
				},
			},
		},
	}
	deck.ID = 1
	deck.FlashCards[0].ID = 1
	deck.FlashCards[0].Answers[0].ID = 1
	if _, err := repo.SaveDeck(&deck); err != nil {
		t.Fatalf("There should be no error when saving the deck, but recieved %#v", err)
	}

	newDeck := models.Deck{
		Description: "new description",
		OwnerId: 1,
		FlashCards: []models.FlashCard{
			{
				Question: "new question",
				DeckId: 1,
				Answers: []models.Answer{
					{
						Name: "new name",
						Value: "new value",
						IsCorrect: true,
						FlashCardId: 1,
					},
				},
			},
		},
	}
	newDeck.ID = 1
	newDeck.FlashCards[0].ID = 1
	newDeck.FlashCards[0].Answers[0].ID = 1
	if _, err := repo.SaveDeck(&newDeck); err != nil {
		t.Fatalf("There should be no error when saving the deck, but recieved %#v", err)
	}

	var retrievedDeck models.Deck
	repo.GetDeckById(&retrievedDeck, deck.ID)
	assert := assert.New(t)
	assert.Equal(uint(1), retrievedDeck.ID, "Deck IDs don't match. May have created a new deck instead of updating.")
	assert.Equal(uint(1), retrievedDeck.FlashCards[0].ID, "FlashCard IDs don't match.")
	assert.Equal(uint(1), retrievedDeck.FlashCards[0].Answers[0].ID, "Answer IDs don't match.")
	assert.Equal("new description", retrievedDeck.Description, "Deck description did not update correctly.")
	assert.Equal("new question", retrievedDeck.FlashCards[0].Question, "Question did not update correctly.")
	assert.Equal("new name", retrievedDeck.FlashCards[0].Answers[0].Name, "Answer name did not update correctly.")
	assert.Equal("new value", retrievedDeck.FlashCards[0].Answers[0].Value, "Answer value did not update correctly.")
	assert.Equal(true, retrievedDeck.FlashCards[0].Answers[0].IsCorrect, "Answer correctness did not update correctly.")
}

func TestRepository_GetDeckById(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Deck{}, "1=1")
	oldDeck := models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	_, err := repo.SaveDeck(&oldDeck)
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	var newDeck models.Deck
	repo.GetDeckById(&newDeck, oldDeck.ID)

	if len(newDeck.FlashCards) == 0 {
		t.Fatalf("there's no flashcards :(\n")
	}
	newFlashcard := newDeck.FlashCards[0]

	if len(newFlashcard.Answers) == 0 {
		t.Fatalf("there's no answers :(\n")
	}
	newAnswer := newFlashcard.Answers[0]

	if newFlashcard.DeckId != newDeck.ID ||
		newAnswer.FlashCardId != newFlashcard.ID {
		t.Fatalf("the sub element should have a reference to parent ID")
	}
	if newDeck.ID != oldDeck.ID {
		t.Fatalf("these decks should be equal, but they're not. Instead have \n\n -> %#v \n\n -> %#v", newDeck, oldDeck)
	}

}

func TestRepository_GetAllDecks(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Deck{}, "1=1")
	oldDeck := models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	_, err := repo.SaveDeck(&oldDeck)
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	_, err = repo.SaveDeck(&models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	_, err = repo.SaveDeck(&models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		//FlashCards:  nil,
	})
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	var all []models.Deck
	err = repo.GetAllDecks(&all)
	if err != nil {
		t.Fatalf("should be able to retrieve the decks without an error, but received %#v", err)
	}

	if len(all) != 3 || err != nil {
		t.Fatalf("expected to have 3 elements returned")
	}

	for _, deck := range all {
		if deck.ID == 0 {
			t.Fatalf("expected none of the Ids to be 0")
		}
	}
}

func TestRepository_GetAllDecksByUserId(t *testing.T) {
	repo := getRepo(t)
	repo.DB.Delete(&models.Deck{}, "1=1")
	oldDeck := models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	_, err := repo.SaveDeck(&oldDeck)
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	_, err = repo.SaveDeck(&models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		FlashCards: []models.FlashCard{
			{
				//ID:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						//ID:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	_, err = repo.SaveDeck(&models.Deck{
		//ID:          0,
		Description: "",
		OwnerId:     1,
		//FlashCards:  nil,
	})
	if err != nil {
		t.Fatalf("should be able to save the deck without err, but received %#v", err)
	}
	var all []models.Deck
	err = repo.GetDecksByUserId(&all, 1)
	if err != nil {
		t.Fatalf("should be able to retrieve the decks without an error, but received %#v", err)
	}

	if len(all) != 3 || err != nil {
		t.Fatalf("expected to have 3 elements returned")
	}

	for _, deck := range all {
		if deck.ID == 0 {
			t.Fatalf("expected none of the Ids to be 0")
		}
	}
}
