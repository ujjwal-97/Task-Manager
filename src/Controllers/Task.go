package Controllers

import (
	"net/http"

	"../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"./Connect"
)

// GET all tasks

func HandleGetAllTask(c *gin.Context) {
	var loadedTasks, err = GetAllTask(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": loadedTasks})
}

func GetAllTask(c *gin.Context) ([]*Models.Task, error) {

	var tasks []*Models.Task
	cursor, err := Connect.Collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)
	err = cursor.All(c, &tasks)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return tasks, nil
}

// POST a task

func HandleCreateTask(c *gin.Context) {
	var task Models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	id, err := CreateTask(&task, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func CreateTask(task *Models.Task, c *gin.Context) (primitive.ObjectID, error) {

	task.Id = primitive.NewObjectID()

	result, err := Connect.Collection.InsertOne(c, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}
