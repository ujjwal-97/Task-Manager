package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty" `
	Name      string               `json:"name" bson:"name"`
	TaskList  []primitive.ObjectID `json:"tasklist" bson:"tasklist"`
	Email     string               `json:"email" bson:"email" `
	Password  string               `json:"-" bson:"password"`
	CreatedAt primitive.DateTime   `json:"createdAt" bson:"createdAt"`
}
