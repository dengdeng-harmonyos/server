package config

import "testing"

func TestDefaultServerCompatibilityRequiresDeepLinkScheme(t *testing.T) {
	t.Setenv("SERVER_API_VERSION", "")
	t.Setenv("SERVER_CAPABILITIES", "")

	cfg := Load()
	if cfg.Server.APIVersion != 3 {
		t.Fatalf("default API version = %d, want 3", cfg.Server.APIVersion)
	}
	if !containsCapability(cfg.Server.Capabilities, "push_deep_link_scheme") {
		t.Fatalf("default capabilities = %v, want push_deep_link_scheme", cfg.Server.Capabilities)
	}
	if !containsCapability(cfg.Server.Capabilities, "background_push_wake") {
		t.Fatalf("default capabilities = %v, want background_push_wake", cfg.Server.Capabilities)
	}
}

func containsCapability(capabilities []string, expected string) bool {
	for _, capability := range capabilities {
		if capability == expected {
			return true
		}
	}
	return false
}
