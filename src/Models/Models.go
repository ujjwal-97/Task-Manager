package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id" `
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description" `
	Completed   bool               `json:"completed" bson:"completed" `
	Deadline    primitive.DateTime `json:"deadline,omitempty" bson:"deadline,omitempty"`
	PostedAt    primitive.DateTime `json:"posttime" bson:"posttime"`
	Author      primitive.ObjectID `json:"author" bson:"author"`
}

type User struct {
	Id        primitive.ObjectID   `json:"id" bson:"_id" `
	Name      string               `json:"name" bson:"name"`
	TaskList  []primitive.ObjectID `json:"tasklist" bson:"tasklist"`
	Email     string               `json:"email" bson:"email" `
	Password  string               `json:"password" bson:"password"`
	CreatedAt primitive.DateTime   `json:"creationtime" bson:"creationtime"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
