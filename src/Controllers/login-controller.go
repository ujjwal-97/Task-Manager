package Controllers

import (
	"log"

	"../Service"
	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	jWtService Service.JWTService
}

func NewLoginController(jWtService Service.JWTService) LoginController {
	return &loginController{
		jWtService: jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {

	user, err := Service.Login(ctx)
	if err == nil {
		return controller.jWtService.GenerateToken(user.Id.Hex(), true)
	}
	log.Println(err.Error())
	return ""
}
