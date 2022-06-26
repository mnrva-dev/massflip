package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Login struct {
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`
	RememberMe bool   `json:"remember"`
}

type Session struct {
	Session string `json:"session" bson:"session"`
}

// ExistingAccount is a struct that mirrors how Users are stored in the Database
type ExistingAccount struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Case     string `json:"case" bson:"case"`
	Color    string `json:"color" bson:"color"`
	Points   int    `json:"points" bson:"points"`
	Resets   int    `json:"resets" bson:"resets"`
	Session  string `json:"session" bson:"session"`
}

type ReturnedAccount struct {
	Username string `json:"username" bson:"case"`
	Color    string `json:"color" bson:"color"`
	Points   int    `json:"points" bson:"points"`
	Resets   int    `json:"resets" bson:"resets"`
}

var (
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_]{3,24}$`)
)

func createAccount(w http.ResponseWriter, r *http.Request) {
	// prepare DB
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	userCollection := DB.Database("Users").Collection("Users")

	// var v contains POST credentials
	var v Login
	err = json.NewDecoder(r.Body).Decode(&v)

	if err != nil || len(v.Password) < 8 || len(v.Password) > 255 || !UsernameRegex.MatchString(v.Username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{\"error\":\"there was a problem with your request. Please try again with different values\"}")
		return
	}

	// search if user with that username already exists
	find := userCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: strings.ToLower(v.Username)}})
	if find.Err() == nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "{\"error\":\"user already exists with that username\"}")
		return
	}

	// create a new session for the new user
	sessionID := uuid.NewString()
	var session *http.Cookie
	if v.RememberMe {
		expire := time.Now().Add(30 * 24 * time.Hour)
		session = &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   0,
			Expires:  expire,
		}
	} else {
		session = &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   0,
		}
	}
	http.SetCookie(w, session)

	// hash and store the user's hashed password
	hasedPass, err := bcrypt.GenerateFromPassword([]byte(v.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"error\":\"internal server error, please try again later\"}")
	}

	// add the new user to the database
	var acc ExistingAccount
	acc.Username = strings.ToLower(v.Username)
	acc.Password = string(hasedPass)
	acc.Case = v.Username
	acc.Color = "white"
	acc.Points = 100
	acc.Resets = 0
	acc.Session = sessionID
	_, err = userCollection.InsertOne(r.Context(), acc)
	if err != nil {
		log.Println("* Error inserting new user")
	}

	// return the account information to the user
	var ret ReturnedAccount
	ret.Username = v.Username
	ret.Color = "white"
	ret.Points = 100
	ret.Resets = 0
	account, err := json.Marshal(ret)
	if err != nil {
		fmt.Println("* Error marshalling bson.D response")
	}
	fmt.Fprint(w, string(account))
}

func login(w http.ResponseWriter, r *http.Request) {
	// prepare DB collection
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	w.Header().Set("Content-Type", "application/json")
	userCollection := DB.Database("Users").Collection("Users")

	// decode POST into v struct
	var v Login
	json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"error\":\"bad request. Try again later\"}")
		return
	}

	// cmp struct will be compared with v to verify credentials
	var cmp Login

	found := userCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: strings.ToLower(v.Username)}})
	if found.Err() != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"error\":\"account with that username does not exist\"}")
		return
	}
	err = found.Decode(&cmp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(cmp.Password), []byte(v.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"error\":\"invalid password\"}")
		return
	}

	// prepare ReturnedAccount struct to be sent back to the client
	var account ReturnedAccount
	err = found.Decode(&account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set new session cookie for user, either persistent (remember me) or temporary
	sessionID := uuid.NewString()
	var session *http.Cookie
	if v.RememberMe {
		expire := time.Now().Add(30 * 24 * time.Hour)
		session = &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   0,
			Expires:  expire,
		}
	} else {
		session = &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   0,
		}
	}
	http.SetCookie(w, session)

	// update the new user session in the DB
	filter := bson.D{primitive.E{Key: "username", Value: strings.ToLower(account.Username)}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "session", Value: sessionID}}}}
	userCollection.UpdateOne(context.TODO(), filter, update, opts)

	acc, err := json.Marshal(account)
	if err != nil {
		fmt.Println("Error marshalling bson.D response")
	}
	fmt.Fprint(w, string(acc))
	//fmt.Println("logged in user successfully")
}

func loginBySession(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	userCollection := DB.Database("Users").Collection("Users")

	var id Session
	var account ReturnedAccount

	json.NewDecoder(r.Body).Decode(&id)

	filter := bson.D{primitive.E{Key: "session", Value: id.Session}}
	find := userCollection.FindOne(r.Context(), filter)
	if find.Err() != nil {
		log.Println(find.Err())
		fmt.Fprintf(w, "{\"error\":\"no user with given session id\"}")
		return
	}
	err = find.Decode(&account)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "{\"error\":\"cannot decode user bson from session\"}")
		return
	}

	json.NewEncoder(w).Encode(account)
}

func logout(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		DB = openDB()
	}
	userCollection := DB.Database("Users").Collection("Users")

	var v Credentials

	json.NewDecoder(r.Body).Decode(&v)

	filter := bson.D{primitive.E{Key: "username", Value: strings.ToLower(v.Username)}}
	opts := options.Update().SetUpsert(true)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "session", Value: ""}}}}
	userCollection.UpdateOne(context.TODO(), filter, update, opts)
}
