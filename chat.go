package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ColorRequest struct {
	Username string `json:"username" bson:"username"`
	Color    string `json:"color" bson:"color"`
}

func chatColor(w http.ResponseWriter, r *http.Request) {
	// prepare DB
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	userCollection := DB.Database("Users").Collection("Users")

	// decode PUT into v struct
	var v ColorRequest
	json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := bson.D{primitive.E{Key: "username", Value: strings.ToLower(v.Username)}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "color", Value: v.Color}}}}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return
	}
}
