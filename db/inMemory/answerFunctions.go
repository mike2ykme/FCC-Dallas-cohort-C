package inMemory

import "teamC/models"

/*

type AnswerRepository interface {
	SaveAnswer(answer *models.Answer) (uint, error)
	GetAnswerById(answer *models.Answer) error
	GetAnswersByQuestionId(answers *[]models.Answer) error
	GetAllAnswers(answers *[]models.Answer) error
}

*/

func (m *repository) SaveAnswer(answer *models.Answer) (uint, error) {
	//if answer.Id == 0 {
	//	answer.Id = m.currentHighestAnswerId
	//	m.currentHighestAnswerId++
	//} else if answer.Id > m.currentHighestAnswerId {
	//	m.currentHighestAnswerId = answer.Id + 1
	//}
	//m.answers[answer.Id] = answer

	return answer.Id, nil
}

func (m *repository) GetAnswerById(answer *models.Answer, id uint) error {
	//if val, ok := m.answers[id]; ok {
	//	answer.CopyRef(val)
	//}
	return nil
}
func (m *repository) GetAnswersByQuestionId(answers *[]models.Answer, id uint) error {
	//for _, answer := range m.answers {
	//	if an
	//}
	return nil
}
func (m *repository) GetAllAnswers(answers *[]models.Answer) error {
	return nil
}
