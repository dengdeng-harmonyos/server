package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/dengdeng-harmonyos/server/internal/database"
	"github.com/dengdeng-harmonyos/server/internal/logger"
	"github.com/dengdeng-harmonyos/server/internal/models"
	"github.com/dengdeng-harmonyos/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxMessageURLLength        = 2048
	backgroundPushWakeCooldown = 30 * time.Minute
)

var blockedMessageURLSchemes = map[string]struct{}{
	"app-settings":  {},
	"content":       {},
	"data":          {},
	"facetime":      {},
	"file":          {},
	"hmos-settings": {},
	"intent":        {},
	"javascript":    {},
	"mailto":        {},
	"market":        {},
	"ohos":          {},
	"settings":      {},
	"sms":           {},
	"tel":           {},
}

type PushHandler struct {
	db             *database.Database
	pushService    *service.HuaweiPushService
	deviceHandler  *DeviceHandler
	serverName     string
	cryptoService  *service.CryptoService
	messageHandler *MessageHandler
}

type backgroundSyncSignal struct {
	Type       string `json:"type"`
	ServerName string `json:"server_name"`
	CreatedAt  string `json:"created_at"`
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
// GET /api/v1/push/notification?device_id=xxx&title=xxx&content=xxx&data=[{"key":"__url","value":"https://example.com"}]
func (h *PushHandler) SendNotification(c *gin.Context) {
	var req models.PushNotificationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid request: "+err.Error())
		return
	}

	// 验证 device_id 格式是否为有效的 UUID
	if _, err := uuid.Parse(req.DeviceId); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid device_id format")
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

	dataArray, err := parseNotificationData(req.Data)
	if err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, err.Error())
		return
	}

	messageURL := extractMessageURL(dataArray)
	if err := validateMessageURL(messageURL); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, err.Error())
		return
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

	// 2. 先保存加密消息，确保后台唤醒或普通通知到达时 App 已有 pending 可拉取。
	err = h.messageHandler.SaveEncryptedMessage(req.DeviceId, h.serverName, encryptedMsg)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to save encrypted message for device: %s", req.DeviceId)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to save message: "+err.Error())
		return
	}

	// 3. 有 pending 消息时发送一次低频后台唤醒信号，失败不影响普通通知。
	h.maybeSendBackgroundSyncSignal(req.DeviceId, pushToken)

	// 4. 发送华为推送通知（明文内容，用于显示通知）
	notificationData := map[string]interface{}{
		"type":          "new_message",
		"server_name":   h.serverName,
		"__server_name": h.serverName,
	}
	if messageURL != "" {
		notificationData["__url"] = messageURL
	}
	err = h.pushService.SendNotification(pushToken, req.Title, req.Content, notificationData)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to send push notification for device: %s", req.DeviceId)
		RespondError(c, http.StatusInternalServerError, models.OperationFailed, "Failed to send notification: "+err.Error())
		return
	}

	logger.Info("Successfully sent notification to device: %s, title: %s", req.DeviceId, req.Title)

	RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Notification sent successfully",
	})
}

func (h *PushHandler) maybeSendBackgroundSyncSignal(deviceID string, pushToken string) {
	now := time.Now().UTC()
	shouldSend, err := h.reserveBackgroundPushWake(deviceID, now)
	if err != nil {
		logger.ErrorWithStack(err, "Failed to reserve background push wake for device: %s", deviceID)
		return
	}
	if !shouldSend {
		logger.Info("Skipped background push wake for device: %s within cooldown window", deviceID)
		return
	}

	payload, err := json.Marshal(backgroundSyncSignal{
		Type:       "sync_pending",
		ServerName: h.serverName,
		CreatedAt:  now.Format(time.RFC3339),
	})
	if err != nil {
		logger.ErrorWithStack(err, "Failed to build background sync signal for device: %s", deviceID)
		return
	}

	if err := h.pushService.SendBackgroundMessage(pushToken, string(payload)); err != nil {
		logger.ErrorWithStack(err, "Failed to send background sync signal for device: %s", deviceID)
		return
	}
	logger.Info("Background sync signal sent for device: %s", deviceID)
}

func (h *PushHandler) reserveBackgroundPushWake(deviceID string, now time.Time) (bool, error) {
	cutoff := backgroundPushWakeCutoff(now)
	result, err := h.db.DB.Exec(`
		UPDATE devices
		SET last_background_push_attempt_at = $3,
			updated_at = NOW()
		WHERE device_id = $1
			AND is_active = TRUE
			AND EXISTS (
				SELECT 1
				FROM pending_messages
				WHERE device_id = $1
					AND delivered = false
					AND expires_at > NOW()
			)
			AND (
				last_background_push_attempt_at IS NULL
				OR last_background_push_attempt_at <= $2
			)
	`, deviceID, cutoff, now)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func backgroundPushWakeCutoff(now time.Time) time.Time {
	return now.UTC().Add(-backgroundPushWakeCooldown)
}

func parseNotificationData(rawData string) ([]map[string]interface{}, error) {
	var dataArray []map[string]interface{}
	if rawData == "" {
		return dataArray, nil
	}

	if err := json.Unmarshal([]byte(rawData), &dataArray); err != nil {
		return nil, errInvalidNotificationData()
	}

	return dataArray, nil
}

func extractMessageURL(dataArray []map[string]interface{}) string {
	for _, item := range dataArray {
		keyValue, ok := item["key"].(string)
		if !ok || keyValue != "__url" {
			continue
		}

		if urlValue, ok := item["value"].(string); ok {
			return urlValue
		}
	}
	return ""
}

func validateMessageURL(messageURL string) error {
	if messageURL == "" {
		return nil
	}

	if messageURL != strings.TrimSpace(messageURL) ||
		len(messageURL) > maxMessageURLLength ||
		containsUnsafeURLCharacter(messageURL) {
		return errInvalidMessageURL()
	}

	parsedURL, err := url.ParseRequestURI(messageURL)
	if err != nil || parsedURL.Scheme == "" {
		return errInvalidMessageURL()
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if !isValidMessageURLScheme(scheme) {
		return errInvalidMessageURL()
	}
	if (scheme == "http" || scheme == "https") && parsedURL.Host == "" {
		return errInvalidMessageURL()
	}
	if _, blocked := blockedMessageURLSchemes[scheme]; blocked {
		return errInvalidMessageURL()
	}

	return nil
}

func containsUnsafeURLCharacter(messageURL string) bool {
	for _, r := range messageURL {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func isValidMessageURLScheme(scheme string) bool {
	if scheme == "" || !isASCIIAlpha(scheme[0]) {
		return false
	}

	for i := 1; i < len(scheme); i++ {
		c := scheme[i]
		if isASCIIAlpha(c) || isASCIIDigit(c) || c == '+' || c == '-' || c == '.' {
			continue
		}
		return false
	}

	return true
}

func isASCIIAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isASCIIDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func errInvalidNotificationData() error {
	return &pushValidationError{message: "Invalid data format, expected array of {key, value}"}
}

func errInvalidMessageURL() error {
	return &pushValidationError{message: "Invalid __url format, expected a safe http(s) URL or app URL scheme"}
}

type pushValidationError struct {
	message string
}

func (e *pushValidationError) Error() string {
	return e.message
}
