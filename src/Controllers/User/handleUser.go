package User

import (
	"net/http"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../Connect"
)

// GET all tasks

func HandleGetAllUser(c *gin.Context) {

	loadedUsers, err := GetAllUser(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Users": loadedUsers})
}

// POST a user

func HandleCreateUser(c *gin.Context) {
	var user Models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	id, err := CreateUser(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET a User

func HandleGetSingleUser(c *gin.Context) {
	var user *Models.User

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not a valid Hex ID"})
		return
	}

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user ": user})
}

// Update the status of existing User

func HandleUpdateUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err)
		return
	}

	user := Models.User{}
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
	}

	if err := UpdateUser(c, &id, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
	}

	c.JSON(http.StatusOK, gin.H{"Updated the status of User with Id ": id})
}

// Delete existing User

func HandleDeleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err)
		return
	}

	if err := DeleteUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted User with Id ": id})
}
