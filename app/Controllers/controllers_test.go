package Controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"app/DB"

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
	con, _ := gin.CreateTestContext(w)
	c := *con
	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleTask(&c)

	if w.Code != 400 {
		t.Error()
	}

}

/*
func TestCreateTask(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	var u = []byte(`{"Body":{"name":"username"}}`)
	buf := bytes.NewBuffer(u)
	c.Request.Write(buf)
	c.Request.Header.Set("Content-Type", "application/json")

	HandleCreateTask(&c)
	log.Println(c.Request.Body)
	if w.Code != 200 {
		t.Error()
	}
}
*/
func TestGetSingleTask2(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := "609a0d3b4a6388bd9ca071fa"
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleTask(&c)
	//log.Println(c.Request.Body)
	if w.Code != 200 {
		t.Error()
	}
}

func TestDeleteTask1(t *testing.T) {
	w := httptest.NewRecorder()

	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleDeleteTask(&c)

	if w.Code != 400 {
		t.Error()
	}
}

func TestGetAllUser(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con
	HandleGetAllUser(&c)

	if w.Code != 200 {
		t.Error()
	}

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSingleUser1(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleUser(&c)
	if w.Code != 400 {
		t.Error()
	}
}
func TestGetSingleUser2(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := "608be7afe0071f4741c00c6b"
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleGetSingleUser(&c)
	if w.Code != 200 {
		t.Error()
	}
}

func TestDeleteUser(t *testing.T) {
	w := httptest.NewRecorder()
	con, _ := gin.CreateTestContext(w)
	c := *con

	id := primitive.NewObjectID().Hex()
	c.Params = append(c.Params, gin.Param{"id", id})
	HandleDeleteUser(&c)

	if w.Code != 400 {
		t.Error()
	}
}
