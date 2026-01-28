package handler

import (
	"database/sql"
	"net/http"
	"time"
	"github.com/dengdeng-harmenyos/server/internal/logger"
	"github.com/dengdeng-harmenyos/server/internal/models"
	"github.com/dengdeng-harmenyos/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	db            *sql.DB
	cryptoService *service.CryptoService
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(db *sql.DB) *MessageHandler {
	return &MessageHandler{
		db:            db,
		cryptoService: service.NewCryptoService(),
	}
}

// PendingMessage 待接收消息响应
type PendingMessage struct {
	ID               string `json:"id"`
	ServerName       string `json:"serverName"`
	EncryptedAESKey  string `json:"encryptedAESKey"`
	EncryptedContent string `json:"encryptedContent"`
	IV               string `json:"iv"`
	Timestamp        int64  `json:"timestamp"`
}

// GetPendingMessages 获取待接收的消息
// GET /api/messages/pending?device_id=xxx
func (h *MessageHandler) GetPendingMessages(c *gin.Context) {
	// 从 query 获取 device_id
	deviceId := c.Query("device_id")

	if deviceId == "" {
		RespondError(c, http.StatusUnauthorized, models.Unauthorized, "Missing device key")
		return
	}

	// 查询未投递的消息
	rows, err := h.db.Query(`
		SELECT id::TEXT, server_name, encrypted_aes_key, encrypted_content, iv, 
		       EXTRACT(EPOCH FROM created_at)::BIGINT * 1000 as timestamp
		FROM pending_messages
		WHERE device_id = $1 
		  AND delivered = false 
		  AND expires_at > NOW()
		ORDER BY created_at ASC
		LIMIT 100
	`, deviceId)

	if err != nil {
		logger.ErrorWithStack(err, "Failed to query pending messages for device: %s", deviceId)
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to query messages: "+err.Error())
		return
	}
	defer rows.Close()

	messages := []PendingMessage{}
	for rows.Next() {
		var msg PendingMessage
		if err := rows.Scan(&msg.ID, &msg.ServerName, &msg.EncryptedAESKey,
			&msg.EncryptedContent, &msg.IV, &msg.Timestamp); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// 标记为已发送
	if len(messages) > 0 {
		_, _ = h.db.Exec(`
			UPDATE pending_messages 
			SET notification_sent = true 
			WHERE device_id = $1 AND delivered = false
		`, deviceId)
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
	})
}

// ConfirmMessagesRequest 确认消息请求
type ConfirmMessagesRequest struct {
	DeviceId   string   `json:"device_id" binding:"required"`
	MessageIDs []string `json:"messageIds" binding:"required"`
}

// ConfirmMessages 确认消息已收到
// POST /api/messages/confirm
func (h *MessageHandler) ConfirmMessages(c *gin.Context) {
	var req ConfirmMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.ErrorWithStack(err, "Failed to bind confirm request from device: %s", deviceId)
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())

		return
	}

	if req.DeviceId == "" {
		RespondError(c, http.StatusUnauthorized, models.Unauthorized, "Missing device key")
		return
	}

	if len(req.MessageIDs) == 0 {
		RespondSuccess(c, http.StatusOK, gin.H{
			"confirmedCount": 0,
		})
		return
	}

	// 构建 SQL IN 子句 - 使用 pq.Array 将字符串数组转换为 PostgreSQL 数组
	query := `
		DELETE FROM pending_messages
		WHERE device_id = $1 AND id::TEXT = ANY($2)
	`

	result, err := h.db.Exec(query, req.DeviceId, pq.Array(req.MessageIDs))
	if err != nil {

		logger.ErrorWithStack(err, "Failed to confirm messages for device: %s, messageIDs: %v", req.DeviceId, req.MessageIDs)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to confirm messages: "+err.Error())

		return
	}

	rowsAffected, _ := result.RowsAffected()


	logger.Info("Confirmed %d messages for device: %s", rowsAffected, req.DeviceId)
	RespondSuccess(c, http.StatusOK, gin.H{
		"confirmedCount": rowsAffected,
	})
}

// SaveEncryptedMessage 保存加密消息到数据库
func (h *MessageHandler) SaveEncryptedMessage(
	deviceId string,
	serverName string,
	encryptedMsg *service.EncryptedMessage,
) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天后过期

	_, err := h.db.Exec(`
		INSERT INTO pending_messages 
		(device_id, server_name, encrypted_aes_key, encrypted_content, iv, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, deviceId, serverName, encryptedMsg.EncryptedAESKey,
		encryptedMsg.EncryptedContent, encryptedMsg.IV, expiresAt)

	return err
}
