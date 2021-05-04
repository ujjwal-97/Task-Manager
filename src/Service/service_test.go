package Service

import (
	"net/http/httptest"
	"testing"

	"../DB"
	"../Models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userId primitive.ObjectID
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
	c := *con
	DB.EstablishConnection()
	_, err := GetAllUser(&c)
	if err != nil {
		t.Error()
	}
}

func TestCreateUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	c := *con
	DB.EstablishConnection()
	userId, err = CreateUser(&user, &c)
	if err != nil {
		t.Error()
	}
}

func TestUpdateUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Password = "password"
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	c := *con
	DB.EstablishConnection()
	random := primitive.NewObjectID()
	err := UpdateUser(&c, &random, &user)
	if err == nil {
		t.Error()
	}
	err = UpdateUser(&c, &userId, &user)
	if err != nil {
		t.Error()
	}
}

func TestDeleteUser(t *testing.T) {
	user := Models.User{}
	user.Name = "User1"
	user.Email = "user@email.com"
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	c := *con
	DB.EstablishConnection()
	err := DeleteUser(&c, &userId)
	if err != nil {
		t.Error()
	}
	err = DeleteUser(&c, &userId)
	if err == nil {
		t.Error()
	}
}

func TestGet(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(httptest.NewRecorder())
	c := *con
	DB.EstablishConnection()
	_, err := GetAllTask(&c)
	if err != nil {
		t.Error()
	}
}
