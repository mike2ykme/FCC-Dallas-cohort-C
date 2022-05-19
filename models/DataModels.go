package models

type Deck struct {
	Id          uint
	Description string
	Cards       []FlashCard
}

func (d *Deck) CopyReferences(o *Deck) {
	d.Id = o.Id
	d.Description = o.Description

	if len(d.Cards) < len(o.Cards) && len(d.Cards) == 0 {
		d.Cards = make([]FlashCard, len(o.Cards))

		for idx, card := range o.Cards {
			d.Cards[idx] = card.Copy()
		}
	}

}

func (d *Deck) Copy() Deck {
	return Deck{
		Id:          d.Id,
		Description: d.Description,
		Cards:       d.Cards,
	}
}

type FlashCard struct {
	Id       uint
	Question string
	DeckId   uint
	Answers  []Answer
}

func (f *FlashCard) CopyRef(o *FlashCard) {
	f.Id = o.Id
	f.DeckId = o.DeckId
	f.Answers = o.Answers
	f.Question = o.Question
}

func (f *FlashCard) Copy() FlashCard {
	answers := make([]Answer, len(f.Answers))
	for idx, answer := range f.Answers {
		answers[idx] = answer.Copy()
	}
	return FlashCard{
		Id:       f.Id,
		Question: f.Question,
		Answers:  answers,
	}
}

type Answer struct {
	Id          uint
	Name        string
	Value       string
	IsCorrect   bool
	FlashCardId uint
}

func (a *Answer) Copy() Answer {
	return Answer{
		Id:        a.Id,
		Name:      a.Name,
		Value:     a.Value,
		IsCorrect: a.IsCorrect,
	}
}

func (a *Answer) CopyRef(o *Answer) {
	a.Name = o.Name
	a.Id = o.Id
	a.IsCorrect = o.IsCorrect
	a.Value = o.Value
}
