package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" `
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" `
	Completed   bool               `json:"completed" bson:"completed" `
	Deadline    primitive.DateTime `json:"deadline" bson:"deadline"`
	PostedAt    primitive.DateTime `json:"posttime" bson:"posttime"`
}
