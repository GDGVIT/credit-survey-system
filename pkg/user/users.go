package user

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/activity"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/bson"
)


func Errorln(args ...interface{}) {
	color.Error.Println(args...)
}
func Debugln(args ...interface{}) {
	color.Debug.Println(args...)
}
func Infoln(args ...interface{}) {
	color.Info.Println(args...)
}

// CREATING NEW USERS
func InitializeUser(NewUser model.User, initCredit int64) (string, error) {
	documentID, err := AddUser(NewUser)
	if err != nil {
		Errorln(err.Error())
		return "", err
	}
	err = activity.CreateNewUserActivity(initCredit, documentID)
	if err != nil {
		Errorln(err.Error())
		return "", err
	}
	return documentID, nil
}


func AddUser(newUser model.User) (string, error) {
	client := utils.GetClient()
	users := client.Database("Main").Collection("Users")
	v, err := users.InsertOne(context.TODO(), newUser)
	if err != nil {
		Errorln("Error occurred! " + err.Error())
		return "", err
	}
	return v.InsertedID.(string), nil
}

func UpdateUser() {
	client := utils.GetClient()
	users := client.Database("Main").Collection("Users")
	updatedStatus, err := users.UpdateOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": "1cf35a9b-58c1-44b9-bf33-ebdf007f685b",
		},
	}, bson.M{
		"$set": bson.M{
			"surveysTaken": 19000,
		},
	})
	if err != nil {
		Errorln(err.Error())
		return
	}
	Infoln(updatedStatus)
}

func DeleteUser(id string) (int64, error) {
	client := utils.GetClient()
	users := client.Database("Main").Collection("Users")
	deleteResult, err := users.DeleteOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})
	if err != nil {
		Errorln(err.Error())
		return 0, err
	}
	Infoln("delete count", deleteResult.DeletedCount)
	return deleteResult.DeletedCount, nil
}

func GetUser(id string) (model.User, error) {
	client := utils.GetClient()
	users := client.Database("Main").Collection("Users")
	curr := users.FindOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})
	if curr.Err() != nil {
		Errorln(curr.Err().Error())
		return model.User{}, curr.Err()
	}
	var m model.User

	curr.Decode(&m)

	return m, nil
}
