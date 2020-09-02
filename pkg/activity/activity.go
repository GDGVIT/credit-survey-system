package activity

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"github.com/google/uuid"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TransactionType int8

const (
	Responded TransactionType = 0
	Created	= 1
	ReverseTransfer = 2
	Others = 3
)

func AddActivity(NewActivity model.Activity,typeOfTransaction TransactionType, UserID string) error {
	client := utils.GetClient()
	client.Database("Main").Collection("ActivityLogs")
	activity := client.Database("Main").Collection("ActivityLogs")
	SurveysMade := 0
	SurveysTaken := 0

	if typeOfTransaction == Responded {
		SurveysTaken += 1
	} else if typeOfTransaction == Created {
		SurveysMade += 1
	}

	status, err := activity.UpdateOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": UserID,
		},
	}, bson.M{
		"$push": bson.M{
			"activities": NewActivity,
		},
		"$inc": bson.M{
			"credits": NewActivity.TransactionAmt,
			"surveysMade": SurveysMade,
			"SurveysTaken": SurveysTaken,
		},

	})
	if err != nil {
		Errorln(err.Error())
		return err
	}
	Infoln(status)
	return nil
}

func CreateNewUserActivity(initialCredit int64, UserID string) error {
	client := utils.GetClient()  // GET the client

	// CONNECTION to activity collection
	activity := client.Database("Main").Collection("ActivityLogs")

	// Initial CREDITS ACTIVITY credit
	newActivity := model.Activity{
		ActivityId:     uuid.New().String(),
		Description:    "Starting account credit.",
		TransactionAmt: float64(initialCredit),
		FormId:         "NONE",
		Timestamp:      time.Now(),
	}

	// New activity logs instance created
	newActivityLog := model.ActivityLog{
		UserId: UserID,
		Credits: float64(initialCredit),
		SurveysMade: 0,
		SurveysTaken: 0,
		Activities: []model.Activity{
			newActivity,
		},
	}
	_, err := activity.InsertOne(context.TODO(), newActivityLog)
	if err != nil {
		Errorln(err.Error())
		return err
	}
	return nil
}

func GetActivityLogs(id string) (model.ActivityLog, error) {
	client := utils.GetClient()
	client.Database("Main").Collection("ActivityLogs")
	activity := client.Database("Main").Collection("ActivityLogs")

	curr := activity.FindOne(context.TODO(), bson.M{
		"_id": id,
	})
	if curr.Err() != nil {
		return model.ActivityLog{}, curr.Err()
	}
	var m model.ActivityLog
	curr.Decode(&m)
	return m, nil
}

func GetCredits(id string) (float64, error) {
	client := utils.GetClient()
	client.Database("Main").Collection("ActivityLogs")
	activity := client.Database("Main").Collection("ActivityLogs")

	curr := activity.FindOne(context.TODO(), bson.M{
		"_id": id,
	}, options.FindOne().SetProjection(bson.M{
		"credits": 1,
	}))
	if curr.Err() != nil {
		return 0, curr.Err()
	}
	var m map[string] interface{}
	curr.Decode(&m)
	return m["credits"].(float64), nil
}
func Errorln(args ...interface{}) {
	color.Error.Println(args...)
}
func Debugln(args ...interface{}) {
	color.Debug.Println(args...)
}
func Infoln(args ...interface{}) {
	color.Info.Println(args...)
}
