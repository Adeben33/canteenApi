package helpers

import (
	"canteenApi/database"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

type signinDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRETKEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, uid string) (signedToken string, refershToken string, err error) {
	claims := &signinDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &signinDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRETKEY))
	refershToken, err = jwt.NewWithClaims(jwt.SigningMethodES384, refreshClaims).SignedString([]byte(SECRETKEY))
	if err != nil {
		log.Panic(err)
	}
	return token, refershToken, nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshToken})
	updatedAT, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAT})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt)
	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func ValidateTokens(signedToken string) (claims *signinDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&signinDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRETKEY), nil
		},
	)

	claims, ok := token.Claims.(*signinDetails)
	if !ok {
		msg = fmt.Sprintf("the token is valid")
		msg = err.Error()
		return
	}

	// check if claims has expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("Token has expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
