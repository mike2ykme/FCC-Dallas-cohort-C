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
		Id:          0,
		Description: "",
		FlashCards: []models.FlashCard{
			{
				Id:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						Id:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	repo.SaveDeck(&deck)

	for _, card := range deck.FlashCards {
		if card.DeckId == 0 {
			t.Fatalf("There should be no cards with a 0 deckId")
		}
		for _, answer := range card.Answers {
			if answer.FlashCardId != card.Id {
				t.Fatalf("These should be the same as the flashcard, but it is %d instead", answer.FlashCardId)
			}
		}

	}
}

func TestRepository_GetDeckById(t *testing.T) {
	repo := NewInMemoryRepository()
	oldDeck := models.Deck{
		Id:          0,
		Description: "",
		FlashCards: []models.FlashCard{
			{
				Id:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						Id:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	repo.SaveDeck(&oldDeck)
	var newDeck models.Deck
	repo.GetDeckById(&newDeck, oldDeck.Id)

	newFlashcard := newDeck.FlashCards[0]
	newAnswer := newFlashcard.Answers[0]

	if newFlashcard.DeckId != newDeck.Id ||
		newAnswer.FlashCardId != newFlashcard.Id {
		t.Fatalf("the sub element should have a reference to parent ID")
	}

	if !newDeck.IsEqualTo(&oldDeck) {
		t.Fatalf("these decks should be equal, but they're not. Instead have \n\n -> %#v \n\n -> %#v", newDeck, oldDeck)
	}
}

func TestRepository_GetAllDecks(t *testing.T) {
	repo := NewInMemoryRepository()
	oldDeck := models.Deck{
		Id:          0,
		Description: "",
		FlashCards: []models.FlashCard{
			{
				Id:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						Id:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	}
	repo.SaveDeck(&oldDeck)
	repo.SaveDeck(&models.Deck{
		Id:          0,
		Description: "",
		FlashCards: []models.FlashCard{
			{
				Id:       0,
				Question: "",
				DeckId:   0,
				Answers: []models.Answer{
					{
						Id:          0,
						Name:        "",
						Value:       "",
						IsCorrect:   false,
						FlashCardId: 0,
					},
				},
			},
		},
	})
	repo.SaveDeck(&models.Deck{
		Id:          0,
		Description: "",
		FlashCards:  nil,
	})
	var all []models.Deck
	err := repo.GetAllDecks(&all)

	if len(all) != 3 || err != nil {
		t.Fatalf("expected to have 3 elements returned")
	}

	for _, deck := range all {
		if deck.Id == 0 {
			t.Fatalf("expected none of the Ids to be 0")
		}
	}
}
