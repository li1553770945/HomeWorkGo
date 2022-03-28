package routers

import (
	"HomeWorkGo/controller"
	"HomeWorkGo/setting"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionId", store))
	// 告诉gin框架模板文件引用的静态文件去哪里找
	//r.Static("/static", "static")
	//// 告诉gin框架去哪里找模板文件
	//r.LoadHTMLGlob("templates/*")

	//user
	userRoutes := r.Group("/api")
	{

		userRoutes.POST("/user", controller.Register)
		userRoutes.GET("/user", controller.Me)
		userRoutes.GET("/logout", controller.Logout)
		userRoutes.POST("/login", controller.Login)

	}
	//Group
	groupRoutes := r.Group("/api")
	{

		groupRoutes.POST("/group", controller.CreateGroup)
		groupRoutes.GET("/group", controller.GetGroup)
		groupRoutes.DELETE("/group", controller.DeleteGroup)

		groupRoutes.GET("/group-created", controller.GetGroupsCreated)
		groupRoutes.GET("/group-created-num", controller.GetGroupsCreatedNum)

		groupRoutes.POST("/group-member", controller.JoinGroup)
		groupRoutes.DELETE("/group-member", controller.QuitGroup)

		groupRoutes.GET("/group-joined", controller.GetGroupsJoined)
		groupRoutes.GET("/group-joined-num", controller.GetMyGroupsJoinedNum)

	}

	HomeWorkRoutes := r.Group("/api")
	{

		HomeWorkRoutes.POST("/homework", controller.CreateHomework)
		HomeWorkRoutes.GET("/homework", controller.GetHomework)
		HomeWorkRoutes.DELETE("/homework", controller.DeleteHomework)
		HomeWorkRoutes.GET("/homework-created", controller.GetHomeworkCreated)
		HomeWorkRoutes.GET("/homework-created-num", controller.GetHomeworkCreatedNum)

	}

	SubmissionRoutes := r.Group("/api")
	{

		SubmissionRoutes.GET("/submissions", controller.GetSubmissionsByHomeworkId)
		SubmissionRoutes.GET("/submission-file", controller.GetSubmissionFileById)

		SubmissionRoutes.GET("/homework-joined-num", controller.GetHomeworkJoinedNum)
		SubmissionRoutes.GET("/homework-joined", controller.GetHomeworkJoined)

		SubmissionRoutes.POST("/submit", controller.Submit)
		SubmissionRoutes.GET("/export", controller.Export)
		SubmissionRoutes.GET("/download", controller.Download)
	}
	return r
}
