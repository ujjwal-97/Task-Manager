package Task

import (
	"net/http"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../../Middleware"
	"../Connect"
)

// GET all tasks

func HandleGetAllTask(c *gin.Context) {

	var loadedTasks, err = GetAllTask(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": loadedTasks})
}

// POST a task

func HandleCreateTask(c *gin.Context) {
	var task Models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	id, err := CreateTask(&task, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET a Task

func HandleGetSingleTask(c *gin.Context) {
	var task *Models.Task

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not a valid Hex ID"})
		return
	}

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id, "author": Middleware.UserID}).Decode(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task ": task})
}

// Update the status of existing Task

func HandleUpdateTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := UpdateTask(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated the status of Task with Id ": id})
}

// Delete existing Task

func HandleDeleteTask(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		return
	}

	if err := DeleteTask(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted Task with Id ": id})
}
