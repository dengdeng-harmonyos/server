package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestConfirmMessagesEmptyListDoesNotTouchDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/confirm", NewMessageHandler(nil).ConfirmMessages)

	req := httptest.NewRequest(http.MethodPost, "/confirm", strings.NewReader(`{
		"device_id": "d5e2a0a0-36a8-4d8b-bcb7-469c7f09fc61",
		"messageIds": []
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.Code, resp.Body.String())
	}
	if !strings.Contains(resp.Body.String(), `"confirmedCount":0`) {
		t.Fatalf("response does not contain confirmedCount=0: %s", resp.Body.String())
	}
}
