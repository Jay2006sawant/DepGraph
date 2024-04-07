package analysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/DepGraph/pkg/graph"
)

// DependencyAnalyzer handles dependency analysis operations
type DependencyAnalyzer struct {
	graph *graph.Graph
}

// NewAnalyzer creates a new dependency analyzer
func NewAnalyzer(g *graph.Graph) *DependencyAnalyzer {
	return &DependencyAnalyzer{
		graph: g,
	}
}

// VersionConflict represents a version conflict between repositories
type VersionConflict struct {
	Module      string
	Versions    map[string][]string // version -> repositories using it
	Latest      string
	Recommended string
}

// DependencyStats contains statistical information about dependencies
type DependencyStats struct {
	TotalModules        int
	TotalRepositories   int
	SharedModules       int
	VersionConflicts    int
	AverageDependencies float64
	TopShared          []string
}

// FindVersionConflicts detects modules used with different versions
func (a *DependencyAnalyzer) FindVersionConflicts() []VersionConflict {
	moduleVersions := make(map[string]map[string][]string)

	// Collect all versions of each module
	for _, edge := range a.graph.Edges {
		node := a.graph.Nodes[edge.Source]
		if node.Type != "repository" {
			continue
		}

		module := edge.Target
		version := edge.Version
		repoName := node.Label

		if moduleVersions[module] == nil {
			moduleVersions[module] = make(map[string][]string)
		}
		moduleVersions[module][version] = append(moduleVersions[module][version], repoName)
	}

	// Find conflicts
	var conflicts []VersionConflict
	for module, versions := range moduleVersions {
		if len(versions) > 1 {
			// Sort versions to find latest
			allVersions := make([]string, 0, len(versions))
			for v := range versions {
				allVersions = append(allVersions, v)
			}
			sort.Strings(allVersions)
			latest := allVersions[len(allVersions)-1]

			conflicts = append(conflicts, VersionConflict{
				Module:      module,
				Versions:    versions,
				Latest:      latest,
				Recommended: latest, // In a real implementation, this would use semver comparison
			})
		}
	}

	return conflicts
}

// AnalyzeDependencies performs comprehensive dependency analysis
func (a *DependencyAnalyzer) AnalyzeDependencies() (*DependencyStats, error) {
	stats := &DependencyStats{}

	// Count repositories and modules
	for _, node := range a.graph.Nodes {
		if node.Type == "repository" {
			stats.TotalRepositories++
		} else if node.Type == "module" {
			stats.TotalModules++
		}
	}

	// Find shared modules
	shared := a.graph.GetSharedDependencies(2)
	stats.SharedModules = len(shared)

	// Get top shared modules
	sort.Slice(shared, func(i, j int) bool {
		return len(a.graph.GetDependents(shared[i].ID)) > len(a.graph.GetDependents(shared[j].ID))
	})

	for i := 0; i < min(5, len(shared)); i++ {
		stats.TopShared = append(stats.TopShared, shared[i].Label)
	}

	// Calculate average dependencies per repository
	var totalDeps int
	for _, node := range a.graph.Nodes {
		if node.Type == "repository" {
			deps := a.graph.GetDependencies(node.ID)
			totalDeps += len(deps)
		}
	}
	if stats.TotalRepositories > 0 {
		stats.AverageDependencies = float64(totalDeps) / float64(stats.TotalRepositories)
	}

	// Count version conflicts
	conflicts := a.FindVersionConflicts()
	stats.VersionConflicts = len(conflicts)

	return stats, nil
}

// FindCriticalDependencies identifies critical dependencies
func (a *DependencyAnalyzer) FindCriticalDependencies() []string {
	// Find modules that are used by more than 50% of repositories
	threshold := float64(a.countRepositories()) * 0.5
	var critical []string

	for _, node := range a.graph.Nodes {
		if node.Type != "module" {
			continue
		}

		dependents := a.graph.GetDependents(node.ID)
		if float64(len(dependents)) >= threshold {
			critical = append(critical, node.Label)
		}
	}

	return critical
}

// FindUpdateCandidates identifies modules that should be updated
func (a *DependencyAnalyzer) FindUpdateCandidates() map[string][]string {
	candidates := make(map[string][]string)

	conflicts := a.FindVersionConflicts()
	for _, conflict := range conflicts {
		var outdatedRepos []string
		for version, repos := range conflict.Versions {
			if version != conflict.Recommended {
				outdatedRepos = append(outdatedRepos, repos...)
			}
		}
		if len(outdatedRepos) > 0 {
			candidates[conflict.Module] = outdatedRepos
		}
	}

	return candidates
}

func (a *DependencyAnalyzer) countRepositories() int {
	count := 0
	for _, node := range a.graph.Nodes {
		if node.Type == "repository" {
			count++
		}
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 