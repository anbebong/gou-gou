package handler

import (
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var logService service.LogService

func InjectLogService(s service.LogService) { logService = s }

func GetArchiveLogHandler(c *gin.Context) {
	logutil.Debug("GetArchiveLogHandler called")
	logs, err := logService.GetAllLogs()
	if err != nil {
		logutil.Debug("GetArchiveLogHandler error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	logutil.Debug("GetArchiveLogHandler success, %d logs", len(logs))
	response.Success(c, logs)
}

func GetMyDeviceLogHandler(c *gin.Context) {
	logutil.Debug("GetMyDeviceLogHandler called")
	userID, ok := c.Get("username")
	if !ok {
		logutil.Debug("GetMyDeviceLogHandler missing username in context")
		response.Error(c, http.StatusUnauthorized, "username not found in context")
		return
	}
	agentIDs, err := clientService.GetAgentIDsByUserID(userID.(string))
	if err != nil || len(agentIDs) == 0 {
		logutil.Debug("GetMyDeviceLogHandler: user %v has no agentIDs", userID)
		response.Error(c, http.StatusNotFound, "User chưa được gán thiết bị hoặc không tìm thấy agent_id")
		return
	}
	var allLogs []interface{}
	for _, agentID := range agentIDs {
		logs, err := logService.GetLogsByAgentID(agentID)
		if err == nil && logs != nil {
			logutil.Debug("GetMyDeviceLogHandler: found %d logs for agentID=%s", len(logs), agentID)
			for _, l := range logs {
				allLogs = append(allLogs, l)
			}
		}
	}
	logutil.Debug("GetMyDeviceLogHandler: total logs returned: %d", len(allLogs))
	response.Success(c, allLogs)
}
