package router

import (
	"chat/api"
	"chat/docs"
	"chat/middleware"
	"chat/ws"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.Title = "Chat API"
	docs.SwaggerInfo.Description = "This is a chat server API documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/index", api.GetIndex)
	r.GET("/userList", api.GetUserList)
	r.GET("/user/createUser", api.CreateUser)
	r.GET("/ws", func(c *gin.Context) { ws.ServeWS(c.Writer, c.Request) })

	// JWT middleware and auth routes
	authMiddleware := middleware.JWTMiddleware()
	r.POST("/user/login", authMiddleware.LoginHandler)

	// protected routes
	auth := r.Group("/")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.DELETE("/user/:id", api.DeleteUser)
	auth.GET("/messages", api.GetMessages)
	auth.PUT("/user/:id", api.UpdateUser)
	auth.PATCH("/user/:id", api.PartialUpdateUser)

	return r
}
