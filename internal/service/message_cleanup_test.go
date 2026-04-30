package service

import (
	"strings"
	"testing"
)

func TestExpiredMessageCleanupSQLOnlyTouchesPendingMessages(t *testing.T) {
	if !strings.Contains(expiredMessageCleanupSQL, "pending_messages") {
		t.Fatal("cleanup query must delete from pending_messages")
	}
	if strings.Contains(expiredMessageCleanupSQL, "push_statistics") {
		t.Fatal("cleanup query must not touch or expose push statistics")
	}
}
