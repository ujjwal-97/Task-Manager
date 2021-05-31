package utils

import (
	"app/db"
	"app/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindTask(c *gin.Context) (*mongo.Cursor, error) {
	return db.Collection.Find(c, bson.M{"email": bson.M{
		"$exists": false,
	}})
}

func FindOneTask(c *gin.Context, id *primitive.ObjectID) *mongo.SingleResult {
	return db.Collection.FindOne(c, bson.M{"_id": id})
}

func InsertTask(c *gin.Context, task *models.Task) (*mongo.InsertOneResult, error) {
	return db.Collection.InsertOne(c, task)
}

func UpdateTask(c *gin.Context, id *primitive.ObjectID, update primitive.M) (*mongo.UpdateResult, error) {
	return db.Collection.UpdateByID(c, id, update)
}

func DeleteTask(c *gin.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return db.Collection.DeleteOne(c, bson.M{"_id": &id})
}
