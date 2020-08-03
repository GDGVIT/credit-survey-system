package model

type Response struct {
	FormId  string `bson:"formId"`
	UserId  string `bson:"_id"`
	Answers []Answer
}

func NewResponse() *Response {
	return &Response{}
}
