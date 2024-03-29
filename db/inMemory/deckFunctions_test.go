package inMemory

import (
	"teamC/models"
	"testing"
)

/*
type DeckRepository interface {
	SaveDeck(*models.Deck) (uint, error)
	GetDeckById(*models.Deck, uint) error
	GetAllDecks(*[]models.Deck) error
}
/*

*/

func TestRepository_SaveDeck(t *testing.T) {
	repo := NewInMemoryRepository()
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

	if len(repo.decks) == 0 {
		t.Fatalf("deck repo not updated")
	}

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

func TestRepository_GetDeckById(t *testing.T) {
	repo := NewInMemoryRepository()
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

	newFlashcard := newDeck.FlashCards[0]
	newAnswer := newFlashcard.Answers[0]

	if newFlashcard.DeckId != newDeck.ID ||
		newAnswer.FlashCardId != newFlashcard.ID {
		t.Fatalf("the sub element should have a reference to parent ID")
	}

	if !newDeck.IsEqualTo(&oldDeck) {
		t.Fatalf("these decks should be equal, but they're not. Instead have \n\n -> %#v \n\n -> %#v", newDeck, oldDeck)
	}
}

func TestRepository_GetAllDecks(t *testing.T) {
	repo := NewInMemoryRepository()
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
	repo := NewInMemoryRepository()
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
