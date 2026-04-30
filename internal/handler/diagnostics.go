package handler

import (
	"database/sql"
	"net/http"

	"github.com/dengdeng-harmonyos/server/internal/logger"
	"github.com/dengdeng-harmonyos/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiagnosticsHandler struct {
	db *sql.DB
}

type DeviceDiagnosticsResponse struct {
	Exists              bool   `json:"exists"`
	HasPublicKey        bool   `json:"hasPublicKey"`
	IsActive            bool   `json:"isActive"`
	LastActiveAt        string `json:"lastActiveAt"`
	PendingMessageCount int64  `json:"pendingMessageCount"`
}

func NewDiagnosticsHandler(db *sql.DB) *DiagnosticsHandler {
	return &DiagnosticsHandler{db: db}
}

// Device returns non-sensitive device diagnostics for troubleshooting.
func (h *DiagnosticsHandler) Device(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "device_id is required")
		return
	}

	if _, err := uuid.Parse(deviceID); err != nil {
		RespondError(c, http.StatusBadRequest, models.InvalidParams, "Invalid device_id format")
		return
	}

	response := DeviceDiagnosticsResponse{}
	err := h.db.QueryRow(`
		SELECT
			(public_key IS NOT NULL AND public_key <> '') AS has_public_key,
			is_active,
			to_char(last_active_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') AS last_active_at
		FROM devices
		WHERE device_id = $1
	`, deviceID).Scan(&response.HasPublicKey, &response.IsActive, &response.LastActiveAt)
	if err == sql.ErrNoRows {
		RespondSuccess(c, http.StatusOK, response)
		return
	}
	if err != nil {
		logger.ErrorWithStack(err, "Failed to query diagnostics for device: %s", deviceID)
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to query diagnostics")
		return
	}

	response.Exists = true
	if err := h.db.QueryRow(`
		SELECT COUNT(*)
		FROM pending_messages
		WHERE device_id = $1
		  AND delivered = false
		  AND expires_at > NOW()
	`, deviceID).Scan(&response.PendingMessageCount); err != nil {
		logger.ErrorWithStack(err, "Failed to query pending count for device: %s", deviceID)
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to query diagnostics")
		return
	}

	RespondSuccess(c, http.StatusOK, response)
}
