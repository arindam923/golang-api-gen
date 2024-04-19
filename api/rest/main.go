package main

import (
	"github.com/arindam923/api-generator/api/rest/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Post routes
	posts := r.Group("/posts")
	{
		posts.GET("/", handlers.ListPosts)
		posts.GET("/:id", handlers.GetPost)
		posts.POST("/", handlers.CreatePost)
		posts.PUT("/:id", handlers.UpdatePost)
		posts.DELETE("/:id", handlers.DeletePost)
	}

	r.Run()
}
