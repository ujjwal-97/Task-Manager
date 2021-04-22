package Routes

import (
	"../Controllers/Groups"
	"../Controllers/Groups/GroupTasks"
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

	user := r.Group("/user")
	{
		user.GET("", User.HandleGetAllUser)
		user.POST("/signup", User.Signup)
		user.POST("/login", User.HandleLogin)
		user.DELETE("/:id", User.HandleDeleteUser)
	}

	task := r.Group("/task", Middleware.AuthorizeJWT())
	{
		task.GET("", Task.HandleGetAllTask)
		task.POST("", Task.HandleCreateTask)
		task.GET("/:id", Task.HandleGetSingleTask)
		task.PUT("/:id", Task.HandleUpdateTask)
		task.DELETE("/:id", Task.HandleDeleteTask)
	}

	groups := r.Group("/group", Middleware.AuthorizeJWT())
	{
		groups.GET("", Groups.HandleGetAllGroup)
		groups.POST("", Groups.HandleCreateGroup)
		groups.GET("/:id", Groups.HandleGetSingleGroup)
		groups.DELETE("/:id", Groups.HandleDeleteGroup)
		groups.PUT("/addMember/:id", Groups.HandleAddMember)
		groups.PUT("/removeMember/:id", Groups.HandleRemoveMember)
		groups.PUT("/changeAdmin/:id", Groups.HandleChangeAdmin)
	}

	grouptask := r.Group("/grouptask", Middleware.AuthorizeJWT())
	{
		grouptask.GET("", GroupTasks.HandleGetGroupTask)
		grouptask.POST("/:id", GroupTasks.HandleCreateGroupTask)
		grouptask.DELETE("/:id", GroupTasks.HandleDeleteGroupTask)
		grouptask.PUT("/completion/:id", GroupTasks.HandleAlterCompletion)
		grouptask.PUT("/status/:id", GroupTasks.HandleUpdateStatus)

	}
	return r
}
