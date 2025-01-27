package main

import (
	"github.com/LoveKhatri/basic-auth/controllers"
	"github.com/LoveKhatri/basic-auth/initializers"
	"github.com/LoveKhatri/basic-auth/middlewares"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	protected := r.Group("/")
	protected.Use(middlewares.RequireAuth)
	{
		protected.GET("/validate", controllers.Validate)
	}

	r.Run()
}
