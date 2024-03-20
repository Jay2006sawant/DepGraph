package parser

import (
	"reflect"
	"testing"
)

func TestParseGoMod(t *testing.T) {
	content := `module github.com/yourusername/DepGraph

go 1.19

require (
	github.com/google/go-github/v45 v45.2.0
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/oauth2 v0.7.0
)

replace (
	github.com/old/pkg v1.0.0 => github.com/new/pkg v2.0.0
)`

	info, err := ParseGoMod(content)
	if err != nil {
		t.Fatalf("ParseGoMod() error = %v", err)
	}

	// Test module name
	if info.ModuleName != "github.com/yourusername/DepGraph" {
		t.Errorf("Expected module name 'github.com/yourusername/DepGraph', got %s", info.ModuleName)
	}

	// Test Go version
	if info.GoVersion != "1.19" {
		t.Errorf("Expected Go version '1.19', got %s", info.GoVersion)
	}

	// Test dependencies
	expectedDeps := []Dependency{
		{Path: "github.com/google/go-github/v45", Version: "v45.2.0", Indirect: false},
		{Path: "github.com/joho/godotenv", Version: "v1.5.1", Indirect: true},
		{Path: "golang.org/x/oauth2", Version: "v0.7.0", Indirect: false},
	}

	if !reflect.DeepEqual(info.Dependencies, expectedDeps) {
		t.Errorf("Dependencies don't match expected values")
	}

	// Test replacements
	expectedRepl := []Replacement{
		{
			Old: ModuleVersion{Path: "github.com/old/pkg", Version: "v1.0.0"},
			New: ModuleVersion{Path: "github.com/new/pkg", Version: "v2.0.0"},
		},
	}

	if !reflect.DeepEqual(info.Replacements, expectedRepl) {
		t.Errorf("Replacements don't match expected values")
	}
}

func TestParseDependencyLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *Dependency
	}{
		{
			name: "normal dependency",
			line: "github.com/pkg/errors v0.9.1",
			expected: &Dependency{
				Path:     "github.com/pkg/errors",
				Version:  "v0.9.1",
				Indirect: false,
			},
		},
		{
			name: "indirect dependency",
			line: "github.com/pkg/errors v0.9.1 // indirect",
			expected: &Dependency{
				Path:     "github.com/pkg/errors",
				Version:  "v0.9.1",
				Indirect: true,
			},
		},
		{
			name:     "invalid line",
			line:     "invalid",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDependencyLine(tt.line)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseDependencyLine() = %v, want %v", got, tt.expected)
			}
		})
	}
} 