package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dengdeng-harmenyos/server/internal/config"
	"github.com/dengdeng-harmenyos/server/internal/database"
	"github.com/dengdeng-harmenyos/server/internal/logger"
	"github.com/dengdeng-harmenyos/server/internal/models"
	"github.com/dengdeng-harmenyos/server/internal/service"
	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	db             *database.Database
	pushService    *service.HuaweiPushService
	deviceHandler  *DeviceHandler
	serverName     string
	cryptoService  *service.CryptoService
	messageHandler *MessageHandler
}

func NewPushHandler(db *database.Database, deviceHandler *DeviceHandler, cfg config.HuaweiPushConfig, serverName string) (*PushHandler, error) {
	pushService, err := service.NewHuaweiPushService(cfg)
	if err != nil {
		return nil, err
	}

	return &PushHandler{
		db:             db,
		pushService:    pushService,
		deviceHandler:  deviceHandler,
		serverName:     serverName,
		cryptoService:  service.NewCryptoService(),
		messageHandler: NewMessageHandler(db.DB),
	}, nil
}

// SendNotification 发送通知消息（GET方式）
// GET /api/v1/push/notification?device_id=xxx&title=xxx&content=xxx&data=[{"key":"xxx","value":"xxx"}]
func (h *PushHandler) SendNotification(c *gin.Context) {
	var req models.PushNotificationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 根据device_id获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	// 解析额外数据（保持数组格式）
	var dataArray []map[string]interface{}
	if req.Data != "" {
		// 解析为数组格式 [{"key":"xxx", "value":"xxx"}]
		if err := json.Unmarshal([]byte(req.Data), &dataArray); err != nil {
			RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid data format, expected array of {key, value}")
			return
		}
	}

	// 获取设备公钥
	publicKey, err := h.deviceHandler.GetPublicKey(req.DeviceId)
	if err != nil || publicKey == "" {
		// 必须有公钥才能处理推送
		RespondError(c, http.StatusBadRequest, models.OperationFailed, "Device public key not found, please register device first")
		return
	}

	// 1. 加密消息内容
	messageContent := service.MessageContent{
		Title:      req.Title,
		Content:    req.Content,
		Data:       dataArray,
		ServerName: h.serverName,
	}
	encryptedMsg, err := h.cryptoService.EncryptMessage(publicKey, messageContent)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to encrypt message for device: %s", req.DeviceId)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send notification: "+err.Error())
		return
	}

	// 2. 发送华为推送通知（明文内容，用于显示通知）
	notificationData := map[string]interface{}{
		"type":        "new_message",
		"server_name": h.serverName,
	}
	err = h.pushService.SendNotification(pushToken, req.Title, req.Content, notificationData)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to send push notification for device: %s", req.DeviceId)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send notification: "+err.Error())
		return
	}

	// 3. 保存加密消息到数据库
	err = h.messageHandler.SaveEncryptedMessage(req.DeviceId, h.serverName, encryptedMsg)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to save encrypted message for device: %s", req.DeviceId)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to save message: "+err.Error())
		return
	}

	logger.Info("Successfully sent notification to device: %s, title: %s", req.DeviceId, req.Title)

	// 更新统计
	h.updateStatistics("notification", true)

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Notification sent successfully",
	})
}

// SendFormUpdate 发送卡片刷新消息（GET方式）
// GET /api/v1/push/form?device_id=xxx&form_id=xxx&form_data={"key":"value"}
func (h *PushHandler) SendFormUpdate(c *gin.Context) {
	var req models.FormUpdateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 根据device_id获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceId)
	if err != nil {
		RespondError(c, http.StatusNotFound, models.DataNotFound, "Device not found")
		return
	}

	// 解析表单数据
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(req.FormData), &formData); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid form_data format")
		return
	}

	// 将formID从string转换为int64
	formID := int64(0)
	if _, err := fmt.Sscanf(req.FormID, "%d", &formID); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid form_id format, must be a number")
		return
	}

	// 发送推送（使用简化版本）
	err = h.pushService.SendFormUpdateSimple(pushToken, formID, formData)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send form update: "+err.Error())
		return
	}

	// 更新统计
	h.updateStatistics("form", true)

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Form update sent successfully",
	})
}

// SendBackgroundMessage 发送后台消息（GET方式）
// GET /api/v1/push/background?device_id=xxx&data={"action":"sync"}
func (h *PushHandler) SendBackgroundMessage(c *gin.Context) {
	var req models.BackgroundPushRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 根据device_id获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceId)
	if err != nil {
		RespondError(c, http.StatusNotFound, models.DataNotFound, "Device not found")
		return
	}

	// 发送推送
	err = h.pushService.SendBackgroundMessage(pushToken, req.Data)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send background message: "+err.Error())
		return
	}

	// 更新统计
	h.updateStatistics("background", true)

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Background message sent successfully",
	})
}

// SendBatch 批量发送通知消息（GET方式）
// GET /api/v1/push/batch?device_ids=key1,key2,key3&title=xxx&body=xxx&data={"key":"value"}
func (h *PushHandler) SendBatch(c *gin.Context) {
	var req models.BatchPushRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 解析device_ids（逗号分隔）
	deviceIds := strings.Split(req.DeviceIds, ",")
	if len(deviceIds) == 0 {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "No device keys provided")
		return
	}

	if len(deviceIds) > 1000 {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Too many devices (max 1000)")
		return
	}

	// 批量获取push_token
	pushTokens, err := h.deviceHandler.GetPushTokens(deviceIds)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to get device tokens")
		return
	}

	if len(pushTokens) == 0 {
		RespondError(c, http.StatusNotFound, models.DataNotFound, "No valid devices found")
		return
	}

	// 解析额外数据
	var data map[string]interface{}
	if req.Data != "" {
		if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
			RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid data format")
			return
		}
	}

	// 批量发送推送
	err = h.pushService.SendBatchNotification(pushTokens, req.Title, req.Body, data)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send batch push: "+err.Error())
		return
	}

	// 更新统计
	for i := 0; i < len(pushTokens); i++ {
		h.updateStatistics("notification", true)
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message":      "Batch notification sent successfully",
		"total_sent":   len(pushTokens),
		"total_failed": len(deviceIds) - len(pushTokens),
	})
}

// GetStatistics 获取推送统计数据
// GET /api/v1/push/statistics?date=2026-01-13
func (h *PushHandler) GetStatistics(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	rows, err := h.db.DB.Query(`
		SELECT push_type, total_count, success_count, failed_count
		FROM push_statistics
		WHERE date = $1
	`, dateStr)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to get statistics")
		return
	}
	defer rows.Close()

	stats := make([]map[string]interface{}, 0)
	for rows.Next() {
		var pushType string
		var total, success, failed int
		if err := rows.Scan(&pushType, &total, &success, &failed); err == nil {
			stats = append(stats, map[string]interface{}{
				"push_type":     pushType,
				"total_count":   total,
				"success_count": success,
				"failed_count":  failed,
			})
		}
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"date":  dateStr,
		"stats": stats,
	})
}

// updateStatistics 更新推送统计（内部方法）
func (h *PushHandler) updateStatistics(pushType string, success bool) {
	date := time.Now().Format("2006-01-02")

	if success {
		h.db.DB.Exec(`
			INSERT INTO push_statistics (date, push_type, total_count, success_count, failed_count, created_at)
			VALUES ($1, $2, 1, 1, 0, NOW())
			ON CONFLICT (date, push_type) 
			DO UPDATE SET 
				total_count = push_statistics.total_count + 1,
				success_count = push_statistics.success_count + 1
		`, date, pushType)
	} else {
		h.db.DB.Exec(`
			INSERT INTO push_statistics (date, push_type, total_count, success_count, failed_count, created_at)
			VALUES ($1, $2, 1, 0, 1, NOW())
			ON CONFLICT (date, push_type) 
			DO UPDATE SET 
				total_count = push_statistics.total_count + 1,
				failed_count = push_statistics.failed_count + 1
		`, date, pushType)
	}
}
