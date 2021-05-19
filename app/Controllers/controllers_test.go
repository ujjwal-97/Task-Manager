package Controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/CRON"
	"app/DB"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllTask(t *testing.T) {
	w := httptest.NewRecorder()

	err := godotenv.Load("../.env")
	assert.NoError(t, err)

	con, _ := gin.CreateTestContext(w)
	c := *con
	DB.EstablishConnection()

	HandleGetAllTask(&c)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string][]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, ok := result["tasks"]
	assert.Equal(t, ok, true)

}

var (
	taskID string
	userID string
)

func TestCreateTask(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	CRON.C = cron.New()
	HandleCreateTask(&c)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Created Task"]["id"]
	assert.Equal(t, exists, true)

	taskID = value
}

func TestGetSingleTask1(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})

	HandleGetSingleTask(&c)
	assert.Equal(t, w.Code, 400)
	expected := `{"msg":"Can't find"}`
	assert.Equal(t, w.Body.String(), expected)

}

func TestGetSingleTask2(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	c.Params = append(c.Params, gin.Param{"id", taskID})

	HandleGetSingleTask(&c)
	assert.Equal(t, w.Code, 200)
	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Task"]["id"]
	assert.Equal(t, exists, true)
	assert.Equal(t, value, taskID)
}

func TestDeleteTask1(t *testing.T) {
	w := httptest.NewRecorder()

	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})

	HandleDeleteTask(&c)
	assert.Equal(t, w.Code, 400)

	expected := `{"msg":"mongo: no documents in result"}`
	assert.Equal(t, w.Body.String(), expected)

}

func TestDeleteTask2(t *testing.T) {
	w := httptest.NewRecorder()

	con, _ := gin.CreateTestContext(w)
	c := *con

	//id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", taskID})

	HandleDeleteTask(&c)
	assert.Equal(t, w.Code, 200)
	log.Println(w.Body.String())

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	value := result["Deleted Task"]["id"]
	assert.Equal(t, value, taskID)

}

func TestGetAllUser(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	HandleGetAllUser(&c)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string][]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, ok := result["Users"]
	assert.Equal(t, ok, true)

}

func TestCreateUser(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	CRON.C = cron.New()
	HandleCreateUser(&c)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	log.Println(w.Body.String())
	value, exists := result["Created User"]["id"]
	assert.Equal(t, exists, true)
	userID = value
}

func TestGetSingleUser1(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})

	HandleGetSingleUser(&c)
	assert.Equal(t, w.Code, 400)

	expected := `{"msg":"Can't find User"}`
	assert.Equal(t, w.Body.String(), expected)

}
func TestGetSingleUser2(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	c.Params = append(c.Params, gin.Param{"id", userID})

	HandleGetSingleUser(&c)
	assert.Equal(t, w.Code, 200)
	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["User"]["id"]
	assert.Equal(t, exists, true)
	assert.Equal(t, value, userID)
}

func TestDeleteUser1(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})

	HandleDeleteUser(&c)
	assert.Equal(t, w.Code, 400)
	expected := `{"msg":"No such user exists"}`
	assert.Equal(t, w.Body.String(), expected)
}
func TestDeleteUser2(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	c.Params = append(c.Params, gin.Param{"id", userID})

	HandleDeleteUser(&c)
	assert.Equal(t, w.Code, 200)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	value := result["Deleted User"]["id"]
	assert.Equal(t, value, userID)
}
