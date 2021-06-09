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
	"golang.org/x/crypto/bcrypt"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userId primitive.ObjectID
	taskId primitive.ObjectID
)

func TestEncryptPass(t *testing.T) {
	password := "asdffgjk"
	result, err := service.EncryptPass(password)
	assert.NotEmpty(t, result)
	assert.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(result), []byte(password))
	assert.NoError(t, err)
	password = "asdf"
	result, err = service.EncryptPass(password)
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(result), []byte(password))
	assert.NoError(t, err)

	password = ""
	result, err = service.EncryptPass(password)
	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword([]byte(result), []byte(password))
	assert.NoError(t, err)
}

var (
	VMname primitive.ObjectID
)

func TestGetAllUser(t *testing.T) {
	godotenv.Load("../../.env")
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	db.EstablishConnection()
	list, err := service.GetAllUser(con)
	assert.NoError(t, err)
	assert.NotEqual(t, list, nil)
	if len(list) > 0 {
		user := list[0]
		assert.NotEqual(t, user.Id, primitive.NilObjectID)
	}
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
	assert.NotEqual(t, userId, primitive.NilObjectID)

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
	updatedUser, err := service.GetSingleUser(con, &userId)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Name, user.Name)
}

//Get all Tasks
func TestGetAllTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	list, err := service.GetAllTask(con)
	assert.NoError(t, err)
	assert.NotEqual(t, list, nil)
	if len(list) > 0 {
		task := list[0]
		assert.NotEqual(t, task.Id, primitive.NilObjectID)
	}
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

	generatedID, err := service.CreateTask(&task, con)
	assert.Error(t, err)
	assert.Equal(t, generatedID, primitive.NilObjectID)

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
	updatedTask, err := service.GetSingleTask(con, &taskId)
	assert.NoError(t, err)
	assert.Equal(t, task.Status, updatedTask.Status)

	task.Status = "Wrong Status"

	err = service.UpdateTask(con, &taskId, &task)
	assert.Error(t, err)
}

//Delete Task
func TestDeleteTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := service.DeleteTask(con, &taskId)
	assert.NoError(t, err)

	err = service.DeleteTask(con, &taskId)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "task not found")
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
	assert.Equal(t, err.Error(), "no such user exist")
}
