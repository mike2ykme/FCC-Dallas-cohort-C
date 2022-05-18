package inMemory

//
//type memoryDeck struct {
//	Id          uint
//	Description string
//	Cards       []memoryFlashCard
//}
//
//func (d *memoryDeck) CopyReferences(o *memoryDeck) {
//	d.Id = o.Id
//	d.Description = o.Description
//
//	if len(d.Cards) < len(o.Cards) && len(d.Cards) == 0 {
//		d.Cards = make([]memoryFlashCard, len(o.Cards))
//
//		for idx, card := range o.Cards {
//			d.Cards[idx] = card.Copy()
//		}
//	}
//
//}
//
//type memoryFlashCard struct {
//	Id           uint
//	Question     string
//	Answers      []memoryAnswer
//	MemoryDeckId uint
//}
//
//func (f *memoryFlashCard) Copy() memoryFlashCard {
//	answers := make([]memoryAnswer, len(f.Answers))
//	for idx, answer := range f.Answers {
//		answers[idx] = answer.Copy()
//	}
//	return memoryFlashCard{
//		Id:       f.Id,
//		Question: f.Question,
//		Answers:  answers,
//	}
//}
//
//type memoryAnswer struct {
//	Id          uint
//	Name        string
//	Value       string
//	IsCorrect   bool
//	FlashCardId uint
//}
//
//func (a *memoryAnswer) Copy() memoryAnswer {
//	return memoryAnswer{
//		Id:          a.Id,
//		Name:        a.Name,
//		Value:       a.Value,
//		IsCorrect:   a.IsCorrect,
//		FlashCardId: a.FlashCardId,
//	}
//}
//
//func (a *memoryAnswer) CopyRef(o *memoryAnswer) {
//	a.Name = o.Name
//	a.Id = o.Id
//	a.IsCorrect = o.IsCorrect
//	a.Value = o.Value
//	a.FlashCardId = o.FlashCardId
//}
