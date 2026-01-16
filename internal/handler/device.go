package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/dangdangdang-push-server/internal/config"
	"github.com/yourusername/dangdangdang-push-server/internal/database"
	"github.com/yourusername/dangdangdang-push-server/internal/models"
	"github.com/yourusername/dangdangdang-push-server/internal/service"
)

type DeviceHandler struct {
	db         *database.Database
	encryption *service.EncryptionService
	config     config.SecurityConfig
}

func NewDeviceHandler(db *database.Database, cfg config.Config) (*DeviceHandler, error) {
	encryption, err := service.NewEncryptionService(cfg.Security.EncryptionKey)
	if err != nil {
		return nil, err
	}

	return &DeviceHandler{
		db:         db,
		encryption: encryption,
		config:     cfg.Security,
	}, nil
}

// Register 设备注册接口
func (h *DeviceHandler) Register(c *gin.Context) {
	var req models.DeviceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	// 加密push_token
	encryptedToken, err := h.encryption.Encrypt(req.PushToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to encrypt token",
		})
		return
	}

	// 查询是否已存在该push_token
	var existingDevice models.Device
	err = h.db.DB.QueryRow(`
		SELECT id, device_key FROM devices WHERE push_token = $1
	`, encryptedToken).Scan(&existingDevice.ID, &existingDevice.DeviceKey)

	if err == nil {
		// 设备已存在，更新信息
		_, err = h.db.DB.Exec(`
			UPDATE devices 
			SET device_type = $1, os_version = $2, app_version = $3, 
			    is_active = true, last_active_at = NOW(), updated_at = NOW()
			WHERE id = $4
		`, req.DeviceType, req.OSVersion, req.AppVersion, existingDevice.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to update device",
			})
			return
		}

		c.JSON(http.StatusOK, models.DeviceRegisterResponse{
			Success:   true,
			DeviceKey: existingDevice.DeviceKey,
			Message:   "Device updated successfully",
		})
		return
	}

	// 生成新的device_key
	deviceKey := uuid.New().String()

	// 插入新设备
	_, err = h.db.DB.Exec(`
		INSERT INTO devices (device_key, push_token, device_type, os_version, app_version, is_active, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW(), NOW())
	`, deviceKey, encryptedToken, req.DeviceType, req.OSVersion, req.AppVersion)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to register device",
		})
		return
	}

	c.JSON(http.StatusOK, models.DeviceRegisterResponse{
		Success:   true,
		DeviceKey: deviceKey,
		Message:   "Device registered successfully",
	})
}

// UpdateToken 更新Push Token
func (h *DeviceHandler) UpdateToken(c *gin.Context) {
	var req struct {
		DeviceKey    string `json:"device_key" binding:"required"`
		NewPushToken string `json:"new_push_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	// 加密新token
	encryptedToken, err := h.encryption.Encrypt(req.NewPushToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to encrypt token",
		})
		return
	}

	// 更新token
	result, err := h.db.DB.Exec(`
		UPDATE devices 
		SET push_token = $1, last_active_at = NOW(), updated_at = NOW()
		WHERE device_key = $2 AND is_active = true
	`, encryptedToken, req.DeviceKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update token",
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
		"message": "Token updated successfully",
	})
}

// Deactivate 停用设备
func (h *DeviceHandler) Deactivate(c *gin.Context) {
	deviceKey := c.Query("device_key")
	if deviceKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "device_key is required",
		})
		return
	}

	_, err := h.db.DB.Exec(`
		UPDATE devices SET is_active = false, updated_at = NOW()
		WHERE device_key = $1
	`, deviceKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to deactivate device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device deactivated successfully",
	})
}

// GetPushToken 内部方法：根据device_key获取push_token
func (h *DeviceHandler) GetPushToken(deviceKey string) (string, error) {
	var encryptedToken string
	err := h.db.DB.QueryRow(`
		SELECT push_token FROM devices 
		WHERE device_key = $1 AND is_active = true
	`, deviceKey).Scan(&encryptedToken)

	if err != nil {
		return "", err
	}

	// 解密token
	return h.encryption.Decrypt(encryptedToken)
}

// GetPushTokens 批量获取push_token
func (h *DeviceHandler) GetPushTokens(deviceKeys []string) ([]string, error) {
	if len(deviceKeys) == 0 {
		return []string{}, nil
	}

	tokens := make([]string, 0, len(deviceKeys))
	for _, key := range deviceKeys {
		token, err := h.GetPushToken(key)
		if err == nil {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}
