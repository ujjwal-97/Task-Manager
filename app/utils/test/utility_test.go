package utils

import (
	"net/http/httptest"
	"testing"

	"app/db"
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
	user := utils.User{}
	_, err := user.Find(con)
	assert.NoError(t, err)
}

func TestInsertUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = primitive.NewObjectID()
	userId = user.Id
	user.Email = "demo@email.com"
	user.Name = "demo"
	_, err := user.Insert(con)
	assert.NoError(t, err)
}

func TestFindOneUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = userId
	err := user.FindOne(con).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, user.Id, userId)
}

func TestUpdateUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	update = bson.M{"$set": bson.M{"password": "updatedPassword"}}
	update = bson.M{"$set": bson.M{"name": "updatedName"}}
	user := utils.User{}
	user.Id = userId
	result, err := user.Update(con, update)
	assert.NotEqual(t, 0, result.ModifiedCount)
	assert.NoError(t, err)
}

//Task

func TestInsertTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	task.Id = primitive.NewObjectID()
	taskId = task.Id
	task.Title = "demoTask"
	task.Status = "pending"
	_, err := task.Insert(con)
	assert.NoError(t, err)
}

func TestFindOneTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	task.Id = taskId
	err := task.FindOne(con).Decode(&task)
	assert.NoError(t, err)
	assert.Equal(t, task.Id, taskId)
}

func TestFindTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	_, err := task.Find(con)
	assert.NoError(t, err)
}

func TestUpdateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	update = bson.M{"$set": bson.M{"title": "updatedTitle"}}
	update = bson.M{"$set": bson.M{"status": "completed"}}
	task := utils.Task{}
	task.Id = taskId
	result, err := task.Update(con, update)
	assert.NotEqual(t, 0, result.ModifiedCount)
	assert.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	task.Id = taskId
	result, err := task.Delete(con)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
}

//User
func TestDeleteUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = userId
	result, err := user.Delete(con)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
}
