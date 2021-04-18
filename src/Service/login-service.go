package Service

import (
	"encoding/json"
	"errors"
	"log"

	"../Controllers/Connect"
	"../Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func Login(c *gin.Context) (Models.User, error) {

	var user Models.User

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		log.Print(err.Error())
		return user, err
	}

	if len(user.Email) == 0 {
		return user, errors.New("Email is required")
	}

	//checked

	obj, err := TrytoLogin(c, user.Email, user.Password)
	if err != nil {
		return user, err
	}

	return obj, nil
}

func TrytoLogin(c *gin.Context, email string, password string) (Models.User, error) {
	user, found, _ := CheckUserExists(c, email)
	if found == false {
		return user, errors.New("User doesn't exist")
	}
	// Validate password
	passByte := []byte(password)
	passBd := []byte(user.Password)
	// first goes pass encrypted and afer the normal
	err := bcrypt.CompareHashAndPassword(passBd, passByte)
	if err != nil {
		return user, errors.New("Invalid Password")
	}
	return user, nil
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
