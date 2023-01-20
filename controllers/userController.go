package controllers

import (
	"canteenApi/database"
	"canteenApi/helpers"
	"canteenApi/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var userCollections *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		userId := c.Param("userId")
		err := userCollections.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		//	records per page
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		//page
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		//startIndex
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{
					{"$slice", []interface{}{"$data", startIndex, recordPerPage}},
				}},
			}},
		}
		//Aggregate
		result, err := userCollections.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user item"})
		}
		var allUser []bson.M
		if err = result.All(ctx, allUser); err != nil {
			log.Panic(err)
		}
		c.JSON(http.StatusOK, allUser)

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
		//err := userCollections.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		//if err == nil {
		//	msg := fmt.Sprintf("User all exist")
		//	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		//}

		count, err := userCollections.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occuring while counting"})
			return
		}
		passwordHarsh := HashPassword(*user.Password)
		user.Password = &passwordHarsh
		err = userCollections.FindOne(ctx, bson.M{"email": user.Phone}).Decode(&existingUser)
		//if err == nil {
		//	msg := fmt.Sprintf("User all exist")
		//	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		//}
		count, err = userCollections.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occuring while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number alrady exists"})
			return
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		// Generate Token and Refresh token
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserId)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		result, InsertErr := userCollections.InsertOne(ctx, user)
		if InsertErr != nil {
			msg := fmt.Sprintf("user Item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		//	bind the input
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//	find the user in the db

		err := userCollections.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		//	verifying the password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		//	return the user
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	stringHarsh := string(byte)
	return stringHarsh
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Loginor password is incorrect")
		log.Panic(err)
	}
	return check, msg
}
