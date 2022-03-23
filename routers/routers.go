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
	r.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", controller.IndexHandler)

	// v1
	v1Group := r.Group("")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", controller.CreateTodo)
		// 查看所有的待办事项
		v1Group.GET("/todo", controller.GetTodoList)
		// 修改某一个待办事项
		v1Group.PUT("/todo/:id", controller.UpdateATodo)
		// 删除某一个待办事项
		v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	}
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

		groupRoutes.GET("/mygroup", controller.GetMyGroups)
		groupRoutes.GET("/mygroupnum", controller.GetMyGroupsNum)

		groupRoutes.POST("/join", controller.JoinGroup)

	}
	//homework
	//HomeWorkRoutes := r.Group("")
	//{
	//
	//	HomeWorkRoutes.POST("/work", controller.CreateWork)
	//	HomeWorkRoutes.PUT("/work", controller.Me)
	//	HomeWorkRoutes.GET("/login", controller.Login)
	//
	//}
	//file
	//FileRoutes := r.Group("")
	//{
	//
	//	FileRoutes.POST("/file", controller.Upload)
	//	FileRoutes.GET("/get", controller.Me)
	//	FileRoutes.POST("/login", controller.Login)
	//
	//}
	return r
}
