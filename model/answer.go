package model

type Answer struct {
	QuestionId string `bson:"_id"`
	Value      string
	Type       string
}
