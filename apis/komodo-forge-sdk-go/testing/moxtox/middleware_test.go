package moxtox

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestInitMoxtoxMiddleware(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "moxtox_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a minimal config file
	configContent := `
enableMoxtox: true
allowedEnvironments:
  - test
performanceMode: default
mappings: {}
ignoredRoutes: []
`
	configPath := filepath.Join(tempDir, "moxtox_config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Initialize middleware
	middleware := InitMoxtoxMiddleware("test", tempDir)

	// Create a test handler
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("real response"))
	}))

	// Test with a request that doesn't match any mocks
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should get 418 (no mocks found) since no mappings are defined
	if w.Code != http.StatusTeapot {
		t.Errorf("Expected status 418, got %d", w.Code)
	}
}