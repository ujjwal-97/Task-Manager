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
	//log.Println(token)
	if token != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, nil)
	}
}
