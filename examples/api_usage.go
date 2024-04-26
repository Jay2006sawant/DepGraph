package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yourusername/DepGraph/pkg/analysis"
	"github.com/yourusername/DepGraph/pkg/github"
	"github.com/yourusername/DepGraph/pkg/graph"
)

// DependencyReport represents a comprehensive dependency analysis report
type DependencyReport struct {
	Repository string
	Stats     *analysis.DependencyStats
	Conflicts []*analysis.VersionConflict
	Security  []*analysis.SecurityIssue
	Chains    []*analysis.DependencyChain
}

// RepositoryAnalyzer provides high-level dependency analysis functionality
type RepositoryAnalyzer struct {
	client   *github.Client
	graph    *graph.Graph
	analyzer *analysis.DependencyAnalyzer
}

// NewRepositoryAnalyzer creates a new analyzer instance
func NewRepositoryAnalyzer(token string) *RepositoryAnalyzer {
	return &RepositoryAnalyzer{
		client: github.NewClient(token),
		graph:  graph.NewGraph(),
	}
}

// AnalyzeRepository performs a comprehensive analysis of a repository
func (ra *RepositoryAnalyzer) AnalyzeRepository(owner, repo string) (*DependencyReport, error) {
	// Create repository node
	repoID := fmt.Sprintf("%s/%s", owner, repo)
	ra.graph.AddNode(&graph.Node{
		ID:    repoID,
		Label: repoID,
		Type:  "repository",
	})

	// Get dependencies from go.mod
	gomod, err := ra.client.GetGoMod(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get go.mod: %w", err)
	}

	// Add dependencies to graph
	for _, dep := range gomod.Dependencies {
		ra.graph.AddNode(&graph.Node{
			ID:      dep.Path,
			Label:   dep.Path,
			Type:    "module",
			Version: dep.Version,
		})

		ra.graph.AddEdge(&graph.Edge{
			Source:  repoID,
			Target:  dep.Path,
			Version: dep.Version,
		})
	}

	// Initialize analyzer
	ra.analyzer = analysis.NewAnalyzer(ra.graph)

	// Generate report
	report := &DependencyReport{
		Repository: repoID,
	}

	// Get statistics
	stats, err := ra.analyzer.AnalyzeDependencies()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze dependencies: %w", err)
	}
	report.Stats = stats

	// Get version conflicts
	report.Conflicts = ra.analyzer.FindVersionConflicts()

	// Get security issues
	report.Security = ra.analyzer.SimulateSecurityScan()

	// Get dependency chains
	report.Chains = ra.analyzer.FindLongestDependencyChains(5)

	return report, nil
}

// PrintReport prints a formatted dependency report
func PrintReport(report *DependencyReport) {
	fmt.Printf("Dependency Analysis Report for %s\n", report.Repository)
	fmt.Printf("Generated at: %s\n\n", time.Now().Format(time.RFC3339))

	// Print statistics
	fmt.Println("=== Statistics ===")
	fmt.Printf("Total Modules: %d\n", report.Stats.TotalModules)
	fmt.Printf("Shared Modules: %d\n", report.Stats.SharedModules)
	fmt.Printf("Version Conflicts: %d\n", report.Stats.VersionConflicts)
	fmt.Printf("Average Dependencies: %.2f\n\n", report.Stats.AverageDependencies)

	// Print conflicts
	if len(report.Conflicts) > 0 {
		fmt.Println("=== Version Conflicts ===")
		for _, conflict := range report.Conflicts {
			fmt.Printf("Module: %s\n", conflict.Module)
			for version, repos := range conflict.Versions {
				fmt.Printf("  %s: used by %s\n", version, repos)
			}
			fmt.Printf("  Recommended: %s\n\n", conflict.Recommended)
		}
	}

	// Print security issues
	if len(report.Security) > 0 {
		fmt.Println("=== Security Issues ===")
		for _, issue := range report.Security {
			fmt.Printf("Module: %s (Version: %s)\n", issue.Module, issue.Version)
			fmt.Printf("  Risk Level: %s\n", issue.RiskLevel)
			fmt.Printf("  Fix: %s\n\n", issue.RecommendedFix)
		}
	}

	// Print dependency chains
	if len(report.Chains) > 0 {
		fmt.Println("=== Longest Dependency Chains ===")
		for i, chain := range report.Chains {
			fmt.Printf("%d. Length: %d\n", i+1, chain.Length)
			fmt.Printf("   Path: %s\n\n", chain.Path)
		}
	}
}

func main() {
	// Get GitHub token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Create analyzer
	analyzer := NewRepositoryAnalyzer(token)

	// Get repository from command line
	if len(os.Args) != 3 {
		log.Fatal("Usage: api_usage <owner> <repo>")
	}
	owner := os.Args[1]
	repo := os.Args[2]

	// Analyze repository
	report, err := analyzer.AnalyzeRepository(owner, repo)
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Print report
	PrintReport(report)
} 