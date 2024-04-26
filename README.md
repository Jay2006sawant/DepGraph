# DepGraph - Multi-Repository Dependency Analyzer

DepGraph is a powerful tool for analyzing and visualizing dependencies across multiple Go repositories. It helps teams understand dependency relationships, identify potential issues, and make informed decisions about dependency updates.

## Features

- **Multi-Repository Analysis**: Scan and analyze dependencies across multiple GitHub repositories
- **Interactive Visualization**: Web-based visualization with multiple layout options (force, tree, radial)
- **Dependency Analysis**:
  - Version conflict detection
  - Security vulnerability scanning
  - Impact analysis for updates
  - Dependency chain analysis
- **CLI Interface**: Powerful command-line interface for automation and scripting
- **Web Interface**: Interactive web UI for exploring dependency relationships

## Installation

```bash
# Install from source
git clone https://github.com/yourusername/DepGraph
cd DepGraph
go install

# Or using go get
go get -u github.com/yourusername/DepGraph
```

## Quick Start

1. Set up your GitHub token:
```bash
export GITHUB_TOKEN=your_token_here
```

2. Analyze a repository:
```bash
depgraph analyze owner/repo
```

3. Start the web interface:
```bash
depgraph serve --port 8080
```

## CLI Usage

### Analyzing Dependencies

```bash
# Analyze a single repository
depgraph analyze owner/repo

# Analyze multiple repositories
depgraph analyze owner/repo1 owner/repo2

# Export results as JSON
depgraph analyze owner/repo --output json > analysis.json
```

### Security Scanning

```bash
# Scan for security vulnerabilities
depgraph scan owner/repo

# Filter by risk level
depgraph scan owner/repo --risk-level high
```

### Impact Analysis

```bash
# Analyze impact of updating a module
depgraph impact --module github.com/example/module

# Show affected repositories
depgraph impact --module github.com/example/module --show-affected
```

### Dependency Chains

```bash
# Find longest dependency chains
depgraph chains owner/repo

# Limit the number of chains
depgraph chains owner/repo --limit 5
```

## Web Interface

The web interface provides an interactive visualization of your dependency graph:

1. Start the server:
```bash
depgraph serve
```

2. Open http://localhost:8080 in your browser

### Features:
- Interactive graph visualization
- Multiple layout options
- Dependency statistics
- Security scan results
- Impact analysis
- Node details on click

## Configuration

Create a `.depgraph.yaml` file in your project root:

```yaml
github:
  token: ${GITHUB_TOKEN}
  organizations:
    - your-org
  exclude:
    - archived-repo
    - test-repo

analysis:
  scan_depth: 3
  include_dev_deps: false
  version_check: true
  security_scan: true

visualization:
  default_layout: force
  theme: light
  node_size: 8
  edge_width: 1.5
```

## API Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/DepGraph/pkg/analysis"
    "github.com/yourusername/DepGraph/pkg/github"
    "github.com/yourusername/DepGraph/pkg/graph"
)

func main() {
    // Initialize GitHub client
    client := github.NewClient(os.Getenv("GITHUB_TOKEN"))

    // Create dependency graph
    g := graph.NewGraph()

    // Add repositories to scan
    scanner := scanner.NewScanner(client, g)
    err := scanner.ScanRepository("owner/repo")
    if err != nil {
        panic(err)
    }

    // Analyze dependencies
    analyzer := analysis.NewAnalyzer(g)
    
    // Get dependency statistics
    stats, _ := analyzer.AnalyzeDependencies()
    fmt.Printf("Total Repositories: %d\n", stats.TotalRepositories)
    
    // Find version conflicts
    conflicts := analyzer.FindVersionConflicts()
    for _, conflict := range conflicts {
        fmt.Printf("Module %s has conflicting versions\n", conflict.Module)
    }
    
    // Perform security scan
    results := analyzer.SimulateSecurityScan()
    for _, result := range results {
        fmt.Printf("Module %s: Risk Level %s\n", result.Module, result.RiskLevel)
    }
}
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 