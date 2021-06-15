package service

import (
	"errors"
	"time"

	"app/models"
	"app/utils"

	"log"

	"github.com/gin-gonic/gin"

	"app/cronjob"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTask(c *gin.Context) ([]*models.Task, error) {

	var task []*models.Task
	utilsTask := utils.Task{}
	cursor, err := utilsTask.Find(c)
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

func GetSingleTask(c *gin.Context, id *primitive.ObjectID) (*models.Task, error) {

	var task *models.Task
	utilTask := utils.Task{}
	utilTask.Id = *id
	cursor := utilTask.FindOne(c)
	err := cursor.Decode(&task)
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
		utilsUser := utils.User{}
		utilsUser.Id = task.TaskUser.Id
		if err := utilsUser.FindOne(c).Decode(&user); err == nil {
			task.TaskUser = &user
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
		}
	}

	// Create the cron

	cronID, err := cronjob.CreateSnapshotCron(task.SnapshotSchedule, &user)
	if err != nil {
		task.CronID = cronID
	}

	var utilsTask utils.Task = utils.Task(*task)
	result, err := utilsTask.Insert(c)
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
	utilsTask := utils.Task{}
	utilsTask.Id = *id
	if err := utilsTask.FindOne(c).Decode(&task); err != nil {
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
		utilsUser := utils.User{}
		utilsUser.Id = taskUpdate.TaskUser.Id
		if err := utilsUser.FindOne(c).Decode(&user); err == nil {
			taskUpdate.TaskUser = &user
			RemoveTaskFromList(c, task)
			user.TaskList = append(user.TaskList, task.Id)
			AddTaskTOList(c, &user)
			update = bson.M{"$set": bson.M{"user": taskUpdate.TaskUser}}
		}
	}

	utilsTask.Id = *id
	if _, err := utilsTask.Update(c, update); err != nil {
		return err
	}
	return nil
}

func DeleteTask(c *gin.Context, id *primitive.ObjectID) error {

	var task models.Task
	utilsTask := utils.Task{}
	utilsTask.Id = *id
	if err := utilsTask.FindOne(c).Decode(&task); err != nil {
		return errors.New("task not found")
	}

	RemoveTaskFromList(c, task)
	cronjob.C.Remove(task.CronID)

	if result, err := utilsTask.Delete(c); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("task not found")
		}
		return err
	}
	return nil
}

func UpdateTaskCronID(c *gin.Context, task *models.Task) {
	if task.CronID != 0 {
		update := bson.M{"$set": bson.M{"cronid": task.CronID}}
		utilsTask := utils.Task{}
		utilsTask.Id = *&task.Id
		utilsTask.Update(c, update)
	}
}
