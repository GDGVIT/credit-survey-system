package model

import "github.com/google/uuid"

type User struct {
	Name         string `bson:"name";json:"name"`
	UserId       string `bson:"_id";json:"userId"`
	Email        string `bson:"email";json:"email"`
}

func NewUser(email string, name string) *User {
	return &User{
		Email:        email,
		UserId:       uuid.New().String(),
		Name:         name,
	}
}
