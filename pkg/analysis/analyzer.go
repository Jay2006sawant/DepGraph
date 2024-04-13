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

// DependencyChain represents a chain of dependencies
type DependencyChain struct {
	Path     []string
	Length   int
	Circular bool
}

// ImpactAnalysis represents the impact of a module change
type ImpactAnalysis struct {
	Module            string
	AffectedRepos    []string
	ImpactScore      float64
	TransitiveDeps   int
	BreakingChanges  bool
}

// SecurityScan represents a simulated security scan result
type SecurityScan struct {
	Module          string
	Version         string
	RiskLevel       string
	AffectedRepos   []string
	RecommendedFix  string
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

// FindLongestDependencyChains identifies the longest dependency chains
func (a *DependencyAnalyzer) FindLongestDependencyChains(limit int) []DependencyChain {
	var chains []DependencyChain
	visited := make(map[string]bool)

	for _, node := range a.graph.Nodes {
		if node.Type != "repository" || visited[node.ID] {
			continue
		}

		chain := a.traverseDependencyChain(node.ID, []string{}, make(map[string]bool))
		if len(chain.Path) > 0 {
			chains = append(chains, chain)
		}
		visited[node.ID] = true
	}

	// Sort by length descending
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].Length > chains[j].Length
	})

	if len(chains) > limit {
		chains = chains[:limit]
	}

	return chains
}

// AnalyzeModuleImpact calculates the potential impact of updating a module
func (a *DependencyAnalyzer) AnalyzeModuleImpact(moduleID string) (*ImpactAnalysis, error) {
	node := a.graph.Nodes[moduleID]
	if node == nil || node.Type != "module" {
		return nil, fmt.Errorf("invalid module ID: %s", moduleID)
	}

	affected := a.graph.GetDependents(moduleID)
	transitive := a.countTransitiveDependents(moduleID)

	// Calculate impact score based on number of affected repos and their importance
	impactScore := float64(len(affected)) * 0.6 + float64(transitive) * 0.4

	affectedRepos := make([]string, 0, len(affected))
	for _, repo := range affected {
		if a.graph.Nodes[repo].Type == "repository" {
			affectedRepos = append(affectedRepos, a.graph.Nodes[repo].Label)
		}
	}

	return &ImpactAnalysis{
		Module:          node.Label,
		AffectedRepos:   affectedRepos,
		ImpactScore:     impactScore,
		TransitiveDeps:  transitive,
		BreakingChanges: impactScore > 0.7, // Simulate breaking change detection
	}, nil
}

// SimulateSecurityScan performs a simulated security vulnerability scan
func (a *DependencyAnalyzer) SimulateSecurityScan() []SecurityScan {
	var results []SecurityScan

	for _, node := range a.graph.Nodes {
		if node.Type != "module" {
			continue
		}

		// Simulate vulnerability detection based on version patterns
		var riskLevel string
		var recommendedFix string

		version := node.Version
		if strings.HasPrefix(version, "v0.") {
			riskLevel = "HIGH"
			recommendedFix = "Upgrade to stable version"
		} else if strings.Contains(version, "alpha") || strings.Contains(version, "beta") {
			riskLevel = "MEDIUM"
			recommendedFix = "Consider upgrading to stable release"
		} else if strings.HasPrefix(version, "v1.0.") {
			riskLevel = "LOW"
			recommendedFix = "No action required"
		}

		if riskLevel != "" {
			affected := a.graph.GetDependents(node.ID)
			affectedRepos := make([]string, 0)
			for _, repo := range affected {
				if a.graph.Nodes[repo].Type == "repository" {
					affectedRepos = append(affectedRepos, a.graph.Nodes[repo].Label)
				}
			}

			results = append(results, SecurityScan{
				Module:         node.Label,
				Version:       version,
				RiskLevel:     riskLevel,
				AffectedRepos: affectedRepos,
				RecommendedFix: recommendedFix,
			})
		}
	}

	return results
}

func (a *DependencyAnalyzer) traverseDependencyChain(nodeID string, path []string, visited map[string]bool) DependencyChain {
	if visited[nodeID] {
		return DependencyChain{
			Path:     append(path, nodeID),
			Length:   len(path),
			Circular: true,
		}
	}

	visited[nodeID] = true
	path = append(path, nodeID)

	maxChain := DependencyChain{
		Path:   path,
		Length: len(path),
	}

	deps := a.graph.GetDependencies(nodeID)
	for _, dep := range deps {
		chain := a.traverseDependencyChain(dep, path, visited)
		if chain.Length > maxChain.Length {
			maxChain = chain
		}
	}

	return maxChain
}

func (a *DependencyAnalyzer) countTransitiveDependents(moduleID string) int {
	visited := make(map[string]bool)
	var count int

	var traverse func(string)
	traverse = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		count++

		for _, dep := range a.graph.GetDependents(id) {
			traverse(dep)
		}
	}

	traverse(moduleID)
	return count - 1 // Subtract 1 to exclude the starting module
} 