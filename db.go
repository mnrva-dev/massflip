package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Client

func openDB() *mongo.Client {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + dbUsername + ":" + dbPassword + "@cluster0.tqrat.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func DBSubtractPoints(user string, p int) error {
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	var userCollection *mongo.Collection
	if os.Getenv("ENVIRONMENT") == "production" {
		userCollection = DB.Database("Users").Collection("Users")
	} else {
		userCollection = DB.Database("Development").Collection("Users")
	}

	var v ExistingAccount
	err = userCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: user}}).Decode(&v)

	if err != nil {
		fmt.Println("* Error editing points for user " + user + ": User not found in DB")
		return err
	}

	filter := bson.D{primitive.E{Key: "username", Value: user}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "points", Value: v.Points - p}}}}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println("* Error updating points in ")
	}
	return nil
}
func DBAddPoints(user string, p int) error {
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	var userCollection *mongo.Collection
	if os.Getenv("ENVIRONMENT") == "production" {
		userCollection = DB.Database("Users").Collection("Users")
	} else {
		userCollection = DB.Database("Development").Collection("Users")
	}
	var v ExistingAccount
	err = userCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: user}}).Decode(&v)

	if err != nil {
		fmt.Println("* Error editing points for user " + user + ": User not found in DB")
		return err
	}

	filter := bson.D{primitive.E{Key: "username", Value: user}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "points", Value: v.Points + p}}}}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println("* Error updating points in ")
	}
	return nil

}

func DBGetUserByUsername(user string) ExistingAccount {
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	var userCollection *mongo.Collection
	if os.Getenv("ENVIRONMENT") == "production" {
		userCollection = DB.Database("Users").Collection("Users")
	} else {
		userCollection = DB.Database("Development").Collection("Users")
	}
	var v ExistingAccount
	err = userCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: user}}).Decode(&v)

	if err != nil {
		fmt.Println("* Error getting user " + user + ": Not found in DB")
		var e ExistingAccount
		return e
	}

	return v
}
