package scanner

import (
	"testing"

	"github.com/yourusername/DepGraph/pkg/github"
)

type mockGitHubClient struct {
	repos     []*github.Repository
	moduleFile string
	err       error
}

func (m *mockGitHubClient) ListRepositories(org string) ([]*github.Repository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.repos, nil
}

func (m *mockGitHubClient) GetModuleFile(owner, repo, path string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.moduleFile, nil
}

func TestScanRepository(t *testing.T) {
	mockClient := &mockGitHubClient{
		moduleFile: `module example.com/test
go 1.19
require (
	github.com/pkg/errors v0.9.1
)`,
	}

	scanner := NewScanner(mockClient)
	info := scanner.scanRepository("testorg", "testrepo")

	if info.Name != "testrepo" {
		t.Errorf("Expected repo name 'testrepo', got %s", info.Name)
	}

	if info.Error != nil {
		t.Errorf("Expected no error, got %v", info.Error)
	}

	if info.ModuleFile == "" {
		t.Error("Expected module file content to be non-empty")
	}
} 