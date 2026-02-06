package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/dengdeng-harmonyos/server/internal/database"
	"github.com/dengdeng-harmonyos/server/internal/logger"
	"github.com/dengdeng-harmonyos/server/internal/models"
	"github.com/dengdeng-harmonyos/server/internal/service"
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
