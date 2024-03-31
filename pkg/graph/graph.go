package graph

import (
	"encoding/json"
	"fmt"
)

// Node represents a module in the dependency graph
type Node struct {
	ID            string                 `json:"id"`
	Label         string                 `json:"label"`
	Type          string                 `json:"type"`           // "repository" or "module"
	Version       string                 `json:"version"`
	Repositories  []string              `json:"repositories"`   // For module nodes
	Dependencies  []string              `json:"dependencies"`   // List of module IDs
	Metadata      map[string]interface{} `json:"metadata"`
}

// Edge represents a dependency relationship
type Edge struct {
	Source    string `json:"source"`
	Target    string `json:"target"`
	Type      string `json:"type"`      // "direct" or "indirect"
	Version   string `json:"version"`
}

// Graph represents the dependency network
type Graph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges []*Edge         `json:"edges"`
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
	}
}

// AddNode adds a node to the graph if it doesn't exist
func (g *Graph) AddNode(node *Node) error {
	if node.ID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}

	if _, exists := g.Nodes[node.ID]; !exists {
		g.Nodes[node.ID] = node
	}
	return nil
}

// AddEdge adds an edge between two nodes
func (g *Graph) AddEdge(edge *Edge) error {
	if edge.Source == "" || edge.Target == "" {
		return fmt.Errorf("edge source and target cannot be empty")
	}

	// Verify nodes exist
	if _, exists := g.Nodes[edge.Source]; !exists {
		return fmt.Errorf("source node %s does not exist", edge.Source)
	}
	if _, exists := g.Nodes[edge.Target]; !exists {
		return fmt.Errorf("target node %s does not exist", edge.Target)
	}

	g.Edges = append(g.Edges, edge)
	return nil
}

// GetDependents returns all nodes that depend on the given node
func (g *Graph) GetDependents(nodeID string) []*Node {
	var dependents []*Node
	for _, edge := range g.Edges {
		if edge.Target == nodeID {
			if node, exists := g.Nodes[edge.Source]; exists {
				dependents = append(dependents, node)
			}
		}
	}
	return dependents
}

// GetDependencies returns all nodes that the given node depends on
func (g *Graph) GetDependencies(nodeID string) []*Node {
	var dependencies []*Node
	for _, edge := range g.Edges {
		if edge.Source == nodeID {
			if node, exists := g.Nodes[edge.Target]; exists {
				dependencies = append(dependencies, node)
			}
		}
	}
	return dependencies
}

// FindCycles detects dependency cycles in the graph
func (g *Graph) FindCycles() [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	path := make(map[string]bool)

	var dfs func(node string, current []string)
	dfs = func(node string, current []string) {
		visited[node] = true
		path[node] = true
		current = append(current, node)

		for _, edge := range g.Edges {
			if edge.Source != node {
				continue
			}

			if path[edge.Target] {
				// Found a cycle
				cycleStart := -1
				for i, n := range current {
					if n == edge.Target {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := append(current[cycleStart:], edge.Target)
					cycles = append(cycles, cycle)
				}
			} else if !visited[edge.Target] {
				dfs(edge.Target, current)
			}
		}

		path[node] = false
	}

	for node := range g.Nodes {
		if !visited[node] {
			dfs(node, nil)
		}
	}

	return cycles
}

// ToJSON converts the graph to JSON format
func (g *Graph) ToJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nodes []*Node `json:"nodes"`
		Edges []*Edge `json:"edges"`
	}{
		Nodes: g.NodesList(),
		Edges: g.Edges,
	})
}

// NodesList returns a slice of all nodes
func (g *Graph) NodesList() []*Node {
	nodes := make([]*Node, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetSharedDependencies returns modules that are dependencies of multiple repositories
func (g *Graph) GetSharedDependencies(minShared int) []*Node {
	depCount := make(map[string]int)
	
	// Count repository dependencies
	for _, node := range g.Nodes {
		if node.Type == "repository" {
			deps := g.GetDependencies(node.ID)
			for _, dep := range deps {
				depCount[dep.ID]++
			}
		}
	}

	// Filter shared dependencies
	var shared []*Node
	for nodeID, count := range depCount {
		if count >= minShared {
			if node, exists := g.Nodes[nodeID]; exists {
				shared = append(shared, node)
			}
		}
	}

	return shared
} 