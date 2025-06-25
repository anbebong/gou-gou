package handler

import (
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var logService service.LogService

func InjectLogService(s service.LogService) { logService = s }

func GetArchiveLogHandler(c *gin.Context) {
	logs, err := logService.GetAllLogs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, logs)
}

func GetMyDeviceLogHandler(c *gin.Context) {
	agentID := c.Query("agent_id")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	logs, err := logService.GetLogsByAgentID(agentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, logs)
}
