package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/DepGraph/pkg/analysis"
	"github.com/yourusername/DepGraph/pkg/github"
	"github.com/yourusername/DepGraph/pkg/graph"
	"github.com/yourusername/DepGraph/pkg/web"
)

func main() {
	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Initialize GitHub client
	client := github.NewClient(token)

	// Create dependency graph
	g := graph.NewGraph()

	// Get organization name from command line
	if len(os.Args) < 2 {
		log.Fatal("Please provide an organization name")
	}
	org := os.Args[1]

	// List repositories in the organization
	repos, err := client.ListOrgRepositories(org)
	if err != nil {
		log.Fatalf("Failed to list repositories: %v", err)
	}

	// Analyze each repository
	for _, repo := range repos {
		fmt.Printf("Analyzing %s/%s...\n", org, repo)

		// Add repository node
		g.AddNode(&graph.Node{
			ID:    repo,
			Label: fmt.Sprintf("%s/%s", org, repo),
			Type:  "repository",
		})

		// Get go.mod file
		gomod, err := client.GetGoMod(org, repo)
		if err != nil {
			fmt.Printf("Warning: Failed to get go.mod for %s: %v\n", repo, err)
			continue
		}

		// Add dependencies to graph
		for _, dep := range gomod.Dependencies {
			g.AddNode(&graph.Node{
				ID:      dep.Path,
				Label:   dep.Path,
				Type:    "module",
				Version: dep.Version,
			})

			g.AddEdge(&graph.Edge{
				Source:  repo,
				Target:  dep.Path,
				Version: dep.Version,
			})
		}
	}

	// Create analyzer
	analyzer := analysis.NewAnalyzer(g)

	// Get dependency statistics
	stats, err := analyzer.AnalyzeDependencies()
	if err != nil {
		log.Fatalf("Failed to analyze dependencies: %v", err)
	}

	// Print statistics
	fmt.Printf("\nDependency Analysis Results:\n")
	fmt.Printf("Total Repositories: %d\n", stats.TotalRepositories)
	fmt.Printf("Total Modules: %d\n", stats.TotalModules)
	fmt.Printf("Shared Modules: %d\n", stats.SharedModules)
	fmt.Printf("Version Conflicts: %d\n", stats.VersionConflicts)
	fmt.Printf("Average Dependencies: %.2f\n", stats.AverageDependencies)

	// Find version conflicts
	conflicts := analyzer.FindVersionConflicts()
	if len(conflicts) > 0 {
		fmt.Printf("\nVersion Conflicts:\n")
		for _, conflict := range conflicts {
			fmt.Printf("Module: %s\n", conflict.Module)
			fmt.Printf("  Current Versions:\n")
			for version, repos := range conflict.Versions {
				fmt.Printf("    %s: %s\n", version, repos)
			}
			fmt.Printf("  Recommended Version: %s\n", conflict.Recommended)
		}
	}

	// Perform security scan
	fmt.Printf("\nPerforming Security Scan...\n")
	results := analyzer.SimulateSecurityScan()
	if len(results) > 0 {
		fmt.Printf("\nSecurity Issues:\n")
		for _, result := range results {
			fmt.Printf("Module: %s (Version: %s)\n", result.Module, result.Version)
			fmt.Printf("  Risk Level: %s\n", result.RiskLevel)
			fmt.Printf("  Fix: %s\n", result.RecommendedFix)
		}
	}

	// Start web visualization
	fmt.Printf("\nStarting Web Interface...\n")
	server, err := web.NewServer(g)
	if err != nil {
		log.Fatalf("Failed to create web server: %v", err)
	}

	fmt.Printf("Open http://localhost:8080 in your browser\n")
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
} 