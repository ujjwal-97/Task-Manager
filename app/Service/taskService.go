package Service

import (
	"errors"
	"fmt"
	"time"

	"app/Models"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"app/CRON"
	"app/DB"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTask(c *gin.Context) ([]*Models.Task, error) {

	var task []*Models.Task
	cursor, err := DB.Collection.Find(c, bson.M{"email": bson.M{
		"$exists": false,
	}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &task)
	if err != nil {
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}
	return task, nil
}

func CreateTask(task *Models.Task, c *gin.Context) (primitive.ObjectID, error) {

	var user Models.User
	task.Id = primitive.NewObjectID()
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	statusList := map[string]int{"pending": 1, "inprogress": 2, "complete": 3}

	if _, exist := statusList[task.Status]; !exist {
		if task.Status != "" {
			return primitive.NilObjectID, errors.New("status only accepts 'pending, inprogress, complete' as options")
		}
		task.Status = "pending"
	}

	if task.TaskUser != nil {
		if err := DB.Collection.FindOne(c, bson.M{"_id": &task.TaskUser.Id}).Decode(&user); err == nil {
			task.TaskUser = &user
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
		}
	}

	// Create the cron
	if _, err := cron.ParseStandard(task.SnapshotSchedule); err != nil {
		task.SnapshotSchedule = "@every daily"
	}
	task.CronID, _ = CRON.C.AddFunc(task.SnapshotSchedule, func() {
		out, err := CRON.TakeSnapshot(user.Id.Hex(), user.Email)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(out))
	})

	result, err := DB.Collection.InsertOne(c, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err.Error())
		CRON.C.Remove(task.CronID)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)

	return oid, nil
}

func UpdateTask(c *gin.Context, id *primitive.ObjectID, taskUpdate *Models.Task) error {

	var task Models.Task

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		return err
	}
	var update bson.M
	if taskUpdate.Title != "" {
		update = bson.M{"$set": bson.M{"title": taskUpdate.Title}}
	}
	if taskUpdate.Description != "" {
		update = bson.M{"$set": bson.M{"description": taskUpdate.Description}}
	}

	statusList := map[string]int{"pending": 1, "inprogress": 2, "complete": 3}

	if _, exist := statusList[taskUpdate.Status]; exist {
		update = bson.M{"$set": bson.M{"status": taskUpdate.Status}}
		log.Println(taskUpdate.Status)
	} else {
		if taskUpdate.Status != "" {
			return errors.New("status only accepts 'pending, inprogress, complete' as options")
		}
	}

	if taskUpdate.TaskUser != nil {
		var user Models.User
		if err := DB.Collection.FindOne(c, bson.M{"_id": &taskUpdate.TaskUser.Id}).Decode(&user); err == nil {
			taskUpdate.TaskUser = &user
			RemoveTaskFromList(c, task)
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
			update = bson.M{"$set": bson.M{"user": taskUpdate.TaskUser}}
		}
	}

	if _, err := DB.Collection.UpdateByID(c, id, update); err != nil {
		return err
	}
	return nil
}

func DeleteTask(c *gin.Context, id *primitive.ObjectID) error {

	var task Models.Task
	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&task); err != nil {
		return errors.New("task not found")
	}

	RemoveTaskFromList(c, task)
	CRON.C.Remove(task.CronID)
	if result, err := DB.Collection.DeleteOne(c, bson.M{"_id": &id}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("task not found")
		}
		return err
	}
	return nil
}
