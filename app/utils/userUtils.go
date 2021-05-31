package utils

import (
	"app/db"
	"app/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(c *gin.Context) (*mongo.Cursor, error) {
	return db.Collection.Find(c, bson.M{"email": bson.M{
		"$exists": true,
	}})
}

func FindOneUser(c *gin.Context, id *primitive.ObjectID) *mongo.SingleResult {
	return db.Collection.FindOne(c, bson.M{"_id": id})
}

func InsertUser(c *gin.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return db.Collection.InsertOne(c, user)
}

func UpdateUser(c *gin.Context, id *primitive.ObjectID, update primitive.M) (*mongo.UpdateResult, error) {
	return db.Collection.UpdateByID(c, id, update)
}

func DeleteUser(c *gin.Context, id *primitive.ObjectID) (*mongo.DeleteResult, error) {
	return db.Collection.DeleteOne(c, bson.M{"_id": &id})
}
