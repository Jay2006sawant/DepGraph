package graph

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGraphOperations(t *testing.T) {
	g := NewGraph()

	// Test adding nodes
	repoNode := &Node{
		ID:    "repo1",
		Label: "test-repo",
		Type:  "repository",
	}

	moduleNode := &Node{
		ID:      "mod1",
		Label:   "github.com/test/module",
		Type:    "module",
		Version: "v1.0.0",
	}

	if err := g.AddNode(repoNode); err != nil {
		t.Errorf("Failed to add repository node: %v", err)
	}

	if err := g.AddNode(moduleNode); err != nil {
		t.Errorf("Failed to add module node: %v", err)
	}

	// Test adding edge
	edge := &Edge{
		Source:  "repo1",
		Target:  "mod1",
		Type:    "direct",
		Version: "v1.0.0",
	}

	if err := g.AddEdge(edge); err != nil {
		t.Errorf("Failed to add edge: %v", err)
	}

	// Test getting dependencies
	deps := g.GetDependencies("repo1")
	if len(deps) != 1 || deps[0].ID != "mod1" {
		t.Error("GetDependencies did not return expected dependencies")
	}

	// Test getting dependents
	dependents := g.GetDependents("mod1")
	if len(dependents) != 1 || dependents[0].ID != "repo1" {
		t.Error("GetDependents did not return expected dependents")
	}
}

func TestCycleDetection(t *testing.T) {
	g := NewGraph()

	// Create a cycle: A -> B -> C -> A
	nodes := []string{"A", "B", "C"}
	for _, id := range nodes {
		g.AddNode(&Node{ID: id})
	}

	g.AddEdge(&Edge{Source: "A", Target: "B"})
	g.AddEdge(&Edge{Source: "B", Target: "C"})
	g.AddEdge(&Edge{Source: "C", Target: "A"})

	cycles := g.FindCycles()
	if len(cycles) != 1 {
		t.Errorf("Expected 1 cycle, got %d", len(cycles))
	}

	// The cycle should contain all three nodes
	expectedCycle := []string{"A", "B", "C", "A"}
	if !reflect.DeepEqual(cycles[0], expectedCycle) {
		t.Errorf("Expected cycle %v, got %v", expectedCycle, cycles[0])
	}
}

func TestSharedDependencies(t *testing.T) {
	g := NewGraph()

	// Add repositories
	repos := []string{"repo1", "repo2", "repo3"}
	for _, id := range repos {
		g.AddNode(&Node{
			ID:   id,
			Type: "repository",
		})
	}

	// Add shared module
	sharedMod := &Node{
		ID:   "shared",
		Type: "module",
	}
	g.AddNode(sharedMod)

	// All repos depend on shared module
	for _, repo := range repos {
		g.AddEdge(&Edge{
			Source: repo,
			Target: "shared",
		})
	}

	// Test finding shared dependencies
	shared := g.GetSharedDependencies(2)
	if len(shared) != 1 || shared[0].ID != "shared" {
		t.Error("GetSharedDependencies did not identify shared module")
	}
}

func TestGraphJSON(t *testing.T) {
	g := NewGraph()

	// Add test data
	g.AddNode(&Node{
		ID:    "repo1",
		Label: "test-repo",
		Type:  "repository",
	})
	g.AddNode(&Node{
		ID:      "mod1",
		Label:   "test-module",
		Type:    "module",
		Version: "v1.0.0",
	})
	g.AddEdge(&Edge{
		Source:  "repo1",
		Target:  "mod1",
		Type:    "direct",
		Version: "v1.0.0",
	})

	// Test JSON conversion
	data, err := g.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert graph to JSON: %v", err)
	}

	// Verify JSON structure
	var result struct {
		Nodes []*Node `json:"nodes"`
		Edges []*Edge `json:"edges"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(result.Nodes) != 2 {
		t.Errorf("Expected 2 nodes in JSON, got %d", len(result.Nodes))
	}

	if len(result.Edges) != 1 {
		t.Errorf("Expected 1 edge in JSON, got %d", len(result.Edges))
	}
} 