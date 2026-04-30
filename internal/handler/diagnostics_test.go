package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeviceDiagnosticsRejectsInvalidDeviceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/diagnostics/device", NewDiagnosticsHandler(nil).Device)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/diagnostics/device?device_id=not-a-uuid", nil)
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", resp.Code, resp.Body.String())
	}
	if !strings.Contains(resp.Body.String(), "Invalid device_id format") {
		t.Fatalf("unexpected response body: %s", resp.Body.String())
	}
}
