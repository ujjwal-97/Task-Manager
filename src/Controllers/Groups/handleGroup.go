package Groups

import (
	"encoding/json"
	"net/http"

	"../../Models"

	"log"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"../../Middleware"
	"../Connect"
)

// GET all groups that a user is member of

func HandleGetAllGroup(c *gin.Context) {

	var loadedGroups, err = GetAllGroup(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"groups": loadedGroups})
}

// POST a task

func HandleCreateGroup(c *gin.Context) {
	var group Models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	id, err := CreateGroup(&group, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET a Task

func HandleGetSingleGroup(c *gin.Context) {
	var group *Models.Group

	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not a valid Hex ID"})
		return
	}

	if err := Connect.Collection.FindOne(c, bson.M{"_id": &id, "members": bson.M{"$in": []primitive.ObjectID{Middleware.UserID}}}).Decode(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Can't find"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group ": group})
}

// Update the status of group
//id: groupId, body contains memebers
func HandleAddMember(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		return
	}

	group := Models.Group{}

	if err := json.NewDecoder(c.Request.Body).Decode(&group); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	//log.Println(group.Name)
	if err := AddMember(c, &id, &group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		log.Print(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"Updated the status of group with Id ": id})
}

// Delete existing Group

func HandleDeleteGroup(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err := DeleteGroup(c, &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Deleted Group with Id ": id})
}

func HandleRemoveMember(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	var member Models.User
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := RemoveMember(c, &member.Id, &groupID); err != nil {
		if err.Error() == "Unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Removed member from Group with Id ": member.Id})
}

func HandleChangeAdmin(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	var member Models.User
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := ChangeAdmin(c, &groupID, &member.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Membership Upgraded to Admin": member.Id})
}
