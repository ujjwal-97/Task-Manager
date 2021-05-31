package utils

import (
	"net/http/httptest"
	"testing"

	"app/db"
	"app/models"
	"app/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userId primitive.ObjectID
	taskId primitive.ObjectID
)

func TestFindUser(t *testing.T) {
	godotenv.Load("../../.env")

	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	db.EstablishConnection()
	_, err := utils.FindUser(con)
	assert.NoError(t, err)
}

func TestInsertUser(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var user models.User
	user.Id = primitive.NewObjectID()
	userId = user.Id
	user.Email = "demo@email.com"
	user.Name = "demo"
	_, err := utils.InsertUser(con, &user)
	assert.NoError(t, err)
}

func TestFindOneUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var user models.User
	err := utils.FindOneUser(con, &userId).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, user.Id, userId)
}

func TestUpdateUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	update = bson.M{"$set": bson.M{"name": "updatedName"}}
	update = bson.M{"$set": bson.M{"password": "updatedPassword"}}
	result, err := utils.UpdateUser(con, &taskId, update)
	assert.NotEqual(t, 0, result.ModifiedCount)
	assert.NoError(t, err)
}

//Task

func TestInsertTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var task models.Task
	task.Id = primitive.NewObjectID()
	taskId = task.Id
	task.Title = "demoTask"
	task.Status = "pending"
	_, err := utils.InsertTask(con, &task)
	assert.NoError(t, err)
}

func TestFindOneTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var task models.Task
	err := utils.FindOneTask(con, &taskId).Decode(&task)
	assert.NoError(t, err)
	assert.Equal(t, task.Id, taskId)
}

func TestFindTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, err := utils.FindTask(con)
	assert.NoError(t, err)
}

func TestUpdateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	update = bson.M{"$set": bson.M{"title": "updatedTitle"}}
	update = bson.M{"$set": bson.M{"status": "completed"}}
	result, err := utils.UpdateTask(con, &taskId, update)
	assert.NotEqual(t, 0, result.ModifiedCount)
	assert.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	result, err := utils.DeleteTask(con, taskId)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
}

//User
func TestDeleteUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	result, err := utils.DeleteUser(con, &userId)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
}
