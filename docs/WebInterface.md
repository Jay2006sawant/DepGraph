# DepGraph Web Interface Guide

The DepGraph web interface provides an interactive visualization of your dependency graph, along with powerful analysis tools and insights.

## Getting Started

1. Start the web server:
```bash
depgraph serve --port 8080
```

2. Open your browser and navigate to `http://localhost:8080`

## Interface Overview

### Main Components

1. **Graph Visualization**
   - Central area showing the dependency graph
   - Nodes represent repositories and modules
   - Edges represent dependency relationships
   - Different colors indicate node types:
     - Blue: Repositories
     - Green: Modules
     - Red: Modules with conflicts

2. **Control Panel**
   - Layout selector (Force, Tree, Radial)
   - Show/Hide statistics
   - Security scan toggle
   - Zoom controls

3. **Sidebar**
   - Statistics panel
   - Node details
   - Security scan results

## Interacting with the Graph

### Navigation

- **Pan**: Click and drag in empty space
- **Zoom**: Mouse wheel or pinch gesture
- **Center**: Double-click in empty space

### Node Interaction

- **Select**: Click on a node
- **Move**: Drag a node
- **Details**: Click a node to view details in sidebar
- **Highlight**: Hover over a node to highlight its connections

### Layout Options

1. **Force Layout**
   - Default layout
   - Dynamic, force-directed positioning
   - Nodes can be dragged to adjust layout
   - Best for exploring relationships

2. **Tree Layout**
   - Hierarchical visualization
   - Shows dependency chains clearly
   - Fixed node positions
   - Best for understanding dependency flow

3. **Radial Layout**
   - Circular arrangement
   - Shows module relationships around central nodes
   - Good for identifying central dependencies

## Analysis Features

### Dependency Statistics

Click "Show Stats" to view:
- Total repositories and modules
- Shared module count
- Version conflicts
- Average dependencies per repository
- Most shared modules

### Security Scanning

Click "Security Scan" to:
- View security vulnerabilities
- See risk levels for each issue
- Get recommended fixes
- Filter by severity

### Node Details

Click any node to see:
1. For Repositories:
   - Repository name
   - Total dependencies
   - Direct dependencies
   - Dependency chains

2. For Modules:
   - Module path
   - Current version
   - Used by (repositories)
   - Version conflicts
   - Security status

## Tips and Tricks

1. **Finding Dependencies**
   - Use the force layout for general exploration
   - Switch to tree layout to trace dependency chains
   - Use radial layout to identify central modules

2. **Analyzing Conflicts**
   - Red nodes indicate version conflicts
   - Click to see affected repositories
   - View recommended versions

3. **Security Analysis**
   - Regular scans recommended
   - Check high-risk modules first
   - Review recommended fixes
   - Monitor affected repositories

4. **Performance Tips**
   - Limit visible nodes for large graphs
   - Use filters to focus on specific patterns
   - Switch layouts for different perspectives

## Keyboard Shortcuts

- `Space`: Reset zoom and center
- `+/-`: Zoom in/out
- `F`: Switch to force layout
- `T`: Switch to tree layout
- `R`: Switch to radial layout
- `S`: Toggle statistics panel
- `H`: Toggle help overlay

## Customization

The web interface can be customized through the `.depgraph.yaml` configuration:

```yaml
visualization:
  default_layout: force
  theme: light
  graph:
    node_size: 8
    edge_width: 1.5
    show_labels: true
    label_size: 12
  colors:
    repository: "#3498db"
    module: "#2ecc71"
    conflict: "#e74c3c"
    background: "#ffffff"
```

## Troubleshooting

1. **Graph Not Loading**
   - Check server status
   - Verify browser console for errors
   - Ensure WebSocket connection

2. **Performance Issues**
   - Reduce visible nodes
   - Use simpler layouts
   - Clear browser cache

3. **Layout Problems**
   - Reset zoom and center
   - Switch layouts
   - Refresh page

4. **Data Not Updating**
   - Check server connection
   - Verify data source
   - Refresh browser

## Best Practices

1. **Regular Updates**
   - Scan repositories frequently
   - Monitor security issues
   - Update dependencies promptly

2. **Organization**
   - Group related repositories
   - Tag important modules
   - Document known issues

3. **Monitoring**
   - Check statistics regularly
   - Review security scans
   - Track version conflicts

4. **Maintenance**
   - Update dependencies
   - Resolve conflicts
   - Fix security issues 