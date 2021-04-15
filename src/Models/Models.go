package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" `
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" `
	Completed   bool               `json:"completed" bson:"completed" `
	Deadline    primitive.DateTime `json:"deadline,omitempty" bson:"deadline,omitempty"`
	PostedAt    primitive.DateTime `json:"posttime" bson:"posttime"`
}

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id" `
	Name      string             `json:"name" bson:"name"`
	TaskList  []Task             `json:"tasklist" bson:"tasklist"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	CreatedAt primitive.DateTime `json:"creationtime" bson:"creationtime"`
}
