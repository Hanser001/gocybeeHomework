package api

import (
	"ginDemo/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.POST("/register", register) //注册
	r.POST("/login", login)       //登录

	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.GET("/get", getUsernameFromToken)
	}

	r.POST("/messager", leave)          //留言
	r.POST("/setquestion", addQuestion) //设置密保问题
	r.POST("/findpwd", findpassword)    //找回密码

	r.Run(":8088")
}
