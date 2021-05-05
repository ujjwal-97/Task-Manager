package Controllers

import (
	"encoding/json"
	"net/http"

	"app/Models"

	"log"

	"github.com/gin-gonic/gin"

	"app/DB"
	"app/Service"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET all tasks

func HandleGetAllTask(c *gin.Context) {
	var loadedTasks, err = Service.GetAllTask(c)
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

	id, err := Service.CreateTask(&task, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Created Task": task})
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

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Task": task})
}

// Update the status of existing Task

func HandleUpdateTask(c *gin.Context) {

	var task Models.Task
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&task); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := Service.UpdateTask(c, &id, &task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated Task": task})
}

// Delete existing Task

func HandleDeleteTask(c *gin.Context) {

	var task Models.Task
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := Service.DeleteTask(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted Task": task})
}
