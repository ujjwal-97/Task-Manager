package service

import (
	"net/http/httptest"
	"testing"

	"app/cronjob"
	"app/db"
	"app/models"

	"app/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userId primitive.ObjectID
	taskId primitive.ObjectID
)

func TestEncryptPass(t *testing.T) {
	result, err := service.EncryptPass("asdffgjk")
	assert.NotEmpty(t, result)
	assert.NoError(t, err)

	result, err = service.EncryptPass("asdf")
	assert.NotEmpty(t, result)
	assert.NoError(t, err)

	result, err = service.EncryptPass("")
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
}

var (
	VMname primitive.ObjectID
)

func TestGetAll(t *testing.T) {
	godotenv.Load("../../.env")

	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	db.EstablishConnection()
	_, err := service.GetAllUser(con)
	assert.NoError(t, err)
}

//Create a User
func TestCreateUser(t *testing.T) {
	user := models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"

	godotenv.Load("../../.env")

	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	var err error
	userId, err = service.CreateUser(&user, con)
	assert.NoError(t, err)
}

//Update User details
func TestUpdateUser(t *testing.T) {
	user := models.User{}
	user.Name = "User1"
	user.Password = "password"

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	random := primitive.NewObjectID()

	err := service.UpdateUser(con, &random, &user)
	assert.Error(t, err)

	err = service.UpdateUser(con, &userId, &user)
	assert.NoError(t, err)
}

//Get all Tasks
func TestGetAllTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, err := service.GetAllTask(con)
	assert.NoError(t, err)
}

//Create Task
func TestCreateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	cronjob.C = cron.New()
	cronjob.C.Start()

	task := models.Task{}
	task.Title = "Task1"
	task.Description = "This is the test task"

	ID, err := service.CreateTask(&task, con)
	assert.NoError(t, err)

	taskId = ID
	task.Status = "Wrong Status"

	_, err = service.CreateTask(&task, con)
	assert.Error(t, err)
}

//Update Task
func TestUpdateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	cronjob.C = cron.New()
	cronjob.C.Start()

	task := models.Task{}
	task.Title = "Task2"
	task.Description = "New Descripton"
	task.Status = "inprogress"

	err := service.UpdateTask(con, &taskId, &task)
	assert.NoError(t, err)

	task.Status = "Wrong Status"

	err = service.UpdateTask(con, &taskId, &task)
	assert.Error(t, err)
}

//Delete Task
func TestDeleteTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	//log.Println(taskId)
	err := service.DeleteUser(con, &taskId)
	assert.NoError(t, err)

	err = service.DeleteUser(con, &taskId)
	assert.Error(t, err)
}

//Delete User
func TestDeleteUser(t *testing.T) {
	user := models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := service.DeleteUser(con, &userId)
	assert.NoError(t, err)

	err = service.DeleteUser(con, &userId)
	assert.Error(t, err)
}
