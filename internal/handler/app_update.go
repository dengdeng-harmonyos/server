package handler

import (
	"net/http"

	"github.com/dengdeng-harmonyos/server/internal/config"
	"github.com/gin-gonic/gin"
)

type AppUpdateHandler struct {
	cfg config.AppUpdateConfig
}

type AppUpdateRequest struct {
	VersionCode int64  `form:"version_code"`
	VersionName string `form:"version_name"`
}

func NewAppUpdateHandler(cfg config.AppUpdateConfig) *AppUpdateHandler {
	return &AppUpdateHandler{cfg: cfg}
}

func (h *AppUpdateHandler) Check(c *gin.Context) {
	var req AppUpdateRequest
	_ = c.ShouldBindQuery(&req)

	hasNewVersion := h.cfg.LatestVersionCode > 0 &&
		req.VersionCode > 0 &&
		req.VersionCode < h.cfg.LatestVersionCode

	mustUpdate := false
	if h.cfg.MinVersionCode > 0 && req.VersionCode > 0 && req.VersionCode < h.cfg.MinVersionCode {
		mustUpdate = true
	}
	if h.cfg.ForceUpdate && hasNewVersion {
		mustUpdate = true
	}

	RespondSuccess(c, http.StatusOK, gin.H{
		"currentVersionCode": req.VersionCode,
		"currentVersionName": req.VersionName,
		"latestVersionCode":  h.cfg.LatestVersionCode,
		"latestVersionName":  h.cfg.LatestVersionName,
		"minVersionCode":     h.cfg.MinVersionCode,
		"hasNewVersion":      hasNewVersion,
		"mustUpdate":         mustUpdate,
		"storeUrl":           h.cfg.StoreURL,
		"releaseNotes":       h.cfg.ReleaseNotes,
	})
}
