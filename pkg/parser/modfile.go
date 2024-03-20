package parser

import (
	"bufio"
	"fmt"
	"strings"
)

// ModuleInfo represents parsed go.mod file information
type ModuleInfo struct {
	ModuleName    string
	GoVersion     string
	Dependencies  []Dependency
	Replacements  []Replacement
}

// Dependency represents a single module dependency
type Dependency struct {
	Path    string
	Version string
	Indirect bool
}

// Replacement represents a module replacement
type Replacement struct {
	Old ModuleVersion
	New ModuleVersion
}

// ModuleVersion represents a module path and its version
type ModuleVersion struct {
	Path    string
	Version string
}

// ParseGoMod parses a go.mod file content and returns structured information
func ParseGoMod(content string) (*ModuleInfo, error) {
	info := &ModuleInfo{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	
	var inRequireBlock bool
	var inReplaceBlock bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "module"):
			info.ModuleName = strings.TrimSpace(strings.TrimPrefix(line, "module"))
		
		case strings.HasPrefix(line, "go"):
			info.GoVersion = strings.TrimSpace(strings.TrimPrefix(line, "go"))
		
		case line == "require (":
			inRequireBlock = true
			continue
		
		case line == "replace (":
			inReplaceBlock = true
			continue
		
		case line == ")":
			inRequireBlock = false
			inReplaceBlock = false
			continue
		
		case strings.HasPrefix(line, "require"):
			dep := parseDependencyLine(strings.TrimPrefix(line, "require"))
			if dep != nil {
				info.Dependencies = append(info.Dependencies, *dep)
			}
		
		case strings.HasPrefix(line, "replace"):
			repl := parseReplacementLine(strings.TrimPrefix(line, "replace"))
			if repl != nil {
				info.Replacements = append(info.Replacements, *repl)
			}
		
		default:
			if inRequireBlock {
				dep := parseDependencyLine(line)
				if dep != nil {
					info.Dependencies = append(info.Dependencies, *dep)
				}
			} else if inReplaceBlock {
				repl := parseReplacementLine(line)
				if repl != nil {
					info.Replacements = append(info.Replacements, *repl)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning go.mod: %v", err)
	}

	return info, nil
}

// parseDependencyLine parses a single dependency line
func parseDependencyLine(line string) *Dependency {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	dep := &Dependency{
		Path:    parts[0],
		Version: parts[1],
	}

	// Check if dependency is indirect
	if len(parts) > 2 && parts[len(parts)-1] == "// indirect" {
		dep.Indirect = true
	}

	return dep
}

// parseReplacementLine parses a single replacement line
func parseReplacementLine(line string) *Replacement {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil
	}

	return &Replacement{
		Old: ModuleVersion{
			Path:    parts[0],
			Version: parts[1],
		},
		New: ModuleVersion{
			Path:    parts[2],
			Version: parts[3],
		},
	}
} 