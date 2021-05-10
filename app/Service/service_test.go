package Service

import (
	"net/http/httptest"
	"testing"

	"app/CRON"
	"app/DB"
	"app/Models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userId primitive.ObjectID
	taskId primitive.ObjectID
)

func TestEncryptPass(t *testing.T) {
	result, err := EncryptPass("asdffgjk")
	if len(result) == 0 || err != nil {
		t.Error()
	}
	result, err = EncryptPass("asdf")
	if len(result) == 0 || err != nil {
		t.Error()
	}
	result, err = EncryptPass("")
	if len(result) == 0 || err != nil {
		t.Error()
	}
}

func TestGetAll(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	DB.EstablishConnection()
	_, err := GetAllUser(con)
	if err != nil {
		t.Error()
	}
}

//Create a User
func TestCreateUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	userId, err = CreateUser(&user, con)
	if err != nil {
		t.Error()
	}
}

//Update User details
func TestUpdateUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Password = "password"

	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	random := primitive.NewObjectID()
	err := UpdateUser(con, &random, &user)
	if err == nil {
		t.Error()
	}
	err = UpdateUser(con, &userId, &user)
	if err != nil {
		t.Error()
	}
}

//Get all Tasks
func TestGetAllTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, err := GetAllTask(con)
	if err != nil {
		t.Error()
	}
}

//Create Task
func TestCreateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	CRON.C = cron.New()
	CRON.C.Start()

	task := Models.Task{}
	task.Title = "Task1"
	task.Description = "This is the test task"
	ID, err := CreateTask(&task, con)
	if err != nil {
		t.Error()
	}
	taskId = ID
	task.Status = "Wrong Status"
	_, err = CreateTask(&task, con)
	if err == nil {
		t.Error()
	}
}

//Update Task
func TestUpdateTask(t *testing.T) {
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	CRON.C = cron.New()
	CRON.C.Start()

	task := Models.Task{}
	task.Title = "Task2"
	task.Description = "New Descripton"
	task.Status = "inprogress"
	err := UpdateTask(con, &taskId, &task)
	if err != nil {
		t.Error()
	}
	task.Status = "Wrong Status"
	err = UpdateTask(con, &taskId, &task)
	if err == nil {
		t.Error()
	}
}

//Delete Task
func TestDeleteTask(t *testing.T) {

	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	//log.Println(taskId)
	err := DeleteUser(con, &taskId)
	if err != nil {
		t.Error()
	}
	err = DeleteUser(con, &taskId)
	if err == nil {
		t.Error()
	}
}

//Delete User
func TestDeleteUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"
	con, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := DeleteUser(con, &userId)
	if err != nil {
		t.Error()
	}
	err = DeleteUser(con, &userId)
	if err == nil {
		t.Error()
	}
}
