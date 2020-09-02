package model

type ActivityLog struct {
	UserId       string     `bson:"_id"`
	Credits      float64      `bson:"credits";json:"credits"`
	SurveysMade  int64      `bson:"surveysMade";json:"surveysMade"`
	SurveysTaken int64      `bson:"surveysTaken";json:"surveysTaken"`
	Activities   []Activity `bson:"activities"`
}

