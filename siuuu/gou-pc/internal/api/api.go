package api

import (
	"gou-pc/internal/api/handler"
	"gou-pc/internal/api/middleware"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"
	"time"

	"github.com/gin-gonic/gin"
)

// Start khởi động API server với Gin, inject các service
func Start(port string, userService service.UserService, clientService service.ClientService, logService service.LogService, clientRepo repository.ClientRepository, jwtSecret string, jwtExpire time.Duration) {
	// Inject service vào handler
	handler.InjectUserService(userService)
	handler.InjectClientService(clientService)
	handler.InjectLogService(logService)
	handler.InjectOTPService(service.NewOTPService(clientRepo))
	handler.InjectJWTConfig(jwtSecret, int64(jwtExpire.Seconds()))
	// Inject config JWT cho middleware
	middleware.InitJWT(jwtSecret, jwtExpire)

	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Public route: chỉ login
	r.POST("/api/login", handler.LoginHandler)

	// Protected group: tất cả route còn lại đều cần JWT
	api := r.Group("/api", middleware.JWTAuthMiddlewareFunc())
	{
		// User routes
		api.POST("/users/create", middleware.JWTAuthMiddleware(handler.CreateUserHandler, true)) // admin only
		api.POST("/users/change-password", handler.ChangePasswordHandler)
		api.POST("/users/update", handler.UpdateUserHandler)
		api.GET("/users", middleware.JWTAuthMiddleware(handler.ListUsersHandler, true)) // admin only
		api.POST("/users/update-info", handler.UpdateUserInfoHandler)
		api.POST("/users/delete", middleware.JWTAuthMiddleware(handler.DeleteUserHandler, true)) // admin only

		// Client routes
		api.GET("/clients", handler.ListClientsHandler)
		api.GET("/clients/:agent_id", handler.GetClientByAgentIDHandler)
		api.GET("/clients/by-user/:user_id", handler.GetClientsByUserIDHandler)
		api.POST("/clients/delete", middleware.JWTAuthMiddleware(handler.DeleteClientHandler, true))   // admin only
		api.POST("/clients/assign-user", middleware.JWTAuthMiddleware(handler.HandleAssignUser, true)) // admin only
		// OTP routes
		api.GET("/clients/:agent_id/otp", handler.GetOTPByAgentIDHandler)
		api.GET("/clients/my-otp", handler.GetMyOTPHandler)

		// Log routes
		api.GET("/logs/archive", middleware.JWTAuthMiddleware(handler.GetArchiveLogHandler, true)) // admin only
		api.GET("/logs/my-device", handler.GetMyDeviceLogHandler)
	}

	logutil.Info("API server (Gin) starting on port %s...", port)
	r.Run(":" + port)
}
