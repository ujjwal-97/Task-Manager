package User

import (
	"encoding/json"
	"log"
	"net/http"

	"../../Models"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	user := Models.User{}

	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if len(user.Email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Email is required"})
		return
	}
	if len(user.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Password length must be atleast 8 characters"})
		return
	}
	//checked
	if _, found, _ := CheckUserExists(c, user.Email); found {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User already exists"})
		return
	}

	id, err := CreateUser(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"User created with id": id})
}
