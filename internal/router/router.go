package router

import (
	"chat-app-backend/internal/handlers"
	middleware "chat-app-backend/internal/middlewares"
	"chat-app-backend/internal/services"
	    "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3001"}, // Your frontend URL
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
	// Public routes
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)
	services.InitWebSocketServer()
	router := gin.Default()

	handlers.InitWebSocketRoutes(router)

	// Authenticated routes
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// User-related routes
		authenticated.GET("/users", handlers.GetAllUsersHandler)
		authenticated.PUT("/users/:id", handlers.UpdateUserHandler)
		authenticated.GET("/user", handlers.GetUserFromTokenHandler)

		// Message-related routes
		authenticated.POST("/messages", handlers.CreateMessageHandler)
		authenticated.GET("/messages", handlers.GetAllConversationsHandler)
		authenticated.GET("/messages/:other_user_id", handlers.GetConversationWithUserHandler)
		authenticated.DELETE("/messages/:id", handlers.DeleteMessageHandler)
		authenticated.PUT("/messages/:id", handlers.UpdateMessageHandler)
		authenticated.POST("/messages/:id/reply", handlers.ReplyToMessageHandler)
	}

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Admin user-related routes
		admin.GET("/users", handlers.GetAllUsersHandler)
		admin.DELETE("/users/:id", handlers.DeleteUserHandler)
	}

	return r
}
