package controllers

import (
	"canteenApi/database"
	"canteenApi/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userCollections *mongo.Collection = database.OpenCollection(database.Client, "user")

func Getusers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Bind the incoming user login to struct
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var existingUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding user"})
			return
		}
		validateErr := validate.Struct(&user)

		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating user input"})
			return
		}
		err := userCollections.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			msg := fmt.Sprintf("User all exist")
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}
		passwordHarsh := HashPassword(*user.Password)
		user.Password = & passwordHarsh
		err = userCollections.FindOne(ctx, bson.M{"email": user.Phone}).Decode(&existingUser)
		if err == nil {
			msg := fmt.Sprintf("User all exist")
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}
		user.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	harsh,_:= bcrypt.GenerateFromPassword([]byte(password), 10))
	stringHarsh := string(harsh)
	return stringHarsh
}

func VerifyPassword(userPassword, providedPassword string) bool {
	return false
}
