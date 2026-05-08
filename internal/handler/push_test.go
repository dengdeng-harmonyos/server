package handler

import (
	"strings"
	"testing"
	"time"
)

func TestParseNotificationDataAndExtractURL(t *testing.T) {
	data, err := parseNotificationData(`[{"key":"__url","value":"https://example.com"},{"key":"tag","value":"work"}]`)
	if err != nil {
		t.Fatalf("parseNotificationData returned error: %v", err)
	}

	if got := extractMessageURL(data); got != "https://example.com" {
		t.Fatalf("extractMessageURL = %q, want %q", got, "https://example.com")
	}
}

func TestParseNotificationDataRejectsObjectShape(t *testing.T) {
	if _, err := parseNotificationData(`{"key":"tag","value":"work"}`); err == nil {
		t.Fatal("parseNotificationData accepted object shape")
	}
}

func TestValidateMessageURLAllowsHTTPAndDeepLinks(t *testing.T) {
	validURLs := []string{
		"https://example.com/path?q=1",
		"http://example.com/path?q=1",
		"myapp://page/detail?id=1",
		"app://open/path",
	}

	for _, validURL := range validURLs {
		if err := validateMessageURL(validURL); err != nil {
			t.Fatalf("validateMessageURL rejected %q: %v", validURL, err)
		}
	}
}

func TestValidateMessageURLRejectsUnsafeSchemes(t *testing.T) {
	invalidURLs := []string{
		"file:///tmp/a",
		"javascript:alert(1)",
		"data:text/html,<b>test</b>",
		"content://contacts/1",
		"tel:10086",
		"sms:10086",
		"mailto:test@example.com",
		"facetime:user@example.com",
		"intent://open/path",
		"market://details?id=app",
		"settings://wifi",
		"app-settings://dengdeng",
		"hmos-settings://wifi",
		"ohos://settings",
	}

	for _, invalidURL := range invalidURLs {
		if err := validateMessageURL(invalidURL); err == nil {
			t.Fatalf("validateMessageURL accepted unsafe URL %q", invalidURL)
		}
	}
}

func TestValidateMessageURLRejectsMalformedInput(t *testing.T) {
	tooLongURL := "myapp://" + strings.Repeat("a", maxMessageURLLength)
	invalidURLs := []string{
		"example.com/path",
		"://missing-scheme",
		"1app://open/path",
		"my_app://open/path",
		"https:example.com",
		"https:///path",
		" https://example.com",
		"https://example.com ",
		"https://exa mple.com",
		"https://example.com/\nnext",
		tooLongURL,
	}

	for _, invalidURL := range invalidURLs {
		if err := validateMessageURL(invalidURL); err == nil {
			t.Fatalf("validateMessageURL accepted malformed URL %q", invalidURL)
		}
	}
}

func TestBackgroundPushWakeCutoffUsesThirtyMinuteCooldown(t *testing.T) {
	now := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	want := time.Date(2026, 5, 7, 11, 30, 0, 0, time.UTC)

	if got := backgroundPushWakeCutoff(now); !got.Equal(want) {
		t.Fatalf("backgroundPushWakeCutoff = %s, want %s", got, want)
	}
}
