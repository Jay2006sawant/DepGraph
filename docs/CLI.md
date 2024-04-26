# DepGraph CLI Guide

DepGraph provides a powerful command-line interface for analyzing dependencies across multiple repositories.

## Installation

```bash
# Install from source
git clone https://github.com/yourusername/DepGraph
cd DepGraph
go install

# Or using go get
go get -u github.com/yourusername/DepGraph
```

## Environment Setup

Set your GitHub token:
```bash
export GITHUB_TOKEN=your_token_here
```

## Command Overview

### Global Flags

```bash
--config string     Configuration file (default ".depgraph.yaml")
--output string     Output format: text|json (default "text")
--verbose          Enable verbose logging
--no-color        Disable colored output
```

### Available Commands

1. `analyze`: Analyze repository dependencies
2. `scan`: Perform security scanning
3. `impact`: Analyze dependency impact
4. `chains`: Find dependency chains
5. `serve`: Start web interface
6. `version`: Show version information

## Command Details

### analyze

Analyze dependencies in one or more repositories.

```bash
# Analyze a single repository
depgraph analyze owner/repo

# Analyze multiple repositories
depgraph analyze owner/repo1 owner/repo2

# Analyze an entire organization
depgraph analyze --org your-org

# Export as JSON
depgraph analyze owner/repo --output json > analysis.json

# Additional options
depgraph analyze [flags]
  --depth int          Maximum dependency depth (default 3)
  --include-dev        Include development dependencies
  --exclude string     Repositories to exclude (glob pattern)
  --org string         GitHub organization to analyze
  --output string      Output format: text|json (default "text")
```

Example output:
```
Repository: owner/repo
Total Dependencies: 15
Direct Dependencies: 8
Shared Dependencies: 3
Version Conflicts: 2

Dependency Tree:
└── github.com/pkg/errors v0.9.1
    ├── github.com/stretchr/testify v1.8.0
    │   └── github.com/davecgh/go-spew v1.1.1
    └── golang.org/x/sys v0.6.0
```

### scan

Perform security scanning on dependencies.

```bash
# Scan a repository
depgraph scan owner/repo

# Scan with risk level filter
depgraph scan owner/repo --risk-level high

# Additional options
depgraph scan [flags]
  --risk-level string   Minimum risk level: low|medium|high (default "low")
  --fix                Suggest fixes for vulnerabilities
  --notify             Send notifications for issues
```

Example output:
```
Security Scan Results for owner/repo

HIGH Risk:
  - github.com/vulnerable/pkg v1.2.3
    Issue: Remote code execution vulnerability
    Fix: Upgrade to v1.2.4

MEDIUM Risk:
  - github.com/another/pkg v0.1.0
    Issue: Potential data leak
    Fix: Upgrade to v0.2.0
```

### impact

Analyze the impact of updating a dependency.

```bash
# Analyze impact of updating a module
depgraph impact --module github.com/example/module

# Show affected repositories
depgraph impact --module github.com/example/module --show-affected

# Additional options
depgraph impact [flags]
  --module string      Module to analyze
  --version string     Target version
  --show-affected     Show affected repositories
  --breaking          Check for breaking changes
```

Example output:
```
Impact Analysis for github.com/example/module

Current Versions:
  v1.0.0: 3 repositories
  v1.1.0: 2 repositories

Recommended Version: v1.1.0

Affected Repositories:
  - owner/repo1 (v1.0.0)
  - owner/repo2 (v1.0.0)
  - owner/repo3 (v1.0.0)
```

### chains

Find dependency chains in repositories.

```bash
# Find longest chains
depgraph chains owner/repo

# Limit chain length
depgraph chains owner/repo --max-length 5

# Additional options
depgraph chains [flags]
  --max-length int    Maximum chain length (default 10)
  --circular         Show only circular dependencies
  --limit int        Maximum number of chains to show (default 5)
```

Example output:
```
Dependency Chains for owner/repo

1. Length: 4
   repo → mod1 → mod2 → mod3

2. Length: 3 (Circular)
   repo → mod4 → mod5 → mod4

3. Length: 3
   repo → mod6 → mod7 → mod8
```

### serve

Start the web interface.

```bash
# Start server on default port
depgraph serve

# Specify port
depgraph serve --port 8080

# Additional options
depgraph serve [flags]
  --port int          Port number (default 8080)
  --host string       Host address (default "localhost")
  --tls              Enable HTTPS
  --cert string      TLS certificate file
  --key string       TLS key file
```

## Configuration File

The CLI can be configured using a `.depgraph.yaml` file:

```yaml
github:
  token: ${GITHUB_TOKEN}
  organizations:
    - your-org
  exclude:
    - archived-repo
    - test-*

analysis:
  scan_depth: 3
  include_dev_deps: false
  version_check: true
  security_scan: true

output:
  format: text
  color: true
  verbose: false
```

## Environment Variables

- `GITHUB_TOKEN`: GitHub API token
- `DEPGRAPH_CONFIG`: Configuration file path
- `DEPGRAPH_OUTPUT`: Output format (text/json)
- `DEPGRAPH_NO_COLOR`: Disable colored output
- `DEPGRAPH_VERBOSE`: Enable verbose logging

## Exit Codes

- 0: Success
- 1: General error
- 2: Configuration error
- 3: GitHub API error
- 4: Analysis error
- 5: Security issue found

## Examples

1. Analyze and export dependencies:
```bash
depgraph analyze owner/repo --output json | jq .
```

2. Find high-risk security issues:
```bash
depgraph scan owner/repo --risk-level high --output json > security.json
```

3. Check impact of updating a module:
```bash
depgraph impact --module github.com/pkg/errors --version v0.9.1
```

4. Find circular dependencies:
```bash
depgraph chains owner/repo --circular --output json
```

5. Start web interface with custom port:
```bash
depgraph serve --port 3000
```

## Tips

1. Use `--output json` for machine-readable output
2. Set up aliases for common commands
3. Use configuration files for consistent settings
4. Enable verbose logging for troubleshooting
5. Use environment variables in CI/CD pipelines 