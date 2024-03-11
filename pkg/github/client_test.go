package github

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Save current token and restore after test
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", originalToken)

	tests := []struct {
		name      string
		setToken  bool
		wantError bool
	}{
		{
			name:      "valid token",
			setToken:  true,
			wantError: false,
		},
		{
			name:      "missing token",
			setToken:  false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setToken {
				os.Setenv("GITHUB_TOKEN", "dummy-token")
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			client, err := NewClient()
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && client == nil {
				t.Error("NewClient() returned nil client when error not expected")
			}
		})
	}
} 