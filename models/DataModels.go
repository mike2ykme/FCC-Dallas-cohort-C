package models

import (
	"gorm.io/gorm"
	"reflect"
    "math/rand"
    "time"
)

type Deck struct {
	gorm.Model
	//ID          uint
	Description string
	FlashCards  []FlashCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OwnerId     uint
}

func (d *Deck) Shuffle() {
    deck := d  // For some reason, I have to to this for the anonymous function below to recognize the deck
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(d.FlashCards), func(i, j int) {
        deck.FlashCards[i], deck.FlashCards[j] = deck.FlashCards[j], deck.FlashCards[i]
    })
}

func (d *Deck) CopyReferences(o *Deck) {
	d.ID = o.ID
	d.Description = o.Description
	d.OwnerId = o.OwnerId

	if len(d.FlashCards) < len(o.FlashCards) && len(d.FlashCards) == 0 {
		d.FlashCards = make([]FlashCard, len(o.FlashCards))

		for idx, card := range o.FlashCards {
			d.FlashCards[idx] = card.Copy()
		}
	}

}

func (d *Deck) Copy() Deck {
	return Deck{
		Model:       d.Model,
		Description: d.Description,
		FlashCards:  d.FlashCards,
		OwnerId:     d.OwnerId,
	}
}

func (d *Deck) IsEqualTo(o *Deck) bool {
	return reflect.DeepEqual(d, o)
}

func (d *Deck) ReplaceFields(o *Deck) {
	d.Model = o.Model
	d.Description = o.Description
	d.FlashCards = o.FlashCards
	d.OwnerId = o.OwnerId
}

type FlashCard struct {
	gorm.Model
	//ID       uint
	Question string
	DeckId   uint
	Answers  []Answer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (f *FlashCard) CopyRef(o *FlashCard) {
	f.Model = o.Model
	f.Question = o.Question
	f.DeckId = o.DeckId
	f.Answers = o.Answers
}

func (f *FlashCard) Copy() FlashCard {
	answers := make([]Answer, len(f.Answers))
	for idx, answer := range f.Answers {
		answers[idx] = answer.Copy()
	}
	return FlashCard{
		Model:    f.Model,
		Question: f.Question,
		Answers:  answers,
		DeckId:   f.DeckId,
	}
}

type Answer struct {
	gorm.Model
	//ID   uint
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsCorrect   bool   `json:"isCorrect"`
	FlashCardId uint   `json:"flashCardId "`
}

func (a *Answer) Copy() Answer {
	return Answer{
		Model: a.Model,
		//ID:          a.ID,
		Name:        a.Name,
		Value:       a.Value,
		IsCorrect:   a.IsCorrect,
		FlashCardId: a.FlashCardId,
	}
}

func (a *Answer) CopyRef(o *Answer) {
	a.Name = o.Name
	//a.ID = o.ID
	a.Model = o.Model
	a.IsCorrect = o.IsCorrect
	a.Value = o.Value
	a.FlashCardId = o.FlashCardId
}

func (a *Answer) IsEqual(o *Answer) bool {
	return a.ID == o.ID &&
		a.Name == o.Name &&
		a.IsCorrect == o.IsCorrect &&
		a.Value == o.Value &&
		a.FlashCardId == o.FlashCardId
}
