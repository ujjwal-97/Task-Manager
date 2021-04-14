package Task

import (
	"errors"
	"time"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../Connect"
)

func GetAllTask(c *gin.Context) ([]*Models.Task, error) {

	var tasks []*Models.Task
	cursor, err := Connect.Collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &tasks)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return tasks, nil
}

func CreateTask(task *Models.Task, c *gin.Context) (primitive.ObjectID, error) {

	task.Id = primitive.NewObjectID()
	task.PostedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := Connect.Collection.InsertOne(c, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateTask(c *gin.Context, id *primitive.ObjectID) error {
	var statusUpdate bool
	var task Models.Task

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		return err
	}
	statusUpdate = !task.Completed

	update := bson.M{"$set": bson.M{"completed": statusUpdate}}

	if _, err := Connect.Collection.UpdateByID(c, id, update); err != nil {
		return err
	}
	return nil
}

func DeleteTask(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := Connect.Collection.DeleteOne(c, bson.M{"_id": &id}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such task exist")
		}
		return err
	}
	return nil
}
