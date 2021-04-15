package User

import (
	"errors"
	"time"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../Connect"
)

func GetAllUser(c *gin.Context) ([]*Models.User, error) {

	var users []*Models.User
	cursor, err := Connect.Collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &users)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return users, nil
}

func CreateUser(user *Models.User, c *gin.Context) (primitive.ObjectID, error) {

	user.Id = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := Connect.Collection.InsertOne(c, user)
	if err != nil {
		log.Printf("Could not create User: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateUser(c *gin.Context, id *primitive.ObjectID, userUpdate *Models.User) error {

	var user Models.User

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		return err
	}
	var update primitive.M
	if userUpdate.Password != "" {
		update = bson.M{"$set": bson.M{"password": userUpdate.Password}}

	}
	if len(userUpdate.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"TaskList": append(user.TaskList, userUpdate.TaskList...)}}
	}

	if _, err := Connect.Collection.UpdateByID(c, id, update); err != nil {
		return err
	}

	return nil
}

func DeleteUser(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := Connect.Collection.DeleteOne(c, bson.M{"_id": &id}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such user exist")
		}
		return err
	}
	return nil
}
