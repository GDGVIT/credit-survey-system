package sanitizer

import (
	"CreditBasedSurvey/model"
	"github.com/google/uuid"
	)

func SanitizeQuestions(quests *model.FormQuestion, id string) bool {
	(*quests).FormId = id
	(*quests).NumQuest = int64(len((*quests).Questions))
	for i, v := range (*quests).Questions {
		(*quests).Questions[i].QuestionId = uuid.New().String()
		if v.Type == "" || v.Body == "" {
			return true
		}
	}
	if (*quests).NumQuest == 0 {
		return true
	}
	return false
}

