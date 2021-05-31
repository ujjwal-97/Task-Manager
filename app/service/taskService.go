package service

import (
	"errors"
	"fmt"
	"time"

	"app/models"
	"app/utils"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	"app/cronjob"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTask(c *gin.Context) ([]*models.Task, error) {

	var task []*models.Task
	cursor, err := utils.FindTask(c)
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

func CreateTask(task *models.Task, c *gin.Context) (primitive.ObjectID, error) {

	var user models.User
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
		if err := utils.FindOneUser(c, &task.TaskUser.Id).Decode(&user); err == nil {
			task.TaskUser = &user
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
		}
	}

	// Create the cron
	if _, err := cron.ParseStandard(task.SnapshotSchedule); err != nil {
		task.SnapshotSchedule = "@every daily"
	}
	task.CronID, _ = cronjob.C.AddFunc(task.SnapshotSchedule, func() {
		out, err := cronjob.TakeSnapshot(user.Id.Hex(), user.Email)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(out))
	})

	result, err := utils.InsertTask(c, task)
	if err != nil {
		log.Printf("Could not create Task: %v", err.Error())
		cronjob.C.Remove(task.CronID)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)

	return oid, nil
}

func UpdateTask(c *gin.Context, id *primitive.ObjectID, taskUpdate *models.Task) error {

	var task models.Task

	if err := utils.FindOneTask(c, id).Decode(&task); err != nil {
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
		var user models.User
		if err := utils.FindOneUser(c, &taskUpdate.TaskUser.Id).Decode(&user); err == nil {
			taskUpdate.TaskUser = &user
			RemoveTaskFromList(c, task)
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
			update = bson.M{"$set": bson.M{"user": taskUpdate.TaskUser}}
		}
	}

	if _, err := utils.UpdateTask(c, id, update); err != nil {
		return err
	}
	return nil
}

func DeleteTask(c *gin.Context, id *primitive.ObjectID) error {

	var task models.Task
	if err := utils.FindOneTask(c, id).Decode(&task); err != nil {
		return errors.New("task not found")
	}

	RemoveTaskFromList(c, task)
	cronjob.C.Remove(task.CronID)
	if result, err := utils.DeleteTask(c, *id); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("task not found")
		}
		return err
	}
	return nil
}
