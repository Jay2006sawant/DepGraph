# DepGraph API Documentation

This document describes the public API packages provided by DepGraph.

## Table of Contents

1. [Graph Package](#graph-package)
2. [Analysis Package](#analysis-package)
3. [GitHub Package](#github-package)
4. [Web Package](#web-package)

## Graph Package

Package `github.com/yourusername/DepGraph/pkg/graph`

### Types

#### Graph

```go
type Graph struct {
    Nodes map[string]*Node
    Edges []*Edge
}
```

The main data structure for representing dependency relationships.

##### Methods

```go
func NewGraph() *Graph
```
Creates a new empty graph.

```go
func (g *Graph) AddNode(node *Node) bool
```
Adds a node to the graph. Returns true if the node was added, false if it already exists.

```go
func (g *Graph) AddEdge(edge *Edge) bool
```
Adds an edge to the graph. Returns true if the edge was added, false if it already exists.

#### Node

```go
type Node struct {
    ID      string
    Label   string
    Type    string
    Version string
}
```

Represents a node in the dependency graph (repository or module).

#### Edge

```go
type Edge struct {
    Source  string
    Target  string
    Version string
}
```

Represents a dependency relationship between nodes.

## Analysis Package

Package `github.com/yourusername/DepGraph/pkg/analysis`

### Types

#### DependencyAnalyzer

```go
type DependencyAnalyzer struct {
    graph *graph.Graph
}
```

Provides dependency analysis functionality.

##### Methods

```go
func NewAnalyzer(g *graph.Graph) *DependencyAnalyzer
```
Creates a new analyzer for the given graph.

```go
func (a *DependencyAnalyzer) AnalyzeDependencies() (*DependencyStats, error)
```
Analyzes dependencies and returns statistics.

```go
func (a *DependencyAnalyzer) FindVersionConflicts() []*VersionConflict
```
Finds modules with conflicting versions across repositories.

```go
func (a *DependencyAnalyzer) SimulateSecurityScan() []*SecurityIssue
```
Performs a security scan of dependencies.

```go
func (a *DependencyAnalyzer) FindLongestDependencyChains(limit int) []*DependencyChain
```
Finds the longest dependency chains in the graph.

#### DependencyStats

```go
type DependencyStats struct {
    TotalRepositories   int
    TotalModules       int
    SharedModules      int
    VersionConflicts   int
    AverageDependencies float64
    TopShared          []string
}
```

Statistics about dependencies in the graph.

#### VersionConflict

```go
type VersionConflict struct {
    Module      string
    Versions    map[string][]string
    Recommended string
}
```

Information about version conflicts for a module.

#### SecurityIssue

```go
type SecurityIssue struct {
    Module          string
    Version         string
    RiskLevel       string
    RecommendedFix  string
}
```

Security vulnerability information for a module.

#### DependencyChain

```go
type DependencyChain struct {
    Path     []string
    Length   int
    Circular bool
}
```

Represents a chain of dependencies.

## GitHub Package

Package `github.com/yourusername/DepGraph/pkg/github`

### Types

#### Client

```go
type Client struct {
    // internal fields
}
```

GitHub API client for fetching repository information.

##### Methods

```go
func NewClient(token string) *Client
```
Creates a new GitHub client with the given access token.

```go
func (c *Client) ListOrgRepositories(org string) ([]string, error)
```
Lists all repositories in an organization.

```go
func (c *Client) GetGoMod(owner, repo string) (*GoMod, error)
```
Gets the go.mod file contents for a repository.

#### GoMod

```go
type GoMod struct {
    Module       string
    Go           string
    Dependencies []*Dependency
}
```

Represents a parsed go.mod file.

#### Dependency

```go
type Dependency struct {
    Path    string
    Version string
}
```

Represents a single dependency in a go.mod file.

## Web Package

Package `github.com/yourusername/DepGraph/pkg/web`

### Types

#### Server

```go
type Server struct {
    // internal fields
}
```

Web server for dependency visualization.

##### Methods

```go
func NewServer(g *graph.Graph) (*Server, error)
```
Creates a new web server instance.

```go
func (s *Server) Start(addr string) error
```
Starts the web server on the specified address.

### HTTP Endpoints

#### GET /

Returns the main visualization page.

#### GET /api/graph

Returns the graph data in D3.js format:
```json
{
    "nodes": [
        {
            "id": "string",
            "label": "string",
            "type": "string"
        }
    ],
    "links": [
        {
            "source": "string",
            "target": "string",
            "version": "string"
        }
    ]
}
```

#### GET /api/stats

Returns dependency statistics:
```json
{
    "totalRepositories": 0,
    "totalModules": 0,
    "sharedModules": 0,
    "versionConflicts": 0,
    "averageDependencies": 0.0,
    "topShared": ["string"]
}
```

#### GET /api/chains

Returns dependency chains:
```json
[
    {
        "path": ["string"],
        "length": 0,
        "circular": false
    }
]
```

#### GET /api/security

Returns security scan results:
```json
[
    {
        "module": "string",
        "version": "string",
        "riskLevel": "string",
        "recommendedFix": "string"
    }
]
``` 