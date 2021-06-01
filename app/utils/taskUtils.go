package utils

import (
	"app/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (task *Task) Find(c *gin.Context) (*mongo.Cursor, error) {
	return db.Collection.Find(c, bson.M{"email": bson.M{
		"$exists": true,
	}})
}

func (task *Task) FindOne(c *gin.Context) *mongo.SingleResult {
	return db.Collection.FindOne(c, bson.M{"_id": &task.Id})
}

func (task *Task) Insert(c *gin.Context) (*mongo.InsertOneResult, error) {
	return db.Collection.InsertOne(c, task)
}

func (task *Task) Update(c *gin.Context, update primitive.M) (*mongo.UpdateResult, error) {
	return db.Collection.UpdateByID(c, &task.Id, update)
}

func (task *Task) Delete(c *gin.Context) (*mongo.DeleteResult, error) {
	return db.Collection.DeleteOne(c, bson.M{"_id": &task.Id})
}
