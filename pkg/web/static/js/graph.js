class DependencyGraph {
    constructor(containerId) {
        this.container = d3.select(containerId);
        this.width = this.container.node().getBoundingClientRect().width;
        this.height = this.container.node().getBoundingClientRect().height;
        this.currentLayout = 'force';
        
        this.svg = this.container.append('svg')
            .attr('width', this.width)
            .attr('height', this.height);
            
        this.g = this.svg.append('g');
        
        // Add zoom behavior
        this.zoom = d3.zoom()
            .scaleExtent([0.1, 4])
            .on('zoom', (event) => {
                this.g.attr('transform', event.transform);
            });
            
        this.svg.call(this.zoom);
        
        // Initialize forces
        this.initializeForces();
    }
    
    initializeForces() {
        this.simulation = d3.forceSimulation()
            .force('link', d3.forceLink().id(d => d.id))
            .force('charge', d3.forceManyBody().strength(-100))
            .force('center', d3.forceCenter(this.width / 2, this.height / 2))
            .force('collision', d3.forceCollide().radius(30));
    }
    
    async loadData() {
        try {
            const response = await fetch('/api/graph');
            const data = await response.json();
            this.data = data;
            this.render();
        } catch (error) {
            console.error('Error loading graph data:', error);
        }
    }
    
    render() {
        // Clear previous render
        this.g.selectAll('*').remove();
        
        // Create arrow marker for directed edges
        this.svg.append('defs').append('marker')
            .attr('id', 'arrowhead')
            .attr('viewBox', '-0 -5 10 10')
            .attr('refX', 20)
            .attr('refY', 0)
            .attr('orient', 'auto')
            .attr('markerWidth', 6)
            .attr('markerHeight', 6)
            .append('path')
            .attr('d', 'M 0,-5 L 10,0 L 0,5')
            .attr('fill', '#999');
            
        // Create links
        const links = this.g.append('g')
            .selectAll('line')
            .data(this.data.links)
            .enter()
            .append('line')
            .attr('class', 'link')
            .attr('marker-end', 'url(#arrowhead)');
            
        // Create nodes
        const nodes = this.g.append('g')
            .selectAll('.node')
            .data(this.data.nodes)
            .enter()
            .append('g')
            .attr('class', 'node')
            .call(d3.drag()
                .on('start', this.dragStarted.bind(this))
                .on('drag', this.dragged.bind(this))
                .on('end', this.dragEnded.bind(this)));
                
        // Add circles to nodes
        nodes.append('circle')
            .attr('r', 8)
            .attr('fill', d => this.getNodeColor(d.type));
            
        // Add labels to nodes
        nodes.append('text')
            .attr('dx', 12)
            .attr('dy', '.35em')
            .text(d => d.label);
            
        // Update simulation
        this.simulation
            .nodes(this.data.nodes)
            .on('tick', () => {
                links
                    .attr('x1', d => d.source.x)
                    .attr('y1', d => d.source.y)
                    .attr('x2', d => d.target.x)
                    .attr('y2', d => d.target.y);
                    
                nodes
                    .attr('transform', d => `translate(${d.x},${d.y})`);
            });
            
        this.simulation.force('link')
            .links(this.data.links);
            
        // Add click handler for nodes
        nodes.on('click', (event, d) => {
            this.handleNodeClick(d);
        });
    }
    
    getNodeColor(type) {
        const colors = {
            'repository': '#3498db',
            'module': '#2ecc71'
        };
        return colors[type] || '#95a5a6';
    }
    
    setLayout(layout) {
        this.currentLayout = layout;
        
        if (layout === 'tree') {
            this.renderTreeLayout();
        } else if (layout === 'radial') {
            this.renderRadialLayout();
        } else {
            this.render(); // Default force layout
        }
    }
    
    renderTreeLayout() {
        const treeLayout = d3.tree()
            .size([this.height - 100, this.width - 160]);
            
        const root = d3.stratify()
            .id(d => d.id)
            .parentId(d => {
                const link = this.data.links.find(l => l.target === d.id);
                return link ? link.source : null;
            })(this.data.nodes);
            
        const treeData = treeLayout(root);
        
        // Update node positions based on tree layout
        this.simulation.stop();
        this.data.nodes.forEach(node => {
            const treeNode = treeData.descendants().find(n => n.id === node.id);
            if (treeNode) {
                node.x = treeNode.y; // Swap x and y for horizontal layout
                node.y = treeNode.x;
            }
        });
        
        this.render();
    }
    
    renderRadialLayout() {
        const radius = Math.min(this.width, this.height) / 2 - 100;
        
        const radialLayout = d3.tree()
            .size([2 * Math.PI, radius]);
            
        const root = d3.stratify()
            .id(d => d.id)
            .parentId(d => {
                const link = this.data.links.find(l => l.target === d.id);
                return link ? link.source : null;
            })(this.data.nodes);
            
        const radialData = radialLayout(root);
        
        // Update node positions based on radial layout
        this.simulation.stop();
        this.data.nodes.forEach(node => {
            const radialNode = radialData.descendants().find(n => n.id === node.id);
            if (radialNode) {
                node.x = radialNode.x * Math.cos(radialNode.y) + this.width / 2;
                node.y = radialNode.x * Math.sin(radialNode.y) + this.height / 2;
            }
        });
        
        this.render();
    }
    
    dragStarted(event) {
        if (!event.active) this.simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }
    
    dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }
    
    dragEnded(event) {
        if (!event.active) this.simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
    }
    
    handleNodeClick(node) {
        const event = new CustomEvent('nodeSelected', { detail: node });
        window.dispatchEvent(event);
    }
    
    resize() {
        this.width = this.container.node().getBoundingClientRect().width;
        this.height = this.container.node().getBoundingClientRect().height;
        
        this.svg
            .attr('width', this.width)
            .attr('height', this.height);
            
        this.simulation.force('center', d3.forceCenter(this.width / 2, this.height / 2));
        this.simulation.alpha(0.3).restart();
    }
} 