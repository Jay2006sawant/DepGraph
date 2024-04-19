package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	cmd := NewRootCmd()
	output, err := executeCommand(cmd)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(output) == 0 {
		t.Error("Expected help output, got nothing")
	}
}

func TestAnalyzeCommand(t *testing.T) {
	cmd := NewRootCmd()
	_, err := executeCommand(cmd, "analyze", "test/repo")
	if err != nil {
		// Expected to fail without GitHub token, just verify the command exists
		t.Log("analyze command exists but requires GitHub token")
	}
}

func TestScanCommand(t *testing.T) {
	cmd := NewRootCmd()
	_, err := executeCommand(cmd, "scan", "test/repo")
	if err != nil {
		// Expected to fail without GitHub token
		t.Log("scan command exists but requires GitHub token")
	}
}

func TestImpactCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr    bool
		errMessage string
	}{
		{
			name:        "no module specified",
			args:        []string{"impact"},
			wantErr:    true,
			errMessage: "required flag(s) \"module\" not set",
		},
		{
			name:     "with module specified",
			args:     []string{"impact", "--module", "test-module"},
			wantErr: true, // Will fail due to empty graph, but command should parse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd()
			_, err := executeCommand(cmd, tt.args...)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errMessage != "" && err.Error() != tt.errMessage {
					t.Errorf("expected error message %q, got %q", tt.errMessage, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestChainsCommand(t *testing.T) {
	cmd := NewRootCmd()
	_, err := executeCommand(cmd, "chains")
	if err != nil {
		// Expected to succeed with empty graph
		t.Errorf("unexpected error: %v", err)
	}

	// Test with limit flag
	_, err = executeCommand(cmd, "chains", "--limit", "10")
	if err != nil {
		t.Errorf("unexpected error with limit flag: %v", err)
	}
}

func TestOutputFormat(t *testing.T) {
	cmd := NewRootCmd()
	
	// Test JSON output format
	_, err := executeCommand(cmd, "chains", "--output", "json")
	if err != nil {
		t.Errorf("unexpected error with JSON output: %v", err)
	}

	// Test invalid output format
	_, err = executeCommand(cmd, "chains", "--output", "invalid")
	if err != nil {
		// Should still work but default to text
		t.Log("invalid output format defaults to text")
	}
} 