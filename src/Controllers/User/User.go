package User

import (
	"errors"
	"time"

	"../../Models"
	"golang.org/x/crypto/bcrypt"

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
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}
	return users, nil
}

func CreateUser(user *Models.User, c *gin.Context) (primitive.ObjectID, error) {

	user.Id = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.Password, _ = EncryptPass(user.Password)

	result, err := Connect.Collection.InsertOne(c, user)
	if err != nil {
		log.Printf("Could not create User: %v", err.Error())
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
		userUpdate.Password, _ = EncryptPass(user.Password)
		update = bson.M{"$set": bson.M{"password": userUpdate.Password}}

	}
	if len(userUpdate.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"tasklist": append(user.TaskList, userUpdate.TaskList...)}}
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

func CheckUserExists(c *gin.Context, email string) (Models.User, bool, string) {
	// Filter
	condition := bson.M{"email": email}
	var res Models.User
	// Searching the email
	err := Connect.Collection.FindOne(c, condition).Decode(&res)
	ID := res.Id.Hex()
	if err != nil {
		return res, false, ID
	}
	return res, true, ID
}

func EncryptPass(pass string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(bytes), err
}
