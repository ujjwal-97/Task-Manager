package service

import (
	"app/cronjob"
	"app/models"
	"app/utils"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(c *gin.Context) ([]*models.User, error) {

	var users []*models.User
	cursor, err := utils.FindUser(c)
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

func CreateUser(user *models.User, c *gin.Context) (primitive.ObjectID, error) {

	user.Id = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.Password, _ = EncryptPass(user.Password)

	result, err := utils.InsertUser(c, user)
	if err != nil {
		log.Printf("Could not create User: %v", err.Error())
		return primitive.NilObjectID, err
	}
	//Creation of vm with name as user ID
	go CreateVM(user.Id)
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func UpdateUser(c *gin.Context, id *primitive.ObjectID, userUpdate *models.User) error {

	var user models.User

	if err := utils.FindOneUser(c, id).Decode(&user); err != nil {
		return err
	}
	var update primitive.M

	if userUpdate.Password != "" {
		userUpdate.Password, _ = EncryptPass(user.Password)
		update = bson.M{"$set": bson.M{"password": userUpdate.Password}}
	}
	if userUpdate.Name != "" {
		update = bson.M{"$set": bson.M{"name": userUpdate.Name}}
	}
	if len(userUpdate.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"tasklist": append(user.TaskList, userUpdate.TaskList...)}}
	}

	if _, err := utils.UpdateUser(c, id, update); err != nil {
		return err
	}

	return nil
}

func DeleteUser(c *gin.Context, id *primitive.ObjectID) error {

	if result, err := utils.DeleteUser(c, id); err != nil || result.DeletedCount == 0 {
		if result.DeletedCount == 0 {
			return errors.New("no such user exist")
		}
		return err
	}
	go RemoveVM(*id)
	return nil
}

func EncryptPass(pass string) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(bytes), err
}

func AddTaskTOList(c *gin.Context, user *models.User) {

	var update bson.M
	if len(user.TaskList) != 0 {
		update = bson.M{"$set": bson.M{"tasklist": user.TaskList}}
	}
	utils.UpdateUser(c, &user.Id, update)
}

func RemoveTaskFromList(c *gin.Context, task models.Task) {
	var user *models.User
	if task.TaskUser == nil {
		return
	}
	if err := utils.FindOneUser(c, &task.TaskUser.Id).Decode(&user); err != nil {
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
	utils.UpdateUser(c, &user.Id, update)
}

func CreateVM(userid primitive.ObjectID) (string, error) {

	filename := os.Getenv("VMImageFileName")
	cmd := "VBoxManage " + " import " + filename + " --vsys " + " 0 " + " --vmname " + userid.Hex()

	out, err := cronjob.ExecCommandOnHost(cmd)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	log.Println(out + "\n --! VM created successfully !--")

	return string(out), nil
}

func RemoveVM(userid primitive.ObjectID) (string, error) {

	cmd := "VBoxManage " + " unregistervm " + userid.Hex() + " --delete"
	out, err := cronjob.ExecCommandOnHost(cmd)
	if err != nil {
		return "", nil
	}
	log.Println(out + "\n !-- VM deleted successfully --!")
	return string(out), nil
}
