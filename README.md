# DepGraph - Multi-Repo Dependency Visualizer

DepGraph is a powerful Go-based tool that helps organizations visualize and manage dependencies across multiple repositories. It provides real-time insights into module relationships, outdated packages, and dependency overlaps.

## 🔑 Key Features

- **Cross-Repository Scanning**: Analyze multiple repositories simultaneously
- **GitHub API Integration**: Fetch repository metadata and dependency information
- **Interactive Visualization**: Dynamic dependency graphs with D3.js
- **Update Insights**: Detect shared dependencies and version conflicts
- **Automated Scanning**: Periodic scans to keep dependency information up-to-date

## 🛠️ Tech Stack

- **Backend**: Go
- **API Integration**: GitHub REST API
- **Database**: SQLite/PostgreSQL
- **Visualization**: D3.js
- **Scheduling**: Native Go scheduling

## 📋 Prerequisites

- Go 1.19 or higher
- GitHub Personal Access Token
- SQLite or PostgreSQL
- Node.js (for D3.js visualization)

## 🚀 Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/DepGraph.git
   cd DepGraph
   ```

2. Set up environment variables:
   ```bash
   export GITHUB_TOKEN=your_github_token
   export DB_CONNECTION=your_db_connection_string
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run the application:
   ```bash
   go run cmd/depgraph/main.go
   ```

## 📁 Project Structure

```
DepGraph/
├── cmd/
│   └── depgraph/
│       └── main.go
├── internal/
│   ├── api/
│   ├── database/
│   ├── scanner/
│   └── visualization/
├── pkg/
│   ├── github/
│   ├── parser/
│   └── graph/
└── web/
    ├── static/
    └── templates/
```

## 📄 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 