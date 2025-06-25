package handler

import (
	"gou-pc/internal/api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListClientsHandler(c *gin.Context) {
	clients, err := clientService.GetAllClients()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, clients)
}

func GetClientByAgentIDHandler(c *gin.Context) {
	agentID := c.Param("agent_id")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	client, err := clientService.GetClientByAgentID(agentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if client == nil {
		response.Error(c, http.StatusNotFound, "client not found")
		return
	}
	response.Success(c, client)
}

func GetClientsByUserIDHandler(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "user_id required")
		return
	}
	clients, err := clientService.GetClientsByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, clients)
}

func DeleteClientHandler(c *gin.Context) {}
func HandleAssignUser(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id"`
		UserID   string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ClientID == "" || req.UserID == "" {
		response.Error(c, http.StatusBadRequest, "client_id and user_id required")
		return
	}
	if err := clientService.AssignUserToClient(req.ClientID, req.UserID); err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user assigned to client successfully"})
}
