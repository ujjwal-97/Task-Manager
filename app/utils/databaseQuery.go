package utils

import (
	"app/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type databaseQuery interface {
	Find(c *gin.Context) (*mongo.Cursor, error)
	FindOne(c *gin.Context) *mongo.SingleResult
	Insert(c *gin.Context) (*mongo.InsertOneResult, error)
	Update(c *gin.Context, update primitive.M) (*mongo.UpdateResult, error)
	Delete(c *gin.Context) (*mongo.DeleteResult, error)
}

type User models.User
type Task models.Task
