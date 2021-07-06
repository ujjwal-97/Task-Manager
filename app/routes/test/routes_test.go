package routes

import (
	"app/cronjob"
	"app/db"
	"app/routes"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestWelcomePage(t *testing.T) {
	router := routes.SetupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, w.Code, http.StatusBadRequest)
	expected := `{"message":"TASK MANAGER APPLICATION"}`
	assert.Equal(t, w.Body.String(), expected)
}

func TestUserRoutes(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()
	router := routes.SetupRouter()

	req, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, w.Code, http.StatusBadRequest)
	var result map[string][]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, ok := result["Users"]
	assert.Equal(t, ok, true)
}

var (
	userID string
	taskID string
)

func TestCreateUser(t *testing.T) {

	godotenv.Load("../../.env")
	db.EstablishConnection()

	router := routes.SetupRouter()
	var jsonStr = []byte(`{"name":"username"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, w.Code, http.StatusBadRequest)
	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Created User"]["id"]
	assert.Equal(t, exists, true)

	userID = value

}

func TestUserByID(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	req, _ := http.NewRequest("GET", "/user/"+userID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, w.Code, http.StatusBadRequest)

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["User"]["id"]
	assert.Equal(t, exists, true)

	invalidID := "123"
	req, _ = http.NewRequest("GET", "/user/"+invalidID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, 400)
	expected := `{"msg":"Not a valid Hex ID"}`
	assert.Equal(t, w.Body.String(), expected)

}

func TestUpdateUser(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	var jsonStr = []byte(`{"name":"Adam"}`)

	req, _ := http.NewRequest("PUT", "/user/"+userID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["Updated User"]["id"]
	assert.Equal(t, exists, true)

	invalidID := "123"
	req, _ = http.NewRequest("PUT", "/user/"+invalidID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, 400)
	expected := `{"msg":"encoding/hex: odd length hex string"}`
	assert.Equal(t, w.Body.String(), expected)

	invalidID = primitive.NewObjectIDFromTimestamp(time.Now()).Hex()
	req, _ = http.NewRequest("PUT", "/user/"+invalidID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, 400)
	expected = `{"msg":"Can't find User"}`
	assert.Equal(t, w.Body.String(), expected)

}
func TestCreateTask(t *testing.T) {
	godotenv.Load("../../.env")
	db.EstablishConnection()

	cronjob.C = cron.New()
	cronjob.C.Start()
	router := routes.SetupRouter()
	var jsonStr = []byte(`{"title":"TaskName}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Created Task"]["id"]
	assert.Equal(t, exists, true)
	taskID = value

}

func TestTaskRoutes(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()
	router := routes.SetupRouter()

	req, err := http.NewRequest("GET", "/task", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string][]string
	json.Unmarshal(w.Body.Bytes(), &result)
	log.Println(result)
	_, ok := result["tasks"]
	assert.Equal(t, ok, true)
}
func TestUpdateTask(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	var jsonStr = []byte(`{"description":"This is task"}`)

	req, _ := http.NewRequest("PUT", "/task/"+taskID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["Updated Task"]["id"]
	assert.Equal(t, exists, true)

	invalid := "123"
	req, _ = http.NewRequest("PUT", "/task/"+invalid, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, 400)
	expected := `{"msg":"encoding/hex: odd length hex string"}`
	assert.Equal(t, w.Body.String(), expected)

	invalid = primitive.NewObjectIDFromTimestamp(time.Now()).Hex()
	req, _ = http.NewRequest("PUT", "/task/"+invalid, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Code, 400)
	expected = `{"msg":"no such task found"}`
	assert.Equal(t, w.Body.String(), expected)
}
func TestTaskByID(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	req, _ := http.NewRequest("GET", "/task/"+taskID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Task"]["id"]
	assert.Equal(t, exists, true)
	assert.Equal(t, value, taskID)
}

func TestDeleteTask(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	req, _ := http.NewRequest("DELETE", "/task/"+taskID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["Deleted Task"]["id"]
	assert.Equal(t, exists, true)
}

func TestDeleteUser(t *testing.T) {

	godotenv.Load("../../.env")

	db.EstablishConnection()

	router := routes.SetupRouter()

	req, _ := http.NewRequest("DELETE", "/user/"+userID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["Deleted User"]["id"]
	assert.Equal(t, exists, true)
}
