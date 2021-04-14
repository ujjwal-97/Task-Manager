package Routes

import (
	"../Controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "TASK MANAGER APPLICATION",
		})
	})
	user := r.Group("/user")
	{
		user.GET("/task", Controllers.HandleGetAllTask)
		user.POST("/task", Controllers.HandleCreateTask)
		user.GET("/task/:id", Controllers.HandleGetSingleTask)
		user.PUT("/task/:id", Controllers.HandleUpdateTask)
		user.DELETE("/task/:id", Controllers.HandleDeleteTask)
	}
	return r
}
