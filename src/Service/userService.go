package Service

import (
	"errors"
	"time"

	"../Models"
	"golang.org/x/crypto/bcrypt"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../DB"
)

func GetAllUser(c *gin.Context) ([]*Models.User, error) {

	var users []*Models.User
	cursor, err := DB.Collection.Find(c, bson.M{"email": bson.M{
		"$exists": true,
	}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &users)
	if err != nil {
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}
	return users, nil
}

func CreateUser(user *Models.User, c *gin.Context) (primitive.ObjectID, error) {

	user.Id = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.Password, _ = EncryptPass(user.Password)

	result, err := DB.Collection.InsertOne(c, user)
	if err != nil {
		log.Printf("Could not create User: %v", err.Error())
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateUser(c *gin.Context, id *primitive.ObjectID, userUpdate *Models.User) error {

	var user Models.User

	if err := DB.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&user); err != nil {
		return err
	}
	var update primitive.M

	if userUpdate.Password != "" {
		userUpdate.Password, _ = EncryptPass(user.Password)
		update = bson.M{"$set": bson.M{"password": userUpdate.Password}}
	}
	if userUpdate.Name != "" {
		userUpdate.Name, _ = EncryptPass(user.Name)
		update = bson.M{"$set": bson.M{"name": userUpdate.Name}}

	}
	if len(userUpdate.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"tasklist": append(user.TaskList, userUpdate.TaskList...)}}
	}

	if _, err := DB.Collection.UpdateByID(c, id, update); err != nil {
		return err
	}

	return nil
}

func DeleteUser(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := DB.Collection.DeleteOne(c, bson.M{"_id": &id}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such user exist")
		}
		return err
	}
	return nil
}

func EncryptPass(pass string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(bytes), err
}

func AddTaskTOList(c *gin.Context, user *Models.User) {
	var update bson.M
	if len(user.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"tasklist": user.TaskList}}
	}
	DB.Collection.UpdateByID(c, user.Id, update)
}

func RemoveTaskFromList(c *gin.Context, task Models.Task) {
	var user *Models.User
	if err := DB.Collection.FindOne(c, bson.M{"_id": task.TaskUser}).Decode(&user); err != nil {
		return
	}
	var update bson.M
	updatedList := user.TaskList

	index := -1
	for i, element := range updatedList {
		if element == task.Id {
			index = i
		}
	}
	if index != -1 {
		updatedList = append(updatedList[:index], updatedList[index+1:]...)
	}

	update = bson.M{"$set": bson.M{"tasklist": updatedList}}
	DB.Collection.UpdateByID(c, user.Id, update)
}
