package handler

import (
	"net/http"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/dengdeng-harmonyos/server/internal/database"
	"github.com/dengdeng-harmonyos/server/internal/models"
	"github.com/dengdeng-harmonyos/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceHandler struct {
	db         *database.Database
	encryption *service.EncryptionService
	config     config.SecurityConfig
	serverName string // 服务器名称
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
		serverName: cfg.Server.ServerName,
	}, nil
}

// Register 设备注册接口
func (h *DeviceHandler) Register(c *gin.Context) {
	var req models.DeviceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 加密push_token
	encryptedToken, err := h.encryption.Encrypt(req.PushToken)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to encrypt token")
		return
	}

	// 查询是否已存在该push_token
	var existingDevice models.Device
	err = h.db.DB.QueryRow(`
		SELECT id, device_id FROM devices WHERE push_token = $1
	`, encryptedToken).Scan(&existingDevice.ID, &existingDevice.DeviceId)

	if err == nil {
		// 设备已存在，更新信息（包括公钥）
		_, err = h.db.DB.Exec(`
			UPDATE devices 
			SET device_type = $1, os_version = $2, app_version = $3, public_key = $4,
			    is_active = true, last_active_at = NOW(), updated_at = NOW()
			WHERE id = $5
		`, req.DeviceType, req.OSVersion, req.AppVersion, req.PublicKey, existingDevice.ID)

		if err != nil {
			RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to update device")
			return
		}

		RespondSuccess(c, http.StatusOK, gin.H{
			"device_id":   existingDevice.DeviceId,
			"server_name": h.serverName,
			"message":     "Device updated successfully",
		})
		return
	}

	// 生成新的device_id
	deviceId := uuid.New().String()

	// 插入新设备（包括公钥）
	_, err = h.db.DB.Exec(`
		INSERT INTO devices (device_id, push_token, public_key, device_type, os_version, app_version, is_active, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW(), NOW())
	`, deviceId, encryptedToken, req.PublicKey, req.DeviceType, req.OSVersion, req.AppVersion)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to register device")
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"device_id":   deviceId,
		"server_name": h.serverName,
		"message":     "Device registered successfully",
	})
}

// UpdateToken 更新Push Token
func (h *DeviceHandler) UpdateToken(c *gin.Context) {
	var req struct {
		DeviceId     string `json:"device_id" binding:"required"`
		NewPushToken string `json:"new_push_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request")
		return
	}

	// 加密新token
	encryptedToken, err := h.encryption.Encrypt(req.NewPushToken)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to encrypt token")
		return
	}

	// 更新token
	result, err := h.db.DB.Exec(`
		UPDATE devices 
		SET push_token = $1, last_active_at = NOW(), updated_at = NOW()
		WHERE device_id = $2 AND is_active = true
	`, encryptedToken, req.DeviceId)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to update token")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		RespondError(c, http.StatusNotFound, models.DataNotFound, "Device not found")
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Token updated successfully",
	})
}

// Delete 删除设备及其所有相关数据
func (h *DeviceHandler) Delete(c *gin.Context) {
	deviceId := c.Query("device_id")
	if deviceId == "" {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "device_id is required")
		return
	}

	// 先获取设备的push_token（当前仅用于存在性校验；如未来需要调用华为删除接口再解密使用）
	var encryptedToken string
	err := h.db.DB.QueryRow(`
		SELECT push_token FROM devices 
		WHERE device_id = $1
	`, deviceId).Scan(&encryptedToken)

	if err != nil {
		// 设备不存在
		RespondError(c, http.StatusNotFound, models.DataNotFound, "Device not found")
		return
	}

	// 先删除pending_messages中的相关消息
	_, err = h.db.DB.Exec(`
		DELETE FROM pending_messages WHERE device_id = $1
	`, deviceId)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to delete pending messages")
		return
	}

	// 删除设备记录
	result, err := h.db.DB.Exec(`
		DELETE FROM devices WHERE device_id = $1
	`, deviceId)

	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to delete device")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		RespondError(c, http.StatusNotFound, models.DataNotFound, "Device not found")
		return
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Device deleted successfully",
	})
}

// GetPushToken 内部方法：根据device_id获取push_token
func (h *DeviceHandler) GetPushToken(deviceId string) (string, error) {
	var encryptedToken string
	err := h.db.DB.QueryRow(`
		SELECT push_token FROM devices 
		WHERE device_id = $1 AND is_active = true
	`, deviceId).Scan(&encryptedToken)

	if err != nil {
		return "", err
	}

	// 解密token
	return h.encryption.Decrypt(encryptedToken)
}

// GetPublicKey 内部方法：根据device_id获取public_key
func (h *DeviceHandler) GetPublicKey(deviceId string) (string, error) {
	var publicKey string
	err := h.db.DB.QueryRow(`
		SELECT public_key FROM devices 
		WHERE device_id = $1 AND is_active = true
	`, deviceId).Scan(&publicKey)

	return publicKey, err
}
