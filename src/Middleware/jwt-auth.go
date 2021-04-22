package Middleware

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"../Service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var UserID primitive.ObjectID

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			log.Println("User not logged In")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "User not logged In"})
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := Service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			//log.Println("Claims[Uid]: ", claims["Uid"])
			useridString := claims["Uid"].(string)
			if uid, err := primitive.ObjectIDFromHex(useridString); err != nil {
				log.Println(err.Error())
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			} else {
				UserID = uid
			}
		} else {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		}
	}
}
