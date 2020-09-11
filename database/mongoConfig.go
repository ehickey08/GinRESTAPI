package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Client *mongo.Client

func SetupMongoDatabase() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://dbAdmin:pathfind3r!@cluster0.jsr9u.mongodb.net/student_analytics?retryWrites=true&w=majority\n")
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
}
