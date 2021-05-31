package controllers

import (
	"encoding/json"
	"net/http"

	"app/models"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"app/db"
	"app/service"
)

// GET all tasks

func HandleGetAllUser(c *gin.Context) {

	loadedUsers, err := service.GetAllUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Users": loadedUsers})
}

// POST a user

func HandleCreateUser(c *gin.Context) {
	var user models.User
	c.ShouldBindJSON(&user)

	id, err := service.CreateUser(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if err := db.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Created User": user})
}

// GET a User

func HandleGetSingleUser(c *gin.Context) {
	var user *models.User

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not a valid Hex ID"})
		return
	}

	if err := db.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User": user})
}

// Update the status of existing User

func HandleUpdateUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		return
	}

	user := models.User{}

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := service.UpdateUser(c, &id, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := db.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such user exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Updated User": user})
}

// Delete existing User

func HandleDeleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	var user models.User
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := db.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "No such user exists"})
		return
	}

	if err := service.DeleteUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted User": user})
}
