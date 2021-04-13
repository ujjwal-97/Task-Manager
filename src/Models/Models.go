package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id        primitive.ObjectID `json:"id" bson:"_id" `
	Title     string             `json:"title" bson:"title" binding:"required"`
	Completed bool               `json:"completed" bson:"completed" `
}
