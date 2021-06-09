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
	cursor, err := user.Find(con)
	assert.NoError(t, err)
	assert.NotEqual(t, cursor, nil)
	var list []utils.User
	err = cursor.All(con, &list)
	assert.NoError(t, err)
	if len(list) > 0 {
		user = list[0]
		assert.NotEqual(t, user.Id, primitive.NilObjectID)
	}
}

func TestInsertUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = primitive.NewObjectID()
	userId = user.Id
	user.Email = "demo@email.com"
	user.Name = "demo"
	result, err := user.Insert(con)
	assert.NoError(t, err)
	assert.Equal(t, result.InsertedID, userId)
}

func TestFindOneUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = userId
	res := user.FindOne(con)
	assert.NotNil(t, res)
	err := res.Decode(&user)
	assert.NoError(t, err)
	assert.NotEqual(t, user.Id, primitive.NilObjectID)
	assert.Equal(t, user.Id, userId)
}

func TestUpdateUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	updatedName := "updatedName"
	update = bson.M{"$set": bson.M{"password": "updatedPassword"}}
	update = bson.M{"$set": bson.M{"name": updatedName}}
	user := utils.User{}
	user.Id = userId
	result, err := user.Update(con, update)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, result.ModifiedCount)
	user.FindOne(con).Decode(&user)
	assert.Equal(t, user.Name, updatedName)
}

//Task

func TestInsertTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	task.Id = primitive.NewObjectID()
	taskId = task.Id
	task.Title = "demoTask"
	task.Status = "pending"
	result, err := task.Insert(con)
	assert.NoError(t, err)
	assert.Equal(t, result.InsertedID, taskId)
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
	cursor, err := task.Find(con)
	assert.NoError(t, err)
	assert.NotEqual(t, cursor, nil)
	var list []utils.Task
	err = cursor.All(con, &list)
	assert.NoError(t, err)
	if len(list) > 0 {
		task = list[0]
		assert.NotEqual(t, task.Id, primitive.NilObjectID)
	}
}

func TestUpdateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	var update primitive.M
	updatedStatus := "completed"
	update = bson.M{"$set": bson.M{"title": "updatedTitle"}}
	update = bson.M{"$set": bson.M{"status": updatedStatus}}
	task := utils.Task{}
	task.Id = taskId
	result, err := task.Update(con, update)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, result.ModifiedCount)
	task.FindOne(con).Decode(&task)
	assert.NoError(t, err)
	assert.Equal(t, task.Status, updatedStatus)

}

func TestDeleteTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	task := utils.Task{}
	task.Id = taskId
	result, err := task.Delete(con)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
	res := task.FindOne(con).Decode(&task)
	assert.Error(t, res)
}

//User
func TestDeleteUser(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	user := utils.User{}
	user.Id = userId
	result, err := user.Delete(con)
	assert.NoError(t, err)
	assert.NotEqual(t, result.DeletedCount, 0)
	res := user.FindOne(con).Decode(&user)
	assert.Error(t, res)
}
