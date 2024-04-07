package analysis

import (
	"testing"

	"github.com/yourusername/DepGraph/pkg/graph"
)

func setupTestGraph() *graph.Graph {
	g := graph.NewGraph()

	// Add repositories
	repos := []struct {
		id    string
		label string
	}{
		{"repo1", "test-repo-1"},
		{"repo2", "test-repo-2"},
		{"repo3", "test-repo-3"},
	}

	for _, r := range repos {
		g.AddNode(&graph.Node{
			ID:    r.id,
			Label: r.label,
			Type:  "repository",
		})
	}

	// Add modules with different versions
	modules := []struct {
		id      string
		label   string
		version string
	}{
		{"mod1", "github.com/test/mod1", "v1.0.0"},
		{"mod2", "github.com/test/mod2", "v1.0.0"},
		{"mod2-new", "github.com/test/mod2", "v2.0.0"},
	}

	for _, m := range modules {
		g.AddNode(&graph.Node{
			ID:      m.id,
			Label:   m.label,
			Type:    "module",
			Version: m.version,
		})
	}

	// Add dependencies
	edges := []struct {
		source  string
		target  string
		version string
	}{
		{"repo1", "mod1", "v1.0.0"},
		{"repo2", "mod1", "v1.0.0"},
		{"repo1", "mod2", "v1.0.0"},
		{"repo2", "mod2-new", "v2.0.0"},
		{"repo3", "mod2", "v1.0.0"},
	}

	for _, e := range edges {
		g.AddEdge(&graph.Edge{
			Source:  e.source,
			Target:  e.target,
			Version: e.version,
		})
	}

	return g
}

func TestVersionConflictDetection(t *testing.T) {
	g := setupTestGraph()
	analyzer := NewAnalyzer(g)

	conflicts := analyzer.FindVersionConflicts()

	// Should find one conflict for mod2 (v1.0.0 vs v2.0.0)
	if len(conflicts) != 1 {
		t.Errorf("Expected 1 version conflict, got %d", len(conflicts))
	}

	if len(conflicts) > 0 {
		conflict := conflicts[0]
		if len(conflict.Versions) != 2 {
			t.Errorf("Expected 2 different versions, got %d", len(conflict.Versions))
		}
	}
}

func TestDependencyAnalysis(t *testing.T) {
	g := setupTestGraph()
	analyzer := NewAnalyzer(g)

	stats, err := analyzer.AnalyzeDependencies()
	if err != nil {
		t.Fatalf("AnalyzeDependencies failed: %v", err)
	}

	// Verify statistics
	if stats.TotalRepositories != 3 {
		t.Errorf("Expected 3 repositories, got %d", stats.TotalRepositories)
	}

	if stats.TotalModules != 3 {
		t.Errorf("Expected 3 modules, got %d", stats.TotalModules)
	}

	if stats.VersionConflicts != 1 {
		t.Errorf("Expected 1 version conflict, got %d", stats.VersionConflicts)
	}
}

func TestCriticalDependencies(t *testing.T) {
	g := setupTestGraph()
	analyzer := NewAnalyzer(g)

	critical := analyzer.FindCriticalDependencies()

	// mod1 is used by 2/3 repositories, should be critical
	found := false
	for _, dep := range critical {
		if dep == "github.com/test/mod1" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find mod1 as critical dependency")
	}
}

func TestUpdateCandidates(t *testing.T) {
	g := setupTestGraph()
	analyzer := NewAnalyzer(g)

	candidates := analyzer.FindUpdateCandidates()

	// Should find repos using older version of mod2
	if repos, exists := candidates["mod2"]; !exists {
		t.Error("Expected to find update candidates for mod2")
	} else if len(repos) != 2 {
		t.Errorf("Expected 2 repositories needing updates, got %d", len(repos))
	}
} 