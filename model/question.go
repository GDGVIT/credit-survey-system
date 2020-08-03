package model

type Question struct {
	QuestionId string `bson:"_id"`
	Body       string
	Type       string
	Multipart  []Section
}
