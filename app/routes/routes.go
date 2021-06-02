package routes

import (
	"app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	app := r.Group("/")
	{
		app.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "TASK MANAGER APPLICATION",
			})
		})
	}

	user := r.Group("/user")
	{
		user.POST("", controllers.HandleCreateUser)
		user.GET("", controllers.HandleGetAllUser)
		user.PUT("/:id", controllers.HandleUpdateUser)
		user.POST("/snap/:id", controllers.HandleSnapshot)
		user.GET("/:id", controllers.HandleGetSingleUser)
		user.DELETE("/:id", controllers.HandleDeleteUser)
	}

	task := r.Group("/task")
	{
		task.GET("", controllers.HandleGetAllTask)
		task.POST("", controllers.HandleCreateTask)
		task.GET("/:id", controllers.HandleGetSingleTask)
		task.PUT("/:id", controllers.HandleUpdateTask)
		task.DELETE("/:id", controllers.HandleDeleteTask)
	}

	return r
}
