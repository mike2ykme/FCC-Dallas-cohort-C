package models

type Deck struct {
	Id          uint
	Description string
	Cards       []FlashCard
}

type FlashCard struct {
	Id       uint
	Question string
	Answers  []Answer
}

type Answer struct {
	Id        uint
	Name      string
	Value     string
	IsCorrect bool
}
