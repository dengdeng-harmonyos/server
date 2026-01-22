package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dengdeng-harmenyos/server/internal/service"
	"github.com/gin-gonic/gin"
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
// GET /api/messages/pending
func (h *MessageHandler) GetPendingMessages(c *gin.Context) {
	// 从请求头获取 device_key
	deviceKey := c.GetHeader("X-Device-Key")
	if deviceKey == "" {
		// 也尝试从 Authorization header 获取
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			deviceKey = authHeader[7:]
		}
	}

	if deviceKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing device key",
		})
		return
	}

	// 查询未投递的消息
	rows, err := h.db.Query(`
		SELECT id::TEXT, server_name, encrypted_aes_key, encrypted_content, iv, 
		       EXTRACT(EPOCH FROM created_at)::BIGINT * 1000 as timestamp
		FROM pending_messages
		WHERE device_key = $1 
		  AND delivered = false 
		  AND expires_at > NOW()
		ORDER BY created_at ASC
		LIMIT 100
	`, deviceKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to query messages: " + err.Error(),
		})
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
			WHERE device_key = $1 AND delivered = false
		`, deviceKey)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": messages,
		"count":    len(messages),
	})
}

// ConfirmMessagesRequest 确认消息请求
type ConfirmMessagesRequest struct {
	MessageIDs []string `json:"messageIds" binding:"required"`
}

// ConfirmMessages 确认消息已收到
// POST /api/messages/confirm
func (h *MessageHandler) ConfirmMessages(c *gin.Context) {
	deviceKey := c.GetHeader("X-Device-Key")
	if deviceKey == "" {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			deviceKey = authHeader[7:]
		}
	}

	if deviceKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing device key",
		})
		return
	}

	var req ConfirmMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	if len(req.MessageIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success":        true,
			"confirmedCount": 0,
		})
		return
	}

	// 构建 SQL IN 子句 - 将字符串ID转换为整数数组
	query := `
		UPDATE pending_messages 
		SET delivered = true, confirmed_at = $1
		WHERE device_key = $2 AND id::TEXT = ANY($3)
	`

	result, err := h.db.Exec(query, time.Now(), deviceKey, req.MessageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to confirm messages: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"confirmedCount": rowsAffected,
	})
}

// SaveEncryptedMessage 保存加密消息到数据库
func (h *MessageHandler) SaveEncryptedMessage(
	deviceKey string,
	serverName string,
	encryptedMsg *service.EncryptedMessage,
) error {
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天后过期

	_, err := h.db.Exec(`
		INSERT INTO pending_messages 
		(device_key, server_name, encrypted_aes_key, encrypted_content, iv, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, deviceKey, serverName, encryptedMsg.EncryptedAESKey,
		encryptedMsg.EncryptedContent, encryptedMsg.IV, expiresAt)

	return err
}
