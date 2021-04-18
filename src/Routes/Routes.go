package Routes

import (
	"../Controllers/Task"
	"../Controllers/User"
	"../Middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "TASK MANAGER APPLICATION",
		})
	})

	task := r.Group("/task", Middleware.AuthorizeJWT())
	{
		task.GET("", Task.HandleGetAllTask)
		task.POST("", Task.HandleCreateTask)
		task.GET("/:id", Task.HandleGetSingleTask)
		task.PUT("/:id", Task.HandleUpdateTask)
		task.DELETE("/:id", Task.HandleDeleteTask)
	}
	/*
		user := r.Group("/user")
		{
			user.GET("", User.HandleGetAllUser)
			user.POST("", User.HandleCreateUser)
			user.POST("/signup", User.Signup)
			user.GET("/:id", User.HandleGetSingleUser)
			user.PUT("/:id", User.HandleUpdateUser)
			user.DELETE("/:id", User.HandleDeleteUser)
		}
	*/
	user := r.Group("/user")
	{
		user.GET("", User.HandleGetAllUser)
		user.POST("/signup", User.Signup)
		user.POST("/login", User.HandleLogin)
		user.DELETE("/:id", User.HandleDeleteUser)
	}
	return r
}
