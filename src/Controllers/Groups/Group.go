package Groups

import (
	"errors"
	"time"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../../Middleware"
	"../Connect"
)

func GetAllGroup(c *gin.Context) ([]*Models.Group, error) {

	var groups []*Models.Group
	cursor, err := Connect.Collection.Find(c, bson.M{"members": bson.M{"$in": []primitive.ObjectID{Middleware.UserID}}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(c, &groups)
	if err != nil {
		log.Printf("Failed marshalling %v", err.Error())
		return nil, err
	}
	return groups, nil
}

func CreateGroup(group *Models.Group, c *gin.Context) (primitive.ObjectID, error) {

	group.Id = primitive.NewObjectID()
	group.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	group.Admin = Middleware.UserID

	group.Members = append(group.Members, group.Admin)

	result, err := Connect.Collection.InsertOne(c, group)
	if err != nil {
		log.Printf("Could not create group: %v", err.Error())
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateGroup(c *gin.Context, id *primitive.ObjectID, groupUpdate *Models.Group) error {

	var group Models.Group

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id}).Decode(&group); err != nil {
		log.Println("Can't find group")
		return err
	}
	var update primitive.M

	if len(groupUpdate.Members) != 0 {
		update = bson.M{"$set": bson.M{"members": append(group.Members, groupUpdate.Members...)}}
	}

	if _, err := Connect.Collection.UpdateByID(c, &id, update); err != nil {
		return err
	}

	return nil
}

func DeleteGroup(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := Connect.Collection.DeleteOne(c, bson.M{"_id": &id, "admin": Middleware.UserID}); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such group exist")
		}
		return err
	}
	return nil
}

func RemoveMember(c *gin.Context, memberID *primitive.ObjectID, groupID *primitive.ObjectID) error {

	var newMemberList []primitive.ObjectID
	var group Models.Group

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &groupID}).Decode(&group); err != nil {
		return err
	}

	if *memberID != Middleware.UserID || Middleware.UserID != group.Admin {
		return errors.New("Unauthorized")
	}
	newMemberList = group.Members
	index := -1
	for i, n := range newMemberList {
		if n == *memberID {
			index = i
			break
		}
	}
	if index != -1 {
		newMemberList = append(newMemberList[:index], newMemberList[index+1:]...)
	}
	update := bson.M{"$set": bson.M{"members": newMemberList}}

	if _, err := Connect.Collection.UpdateByID(c, groupID, update); err != nil {
		return err
	}
	return nil
}
