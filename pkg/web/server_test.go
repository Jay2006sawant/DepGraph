package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/DepGraph/pkg/graph"
)

func setupTestServer() (*Server, error) {
	g := graph.NewGraph()

	// Add test data
	g.AddNode(&graph.Node{
		ID:    "repo1",
		Label: "test-repo-1",
		Type:  "repository",
	})

	g.AddNode(&graph.Node{
		ID:      "mod1",
		Label:   "github.com/test/mod1",
		Type:    "module",
		Version: "v1.0.0",
	})

	g.AddEdge(&graph.Edge{
		Source:  "repo1",
		Target:  "mod1",
		Version: "v1.0.0",
	})

	return NewServer(g)
}

func TestHandleIndex(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Expected HTML content type, got %s", w.Header().Get("Content-Type"))
	}
}

func TestHandleGraphData(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/graph", nil)
	w := httptest.NewRecorder()

	server.handleGraphData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Nodes []map[string]interface{} `json:"nodes"`
		Links []map[string]interface{} `json:"links"`
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(response.Nodes))
	}

	if len(response.Links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(response.Links))
	}
}

func TestHandleStats(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	server.handleStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var stats struct {
		TotalRepositories int `json:"totalRepositories"`
		TotalModules     int `json:"totalModules"`
	}

	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if stats.TotalRepositories != 1 {
		t.Errorf("Expected 1 repository, got %d", stats.TotalRepositories)
	}

	if stats.TotalModules != 1 {
		t.Errorf("Expected 1 module, got %d", stats.TotalModules)
	}
}

func TestHandleSecurityScan(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/security", nil)
	w := httptest.NewRecorder()

	server.handleSecurityScan(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var results []struct {
		Module    string `json:"module"`
		RiskLevel string `json:"riskLevel"`
	}

	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should have no high-risk modules in test data
	for _, result := range results {
		if result.RiskLevel == "HIGH" {
			t.Errorf("Unexpected high-risk module: %s", result.Module)
		}
	}
}

func TestHandleChains(t *testing.T) {
	server, err := setupTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/chains", nil)
	w := httptest.NewRecorder()

	server.handleChains(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var chains []struct {
		Path   []string `json:"path"`
		Length int      `json:"length"`
	}

	if err := json.NewDecoder(w.Body).Decode(&chains); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(chains) == 0 {
		t.Error("Expected at least one dependency chain")
	}
} 