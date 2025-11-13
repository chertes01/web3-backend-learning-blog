package routes

import (
	"blog/controllers"
	"blog/middle"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/")
	// 需要认证的路由
	auth.Use(middle.JWTAuthMiddleware())
	{
		auth.POST("/posts", controllers.CreatePost)
		auth.GET("/posts", controllers.GetAllPosts)
		auth.GET("/posts/:post_id", controllers.GetPostByID)
		auth.GET("/users/:user_id/posts", controllers.GetPostsByUser)
		auth.PUT("/posts/:post_id", controllers.UpdatePost)
		auth.DELETE("/posts/:post_id", controllers.DeletePost)

		auth.POST("/comments", controllers.CreateComment)
		auth.GET("/posts/:post_id/comments", controllers.GetCommentsByPost)
		auth.DELETE("/comments/:comment_id", controllers.DeleteComment)

	}
	return r

}
