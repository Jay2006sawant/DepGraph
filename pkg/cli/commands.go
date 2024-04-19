package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourusername/DepGraph/pkg/analysis"
	"github.com/yourusername/DepGraph/pkg/github"
	"github.com/yourusername/DepGraph/pkg/graph"
)

var (
	outputFormat string
	limit        int
	moduleID     string
)

// NewRootCmd creates the root command for DepGraph CLI
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "depgraph",
		Short: "DepGraph - Dependency Analysis Tool",
		Long: `DepGraph is a powerful tool for analyzing dependencies across multiple repositories.
It helps identify version conflicts, security risks, and impact of dependency updates.`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text/json)")

	// Add subcommands
	rootCmd.AddCommand(
		newAnalyzeCmd(),
		newScanCmd(),
		newImpactCmd(),
		newChainsCmd(),
	)

	return rootCmd
}

func newAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze [org/repo]",
		Short: "Analyze dependencies in repositories",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := github.NewClient(os.Getenv("GITHUB_TOKEN"))
			g := graph.NewGraph()
			
			// TODO: Implement repository scanning and graph building
			analyzer := analysis.NewAnalyzer(g)
			stats, err := analyzer.AnalyzeDependencies()
			if err != nil {
				return err
			}

			return outputResults(stats)
		},
	}
	return cmd
}

func newScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [org/repo]",
		Short: "Scan for security vulnerabilities",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := github.NewClient(os.Getenv("GITHUB_TOKEN"))
			g := graph.NewGraph()
			
			// TODO: Implement repository scanning
			analyzer := analysis.NewAnalyzer(g)
			results := analyzer.SimulateSecurityScan()

			return outputResults(results)
		},
	}
	return cmd
}

func newImpactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact",
		Short: "Analyze impact of module updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if moduleID == "" {
				return fmt.Errorf("module ID is required")
			}

			g := graph.NewGraph()
			analyzer := analysis.NewAnalyzer(g)
			impact, err := analyzer.AnalyzeModuleImpact(moduleID)
			if err != nil {
				return err
			}

			return outputResults(impact)
		},
	}

	cmd.Flags().StringVarP(&moduleID, "module", "m", "", "Module ID to analyze")
	cmd.MarkFlagRequired("module")
	return cmd
}

func newChainsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chains",
		Short: "Find longest dependency chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			g := graph.NewGraph()
			analyzer := analysis.NewAnalyzer(g)
			chains := analyzer.FindLongestDependencyChains(limit)

			return outputResults(chains)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 5, "Maximum number of chains to show")
	return cmd
}

func outputResults(data interface{}) error {
	if outputFormat == "json" {
		return outputJSON(data)
	}
	return outputText(data)
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func outputText(data interface{}) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	switch v := data.(type) {
	case *analysis.DependencyStats:
		fmt.Fprintf(w, "Total Repositories:\t%d\n", v.TotalRepositories)
		fmt.Fprintf(w, "Total Modules:\t%d\n", v.TotalModules)
		fmt.Fprintf(w, "Shared Modules:\t%d\n", v.SharedModules)
		fmt.Fprintf(w, "Version Conflicts:\t%d\n", v.VersionConflicts)
		fmt.Fprintf(w, "Average Dependencies:\t%.2f\n", v.AverageDependencies)
		fmt.Fprintf(w, "\nTop Shared Modules:\n")
		for _, mod := range v.TopShared {
			fmt.Fprintf(w, "  %s\n", mod)
		}

	case []analysis.SecurityScan:
		fmt.Fprintln(w, "MODULE\tVERSION\tRISK LEVEL\tRECOMMENDED FIX")
		for _, scan := range v {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				scan.Module, scan.Version, scan.RiskLevel, scan.RecommendedFix)
		}

	case *analysis.ImpactAnalysis:
		fmt.Fprintf(w, "Module:\t%s\n", v.Module)
		fmt.Fprintf(w, "Impact Score:\t%.2f\n", v.ImpactScore)
		fmt.Fprintf(w, "Breaking Changes:\t%v\n", v.BreakingChanges)
		fmt.Fprintf(w, "Transitive Dependencies:\t%d\n", v.TransitiveDeps)
		fmt.Fprintf(w, "\nAffected Repositories:\n")
		for _, repo := range v.AffectedRepos {
			fmt.Fprintf(w, "  %s\n", repo)
		}

	case []analysis.DependencyChain:
		for i, chain := range v {
			fmt.Fprintf(w, "\nChain %d (Length: %d):\n", i+1, chain.Length)
			fmt.Fprintf(w, "  %s\n", strings.Join(chain.Path, " → "))
			if chain.Circular {
				fmt.Fprintf(w, "  (Circular Dependency Detected)\n")
			}
		}

	default:
		return fmt.Errorf("unsupported data type for text output")
	}

	return nil
} 