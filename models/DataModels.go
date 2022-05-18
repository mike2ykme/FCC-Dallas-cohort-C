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

type FlashCard struct {
	Id       uint
	Question string
	Answers  []Answer
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
	Id        uint
	Name      string
	Value     string
	IsCorrect bool
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
