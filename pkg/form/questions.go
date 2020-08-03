package form

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/bson"
)

const dbURL = "mongodb+srv://Admin:look202020@cluster0-fbsrx.gcp.mongodb.net/<dbname>?retryWrites=true&w=majority"

func PostQuestions(question model.FormQuestion) error {
	client := utils.GetClient()
	formQuest := client.Database("Main").Collection("FormQuestions")
	_, err := formQuest.InsertOne(context.TODO(), question)
	if err != nil {
		Errorln(err.Error())
		return err
	}
	return nil
}

func ReplaceQuestions(NewFormQuest model.FormQuestion) error {
	client := utils.GetClient()
	formQuest := client.Database("Main").Collection("FormQuestions")
	_, err := formQuest.ReplaceOne(context.TODO(),
		bson.M{
			"_id": bson.M{
				"$eq": NewFormQuest.FormId,
			},
		}, NewFormQuest)

	if err != nil {
		Errorln(err.Error())
		return err
	}
	return nil
}

func GetQuestions(_id string) (model.FormQuestion, error) {
	client := utils.GetClient()
	formQuestions := client.Database("Main").Collection("FormQuestions")
	curr := formQuestions.FindOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": _id,
		},
	})
	if curr.Err() != nil {
		Errorln(curr.Err().Error())
		return model.FormQuestion{}, curr.Err()
	}
	var m model.FormQuestion
	curr.Decode(&m)

	return m, nil
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
