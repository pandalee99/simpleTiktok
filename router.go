package main

import (
	"simpleTiktok/controller"
	"simpleTiktok/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	g := r.Group("/douyin")

	// Base API
	g.GET("/feed/", jwt.AuthTourist, controller.Feed)
	g.POST("/user/register/", controller.UserRegister)
	g.POST("/user/login/", controller.UserLogin)
	g.GET("/user/", jwt.AuthUser, controller.UserInfo)

	p := g.Group("/publish")

	p.POST("/action/", jwt.AuthFormUser, controller.PublishAction)
	p.GET("/list/", jwt.AuthUser, controller.PublishListVideos)

	// 互动 API
	g.POST("/favorite/action/", jwt.AuthUser, controller.FavoriteAction)
	g.GET("/favorite/list/", jwt.AuthUser, controller.FavoriteList)
	g.POST("/comment/action/", jwt.AuthUser, controller.CommentAction)
	g.GET("/comment/list/", jwt.AuthUser, controller.CommentList)

	// 社交 API
	re := g.Group("/relation")
	re.POST("/action/", jwt.AuthUser, controller.RelationAction)
	re.GET("/follow/list/", jwt.AuthUser, controller.RelationFollowList)
	re.GET("/follower/list/", jwt.AuthUser, controller.RelationFollowerList)

	// 静态文件
	g.GET("/static/:authorId/:fileName", jwt.AuthTourist, controller.GetFileContent)
}
