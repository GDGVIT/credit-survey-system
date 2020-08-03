package sanitizer

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/exchangeLogic"
	"github.com/google/uuid"
)

func SanitizeFormDetails(userID string,m *model.Form ) {
	m.FormId = uuid.New().String()
	m.UserId = userID
	m.NumResp = 0
	m.GainRate = exchangeLogic.GetGainRate()
	m.ResponseRate = exchangeLogic.GetResponseRate()
	m.IsPublished = false
}


func CheckFormDetailsPayload(form *model.Form) bool {
	if (*form).Title == "" || form.Description == "" || form.CreditsAllotted <= 0 {
		return true
	}
	return false
}
