package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dengdeng-harmenyos/server/internal/config"
	"github.com/dengdeng-harmenyos/server/internal/database"
	"github.com/dengdeng-harmenyos/server/internal/models"
	"github.com/dengdeng-harmenyos/server/internal/service"
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
// GET /api/v1/push/notification?device_key=xxx&title=xxx&body=xxx&data={"key":"value"}
func (h *PushHandler) SendNotification(c *gin.Context) {
	var req models.PushNotificationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 根据device_key获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	// 解析额外数据
	var data map[string]interface{}
	if req.Data != "" {
		if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid data format",
			})
			return
		}
	} else {
		data = make(map[string]interface{})
	}

	// 添加服务器标识到data中
	data["__server_name"] = h.serverName

	// 获取设备公钥
	publicKey, err := h.deviceHandler.GetPublicKey(req.DeviceKey)
	if err != nil || publicKey == "" {
		// 如果没有公钥，发送普通推送（兼容旧设备）
		err = h.pushService.SendNotification(pushToken, req.Title, req.Body, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to send push: " + err.Error(),
			})
			return
		}
	} else {
		// 有公钥，使用加密存储方案
		// 1. 加密消息内容
		messageContent := service.MessageContent{
			Title:   req.Title,
			Content: req.Body,
			Data:    data,
		}
		encryptedMsg, err := h.cryptoService.EncryptMessage(publicKey, messageContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to encrypt message: " + err.Error(),
			})
			return
		}

		// 2. 保存加密消息到数据库
		err = h.messageHandler.SaveEncryptedMessage(req.DeviceKey, h.serverName, encryptedMsg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to save message: " + err.Error(),
			})
			return
		}

		// 3. 发送华为推送通知（只包含提示信息）
		notificationData := map[string]interface{}{
			"type":        "new_message",
			"server_name": h.serverName,
		}
		err = h.pushService.SendNotification(pushToken, "新消息", "您有新的消息，请打开查看", notificationData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to send notification: " + err.Error(),
			})
			return
		}
	}

	// 更新统计
	h.updateStatistics("notification", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification sent successfully",
	})
}

// SendFormUpdate 发送卡片刷新消息（GET方式）
// GET /api/v1/push/form?device_key=xxx&form_id=xxx&form_data={"key":"value"}
func (h *PushHandler) SendFormUpdate(c *gin.Context) {
	var req models.FormUpdateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 根据device_key获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	// 解析表单数据
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(req.FormData), &formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid form_data format",
		})
		return
	}

	// 将formID从string转换为int64
	formID := int64(0)
	if _, err := fmt.Sscanf(req.FormID, "%d", &formID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid form_id format, must be a number",
		})
		return
	}

	// 发送推送（使用简化版本）
	err = h.pushService.SendFormUpdateSimple(pushToken, formID, formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to send form update: " + err.Error(),
		})
		return
	}

	// 更新统计
	h.updateStatistics("form", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Form update sent successfully",
	})
}

// SendBackgroundMessage 发送后台消息（GET方式）
// GET /api/v1/push/background?device_key=xxx&data={"action":"sync"}
func (h *PushHandler) SendBackgroundMessage(c *gin.Context) {
	var req models.BackgroundPushRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 根据device_key获取push_token
	pushToken, err := h.deviceHandler.GetPushToken(req.DeviceKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	// 发送推送
	err = h.pushService.SendBackgroundMessage(pushToken, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to send background message: " + err.Error(),
		})
		return
	}

	// 更新统计
	h.updateStatistics("background", true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Background message sent successfully",
	})
}

// SendBatch 批量发送通知消息（GET方式）
// GET /api/v1/push/batch?device_keys=key1,key2,key3&title=xxx&body=xxx&data={"key":"value"}
func (h *PushHandler) SendBatch(c *gin.Context) {
	var req models.BatchPushRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 解析device_keys（逗号分隔）
	deviceKeys := strings.Split(req.DeviceKeys, ",")
	if len(deviceKeys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No device keys provided",
		})
		return
	}

	if len(deviceKeys) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Too many devices (max 1000)",
		})
		return
	}

	// 批量获取push_token
	pushTokens, err := h.deviceHandler.GetPushTokens(deviceKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get device tokens",
		})
		return
	}

	if len(pushTokens) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "No valid devices found",
		})
		return
	}

	// 解析额外数据
	var data map[string]interface{}
	if req.Data != "" {
		if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid data format",
			})
			return
		}
	}

	// 批量发送推送
	err = h.pushService.SendBatchNotification(pushTokens, req.Title, req.Body, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to send batch push: " + err.Error(),
		})
		return
	}

	// 更新统计
	for i := 0; i < len(pushTokens); i++ {
		h.updateStatistics("notification", true)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Batch notification sent successfully",
		"total_sent":   len(pushTokens),
		"total_failed": len(deviceKeys) - len(pushTokens),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get statistics",
		})
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"date":    dateStr,
		"stats":   stats,
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

// DeleteDevice 删除设备及其相关数据
func (h *PushHandler) DeleteDevice(c *gin.Context) {
	deviceKey := c.Query("device_key")
	if deviceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "device_key is required",
		})
		return
	}

	// 删除数据库记录（会自动级联删除pending_messages）
	result, err := h.db.DB.Exec(`DELETE FROM devices WHERE device_key = $1`, deviceKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete device",
		})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device deleted successfully",
	})
}
