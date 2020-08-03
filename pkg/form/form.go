package form

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateForms(NewForm model.Form) (string, error) {
	client := utils.GetClient()
	forms := client.Database("Main").Collection("Forms")

	res, err := forms.InsertOne(context.TODO(), NewForm)
	if err != nil {
		Errorln(err.Error())
		return "", err
	}
	return fmt.Sprintf("%s", res.InsertedID), nil
}

func UpdateForm() {
	client := utils.GetClient()
	_id := "da3d44a6-70ac-47f5-8bc2-7cb998b93d8d"
	var m = make(map[string]interface{})
	m["publicDash"] = true
	m["description"] = "NEW FORM! Awesome. (:->)"
	forms := client.Database("Main").Collection("Forms")
	res, err := forms.UpdateOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": _id,
		},
	}, bson.M{
		"$set": m,
	})
	if err != nil {
		Errorln(err.Error())
		return
	}
	Debugln(res)
}

func GetFormDetails(id string) (model.Form, error) {
	client := utils.GetClient()
	forms := client.Database("Main").Collection("Forms")
	curr := forms.FindOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})
	if curr.Err() != nil {
		Errorln(curr.Err().Error())
		return model.Form{}, curr.Err()
	}
	var m model.Form
	curr.Decode(&m)
	return m, nil
}

func SetPublishStatus(id string, status bool) error {
	client := utils.GetClient()
	var m = make(map[string]interface{})
	m["isPublished"] = status
	forms := client.Database("Main").Collection("Forms")
	_ , err := forms.UpdateOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	}, bson.M{
		"$set": m,
	})
	if err != nil {
		Errorln(err.Error())
		return err
	}
	return nil
}
