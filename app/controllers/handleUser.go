package controllers

import (
	"encoding/json"
	"net/http"

	"app/models"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	var user *models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		user = &models.User{}
	}

	id, err := service.CreateUser(user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	if user, err = service.GetSingleUser(c, &id); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Not a valid Hex ID"})
		return
	}

	if user, err = service.GetSingleUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"User": user})
}

// Update the status of existing User

func HandleUpdateUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	user := models.User{}

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := service.UpdateUser(c, &id, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}
	var updatedUser *models.User
	if updatedUser, err = service.GetSingleUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Updated User": updatedUser})
}

// Delete existing User

func HandleDeleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	var user *models.User
	if err != nil {
		log.Println(err.Error())
		return
	}
	if user, err = service.GetSingleUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find User"})
		return
	}

	if err := service.DeleteUser(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted User": user})
}
func HandleSnapshot(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	uuid := ""
	if uuid, err = service.Snapshot(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Snapshot taken successfully", "uuid": uuid})
}
