package response

import (
	"CreditBasedSurvey/model"
	"CreditBasedSurvey/pkg/utils"
	"context"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/bson"
)

const dbURL = "mongodb+srv://Admin:look202020@cluster0-fbsrx.gcp.mongodb.net/<dbname>?retryWrites=true&w=majority"

func PostResponse(ansPayload model.Response) error {
	client := utils.GetClient()

	answers := client.Database("Main").Collection("Answers")
	_, err := answers.InsertOne(context.TODO(), ansPayload)
	if err != nil {
		Errorln(err.Error())
		return err
	}
	return nil
}

func GetResponse(id string) (model.Response, error) {
	client := utils.GetClient()
	answers := client.Database("Main").Collection("Answers")
	curr := answers.FindOne(context.TODO(), bson.M{
		"_id": bson.M{
			"$eq": id,
		},
	})
	if curr.Err() != nil {
		return model.Response{}, curr.Err()
	}
	var m model.Response
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
