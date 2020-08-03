package model

type FormAudience struct {
	Id    string `bson:"_id"`
	Value string
}

func NewFormAudience(id string, value string) *FormAudience {
	return &FormAudience{Id: id, Value: value}
}
