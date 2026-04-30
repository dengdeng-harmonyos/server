package handler

import "testing"

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

func TestValidateMessageURLOnlyAllowsHTTP(t *testing.T) {
	if err := validateMessageURL("https://example.com/path?q=1"); err != nil {
		t.Fatalf("validateMessageURL rejected https URL: %v", err)
	}

	if err := validateMessageURL("file:///tmp/test"); err == nil {
		t.Fatal("validateMessageURL accepted non-http URL")
	}
}
