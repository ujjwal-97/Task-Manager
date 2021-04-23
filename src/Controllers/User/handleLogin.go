package User

import (
	"net/http"

	"../../Controllers"
	"../../Service"
	"github.com/gin-gonic/gin"
)

var (
	jwtService      Service.JWTService          = Service.NewJWTService()
	loginController Controllers.LoginController = Controllers.NewLoginController(jwtService)
)

func HandleLogin(ctx *gin.Context) {
	token := loginController.Login(ctx)
	if token != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid Credentials"})
	}
}
