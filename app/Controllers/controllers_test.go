package Controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/DB"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
/*
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
}*/

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

/*
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
*/
func TestDeleteUser(t *testing.T) {
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
