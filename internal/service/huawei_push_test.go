package service

import "testing"

func TestBuildAlertNotificationUsesTitleWithoutPrefix(t *testing.T) {
	data := map[string]interface{}{"__url": "https://example.com"}
	notification := buildAlertNotification("通知标题", "通知内容", data)

	if notification.Title != "通知标题" {
		t.Fatalf("notification title = %q, want %q", notification.Title, "通知标题")
	}
	if notification.Body != "通知内容" {
		t.Fatalf("notification body = %q, want %q", notification.Body, "通知内容")
	}
	if notification.ClickAction.Data["__url"] != "https://example.com" {
		t.Fatalf("notification click URL = %v, want https://example.com", notification.ClickAction.Data["__url"])
	}
}
