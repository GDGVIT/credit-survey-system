package model

type FormQuestion struct {
	FormId    string `bson:"_id"`
	NumQuest  int64
	Questions []Question
}
