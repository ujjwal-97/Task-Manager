package Controllers

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"../DB"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllTask(t *testing.T) {
	w := httptest.NewRecorder()
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(w)
	c := *con
	DB.EstablishConnection()
	HandleGetAllTask(&c)

	if w.Code != 200 {
		t.Error()
	}

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSingleTask1(t *testing.T) {
	w := httptest.NewRecorder()
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(w)
	c := *con
	DB.EstablishConnection()
	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleTask(&c)

	if w.Code != 400 {
		t.Error()
	}

}
func TestGetSingleTask2(t *testing.T) {
	w := httptest.NewRecorder()
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(w)
	c := *con
	DB.EstablishConnection()

	id := "608babeb8381719a6f548645"
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleTask(&c)
	log.Println(c.Params.ByName("id"))
	if w.Code != 200 {
		t.Error()
	}
}

func TestDeleteTask1(t *testing.T) {
	w := httptest.NewRecorder()
	if err := godotenv.Load("../.env"); err != nil {
		t.Errorf("Error loading .env file")
	}
	con, _ := gin.CreateTestContext(w)
	c := *con
	DB.EstablishConnection()
	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleTask(&c)

	if w.Code != 400 {
		t.Error()
	}
}
