package model

import (
	"github.com/google/uuid"
	"time"
)

type Activity struct {
	ActivityId     string    `bson:"_id"`
	Description    string    `bson:"description"`
	TransactionAmt float64    `bson:"transactionAmt"`
	FormId         string    `bson:"formId"`
	Timestamp      time.Time `bson:"timestamp"`
}

func NewActivity(description string, transactionAmt float64, formId string) *Activity {
	return &Activity{
		ActivityId:     uuid.New().String(),
		Description:    description,
		TransactionAmt: transactionAmt,
		FormId:         formId,
		Timestamp:      time.Now(),
	}
}
