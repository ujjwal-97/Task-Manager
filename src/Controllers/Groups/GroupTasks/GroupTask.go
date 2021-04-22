package GroupTasks

import (
	"errors"
	"log"
	"time"

	"../../../Middleware"
	"../../../Models"
	"../../Connect"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CREATE
func CreateGroupTask(c *gin.Context, groupid primitive.ObjectID, task *Models.Task) (primitive.ObjectID, error) {

	task.Id = primitive.NewObjectID()
	task.PostedAt = primitive.NewDateTimeFromTime(time.Now())
	task.Author = Middleware.UserID
	task.GroupID = groupid.Hex()
	result, err := Connect.Collection.InsertOne(c, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err.Error())
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

//GET
func GetGroupTask(c *gin.Context) ([]*Models.Task, error) {

	var tasksAdmin []*Models.Task
	cursor1, err := Connect.Collection.Find(c, bson.M{"author": Middleware.UserID, "group": bson.M{
		"$exists": true,
	}})
	if err != nil {
		return nil, err
	}
	err = cursor1.All(c, &tasksAdmin)
	if err != nil {
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}
	log.Println(len(tasksAdmin))
	var tasksMember []*Models.Task
	cursor2, err := Connect.Collection.Find(c, bson.M{"assignedto": Middleware.UserID.Hex()})
	if err != nil {
		return nil, err
	}
	err = cursor2.All(c, &tasksMember)
	if err != nil {
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}

	log.Println(len(tasksMember))
	tasksAdmin = append(tasksAdmin, tasksMember...)
	return tasksAdmin, nil
}

//Alter the completed Flag
func AlterCompletion(c *gin.Context, id *primitive.ObjectID) error {
	var statusUpdate bool
	var task Models.Task

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id, "assignedto": Middleware.UserID.Hex()}).Decode(&task); err != nil {
		return err
	}
	//log.Println(task.Completed)
	statusUpdate = !task.Completed

	update := bson.M{"$set": bson.M{"completed": statusUpdate}}

	if _, err := Connect.Collection.UpdateByID(c, &id, update); err != nil {
		return err
	}
	return nil
}

//Update Status
func UpdateStatus(c *gin.Context, id *primitive.ObjectID, taskUpdate *Models.Task) error {

	var task Models.Task

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id, "assignedto": Middleware.UserID}).Decode(&task); err != nil {
		return err
	}
	var update primitive.M

	if len(taskUpdate.Status) != 0 {
		update = bson.M{"$set": bson.M{"members": append(task.Status, taskUpdate.Status...)}}
	}

	if _, err := Connect.Collection.UpdateByID(c, id, update); err != nil {
		return err
	}
	return nil
}

//DELETE
func DeleteGroupTask(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := Connect.Collection.DeleteOne(c, bson.M{"_id": &id, "author": Middleware.UserID}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such task exist")
		}
		return err
	}
	return nil
}
