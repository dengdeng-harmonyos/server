package handler

import (
	"database/sql"
	"net/http"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/dengdeng-harmonyos/server/internal/models"
	"github.com/gin-gonic/gin"
)

type AppUpdateHandler struct {
	db       *sql.DB
	fallback config.AppUpdateConfig
}

type AppUpdateRequest struct {
	VersionCode int64  `form:"version_code"`
	VersionName string `form:"version_name"`
}

type AppUpdatePolicy struct {
	LatestVersionCode int64
	LatestVersionName string
	MinVersionCode    int64
	ForceUpdate       bool
	StoreURL          string
	ReleaseNotes      string
}

func NewAppUpdateHandler(db *sql.DB, fallback config.AppUpdateConfig) *AppUpdateHandler {
	return &AppUpdateHandler{db: db, fallback: fallback}
}

func (h *AppUpdateHandler) Check(c *gin.Context) {
	var req AppUpdateRequest
	_ = c.ShouldBindQuery(&req)

	policy, err := h.loadPolicy(c)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, models.SystemError, "Failed to load app update policy")
		return
	}

	hasNewVersion := policy.LatestVersionCode > 0 &&
		req.VersionCode > 0 &&
		req.VersionCode < policy.LatestVersionCode

	mustUpdate := false
	if policy.MinVersionCode > 0 && req.VersionCode > 0 && req.VersionCode < policy.MinVersionCode {
		mustUpdate = true
	}
	if policy.ForceUpdate && hasNewVersion {
		mustUpdate = true
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"currentVersionCode": req.VersionCode,
		"currentVersionName": req.VersionName,
		"latestVersionCode":  policy.LatestVersionCode,
		"latestVersionName":  policy.LatestVersionName,
		"minVersionCode":     policy.MinVersionCode,
		"hasNewVersion":      hasNewVersion,
		"mustUpdate":         mustUpdate,
		"storeUrl":           policy.StoreURL,
		"releaseNotes":       policy.ReleaseNotes,
	})
}

func (h *AppUpdateHandler) loadPolicy(c *gin.Context) (AppUpdatePolicy, error) {
	const query = `
		SELECT latest_version_code, latest_version_name, min_version_code, force_update, store_url, release_notes
		FROM app_update_policies
		WHERE platform = $1 AND enabled = TRUE
		LIMIT 1
	`

	var policy AppUpdatePolicy
	err := h.db.QueryRowContext(c.Request.Context(), query, "harmonyos").Scan(
		&policy.LatestVersionCode,
		&policy.LatestVersionName,
		&policy.MinVersionCode,
		&policy.ForceUpdate,
		&policy.StoreURL,
		&policy.ReleaseNotes,
	)
	if err == nil {
		return policy, nil
	}
	if err != sql.ErrNoRows {
		return AppUpdatePolicy{}, err
	}

	return AppUpdatePolicy{
		LatestVersionCode: h.fallback.LatestVersionCode,
		LatestVersionName: h.fallback.LatestVersionName,
		MinVersionCode:    h.fallback.MinVersionCode,
		ForceUpdate:       h.fallback.ForceUpdate,
		StoreURL:          h.fallback.StoreURL,
		ReleaseNotes:      h.fallback.ReleaseNotes,
	}, nil
}
