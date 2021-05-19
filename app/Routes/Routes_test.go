package Routes

import (
	"app/DB"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestWelcomePage(t *testing.T) {
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	expected := `{"message":"TASK MANAGER APPLICATION"}`
	assert.Equal(t, w.Body.String(), expected)

}

func TestUserRoutes(t *testing.T) {

	godotenv.Load("../.env")

	DB.EstablishConnection()
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string][]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, ok := result["Users"]
	assert.Equal(t, ok, true)
}

var (
	userID string
)

func TestCreateUser(t *testing.T) {

	godotenv.Load("../.env")
	DB.EstablishConnection()

	router := SetupRouter()
	var jsonStr = []byte(`{"name":"username"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Created User"]["id"]
	assert.Equal(t, exists, true)

	userID = value

}

func TestUserByID(t *testing.T) {

	godotenv.Load("../.env")

	DB.EstablishConnection()

	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/user/"+userID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["User"]["id"]
	assert.Equal(t, exists, true)
}

func TestUpdateUser(t *testing.T) {

	godotenv.Load("../.env")

	DB.EstablishConnection()

	router := SetupRouter()

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

}
func TestDeleteUser(t *testing.T) {

	godotenv.Load("../.env")

	DB.EstablishConnection()

	router := SetupRouter()

	req, _ := http.NewRequest("DELETE", "/user/"+userID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)

	var result map[string]map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	_, exists := result["Deleted User"]["id"]
	assert.Equal(t, exists, true)
}
