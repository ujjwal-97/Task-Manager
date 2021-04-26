package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" `
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description" `
	Deadline    primitive.DateTime `json:"deadline,omitempty" bson:"deadline,omitempty"`
	CreatedAt   primitive.DateTime `json:"createdAt" bson:"createdAt"`
	TaskUser    *User              `json:"user,omitempty" bson:"user,omitempty"`
	Status      string             `json:"status" bson:"status"`
}
