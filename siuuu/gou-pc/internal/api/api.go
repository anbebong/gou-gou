package api

import (
	"gou-pc/internal/api/handler"
	"gou-pc/internal/api/middleware"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"

	"github.com/gin-gonic/gin"
)

// Start khởi động API server với Gin, inject các service
func Start(port string, userService service.UserService, clientService service.ClientService, logService service.LogService) {
	// Inject service vào handler
	handler.InjectUserService(userService)
	handler.InjectClientService(clientService)
	handler.InjectLogService(logService)

	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Public route: chỉ login
	r.POST("/api/login", handler.LoginHandler)

	// Protected group: tất cả route còn lại đều cần JWT
	api := r.Group("/api", middleware.JWTAuthMiddlewareFunc())
	{
		// User routes
		api.POST("/users/create", handler.CreateUserHandler)
		api.POST("/users/change-password", handler.ChangePasswordHandler)
		api.POST("/users/update", handler.UpdateUserHandler)
		api.GET("/users", handler.ListUsersHandler)
		api.POST("/users/update-info", handler.UpdateUserInfoHandler)

		// Client routes
		api.GET("/clients", handler.ListClientsHandler)
		api.GET("/clients/:agent_id", handler.GetClientByAgentIDHandler)
		api.GET("/clients/by-user/:user_id", handler.GetClientsByUserIDHandler)
		api.POST("/clients/delete", handler.DeleteClientHandler)
		api.POST("/clients/assign-user", handler.HandleAssignUser)

		// Message/OTP
		api.POST("/message/send", handler.SendMessageToClientHandler)
		api.GET("/otp", handler.HandleGetOTP)

		// Log routes
		api.GET("/logs/archive", handler.GetArchiveLogHandler)
		api.GET("/logs/my-device", handler.GetMyDeviceLogHandler)
	}

	logutil.Info("API server (Gin) starting on port %s...", port)
	r.Run(":" + port)
}
