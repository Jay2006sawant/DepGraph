package web

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"path"

	"github.com/yourusername/DepGraph/pkg/analysis"
	"github.com/yourusername/DepGraph/pkg/graph"
)

//go:embed templates/* static/*
var content embed.FS

// Server represents the web server for dependency visualization
type Server struct {
	analyzer *analysis.DependencyAnalyzer
	graph    *graph.Graph
	tmpl     *template.Template
}

// NewServer creates a new web server instance
func NewServer(g *graph.Graph) (*Server, error) {
	tmpl, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Server{
		analyzer: analysis.NewAnalyzer(g),
		graph:    g,
		tmpl:     tmpl,
	}, nil
}

// Start starts the web server on the specified address
func (s *Server) Start(addr string) error {
	// Static file server
	fs := http.FileServer(http.FS(content))
	http.Handle("/static/", fs)

	// Route handlers
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/graph", s.handleGraphData)
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/api/chains", s.handleChains)
	http.HandleFunc("/api/security", s.handleSecurityScan)

	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
		Stats *analysis.DependencyStats
	}{
		Title: "DepGraph Visualization",
	}

	stats, err := s.analyzer.AnalyzeDependencies()
	if err == nil {
		data.Stats = stats
	}

	s.tmpl.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) handleGraphData(w http.ResponseWriter, r *http.Request) {
	// Convert graph data to D3.js format
	data := struct {
		Nodes []map[string]interface{} `json:"nodes"`
		Links []map[string]interface{} `json:"links"`
	}{
		Nodes: make([]map[string]interface{}, 0, len(s.graph.Nodes)),
		Links: make([]map[string]interface{}, 0, len(s.graph.Edges)),
	}

	for id, node := range s.graph.Nodes {
		data.Nodes = append(data.Nodes, map[string]interface{}{
			"id":    id,
			"label": node.Label,
			"type":  node.Type,
		})
	}

	for _, edge := range s.graph.Edges {
		data.Links = append(data.Links, map[string]interface{}{
			"source":  edge.Source,
			"target":  edge.Target,
			"version": edge.Version,
		})
	}

	json.NewEncoder(w).Encode(data)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.analyzer.AnalyzeDependencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleChains(w http.ResponseWriter, r *http.Request) {
	chains := s.analyzer.FindLongestDependencyChains(5)
	json.NewEncoder(w).Encode(chains)
}

func (s *Server) handleSecurityScan(w http.ResponseWriter, r *http.Request) {
	results := s.analyzer.SimulateSecurityScan()
	json.NewEncoder(w).Encode(results)
} 