package Routes

import (
	"../Controllers"

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
		user.POST("", Controllers.HandleCreateUser)
		user.GET("", Controllers.HandleGetAllUser)
		user.GET("/:id", Controllers.HandleGetSingleUser)
		user.DELETE("/:id", Controllers.HandleDeleteUser)
	}

	task := r.Group("/task")
	{
		task.GET("", Controllers.HandleGetAllTask)
		task.POST("", Controllers.HandleCreateTask)
		task.GET("/:id", Controllers.HandleGetSingleTask)
		task.PUT("/:id", Controllers.HandleUpdateTask)
		task.DELETE("/:id", Controllers.HandleDeleteTask)
	}

	return r
}
