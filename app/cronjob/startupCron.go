package cronjob

import (
	"app/db"
	"app/models"
	"app/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/robfig/cron.v2"
)

func StartUpPreviousCron(c *cron.Cron) error {
	var tasklist []*models.Task
	utilsTask := utils.Task{}
	godotenv.Load("../../.env")
	db.EstablishConnection()
	cursor, err := utilsTask.Find(&gin.Context{})
	if err != nil {
		return err
	}
	err = cursor.All(&gin.Context{}, &tasklist)
	if err != nil {
		return err
	}

	if len(tasklist) > 0 {
		for _, task := range tasklist {
			if cronID, err := CreateSnapshotCron(task.SnapshotSchedule, task.TaskUser); err == nil {
				task.CronID = cronID
				var update bson.M = bson.M{"$set": bson.M{"cronid": task.CronID}}
				if _, err := utilsTask.Update(&gin.Context{}, update); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
