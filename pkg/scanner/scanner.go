package scanner

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/yourusername/DepGraph/pkg/github"
)

// Scanner handles repository scanning operations
type Scanner struct {
	client *github.Client
}

// NewScanner creates a new repository scanner
func NewScanner(client *github.Client) *Scanner {
	return &Scanner{
		client: client,
	}
}

// RepoInfo contains repository dependency information
type RepoInfo struct {
	Name         string
	ModuleFile   string
	Dependencies []Dependency
	Error        error
}

// Dependency represents a module dependency
type Dependency struct {
	Module  string
	Version string
}

// ScanOrganization scans all repositories in an organization
func (s *Scanner) ScanOrganization(org string) ([]RepoInfo, error) {
	repos, err := s.client.ListRepositories(org)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %v", err)
	}

	var wg sync.WaitGroup
	results := make([]RepoInfo, len(repos))

	// Process repositories concurrently
	for i, repo := range repos {
		wg.Add(1)
		go func(idx int, repo *github.Repository) {
			defer wg.Done()
			results[idx] = s.scanRepository(org, *repo.Name)
		}(i, repo)
	}

	wg.Wait()
	return results, nil
}

// scanRepository scans a single repository for dependencies
func (s *Scanner) scanRepository(owner, repo string) RepoInfo {
	info := RepoInfo{
		Name: repo,
	}

	// Try to fetch go.mod file
	moduleFile, err := s.client.GetModuleFile(owner, repo, "go.mod")
	if err != nil {
		info.Error = fmt.Errorf("failed to get module file: %v", err)
		return info
	}

	info.ModuleFile = moduleFile
	info.Dependencies = parseDependencies(moduleFile)
	return info
}

// parseDependencies extracts dependencies from go.mod content
func parseDependencies(content string) []Dependency {
	// TODO: Implement proper go.mod parsing in the next commit
	// This is a placeholder that will be enhanced
	return []Dependency{}
} 