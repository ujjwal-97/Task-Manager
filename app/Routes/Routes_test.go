package Routes

import (
	"app/DB"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func TestWelcomePage(t *testing.T) {
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Error()
	}
}

func TestUserRoutes(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Error()
	}

}

var (
	userID string
)

func TestCreateUser(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()
	var jsonStr = []byte(`{"name":"username"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Error()
	}

	var result map[string]map[string]string

	json.Unmarshal(w.Body.Bytes(), &result)
	value, exists := result["Created User"]["id"]

	if exists {
		userID = value
	}
}

func TestUserByID(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/user/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Error()
	}
}

func TestUpdateUser(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()

	var jsonStr = []byte(`{"name":"Adam"}`)

	req, _ := http.NewRequest("PUT", "/user/"+userID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Error()
	}
}
func TestDeleteUser(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	DB.EstablishConnection()

	router := SetupRouter()

	req, _ := http.NewRequest("DELETE", "/user/"+userID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Error()
	}
}
