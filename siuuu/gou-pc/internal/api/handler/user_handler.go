package handler

import (
	"gou-pc/internal/api/model"
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	userService   service.UserService
	clientService service.ClientService
)

func InjectUserService(s service.UserService)     { userService = s }
func InjectClientService(s service.ClientService) { clientService = s }

func LoginHandler(c *gin.Context) {
	// TODO: implement login logic
	response.Error(c, http.StatusNotImplemented, "not implemented")
}

func CreateUserHandler(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.PasswordHash == "" {
		response.Error(c, http.StatusBadRequest, "username and password required")
		return
	}
	if err := userService.CreateUser(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user created successfully"})
}

func ChangePasswordHandler(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == "" || req.NewPassword == "" {
		response.Error(c, http.StatusBadRequest, "user_id and new_password required")
		return
	}
	// TODO: hash password trước khi lưu
	if err := userService.UpdatePassword(req.UserID, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "password changed successfully"})
}

func UpdateUserHandler(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ID == "" {
		response.Error(c, http.StatusBadRequest, "user_id required")
		return
	}
	if err := userService.UpdateUser(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user updated successfully"})
}

func ListUsersHandler(c *gin.Context) {
	users, err := userService.ListUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, users)
}

func UpdateUserInfoHandler(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ID == "" {
		response.Error(c, http.StatusBadRequest, "user_id required")
		return
	}
	if err := userService.UpdateUser(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user info updated successfully"})
}
