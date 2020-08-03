package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func SetClient() (err error){
	URL :="mongodb+srv://Admin:LookAtMe2018@cluster0-fbsrx.gcp.mongodb.net/?retryWrites=true&w=majority"
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(URL))
	return err
}

func GetClient() *mongo.Client {
	return client
}
func CloseClient() error {
	return client.Disconnect(context.TODO())
}