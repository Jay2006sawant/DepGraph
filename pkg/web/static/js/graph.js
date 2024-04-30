class DependencyGraph {
    constructor(containerId) {
        this.container = d3.select(containerId);
        this.width = this.container.node().getBoundingClientRect().width;
        this.height = this.container.node().getBoundingClientRect().height;
        this.currentLayout = 'force';
        this.nodeCache = new Map();
        this.linkCache = new Map();
        this.transform = d3.zoomIdentity;
        
        this.svg = this.container.append('svg')
            .attr('width', this.width)
            .attr('height', this.height);
            
        this.g = this.svg.append('g');
        
        // Add zoom behavior with performance optimization
        this.zoom = d3.zoom()
            .scaleExtent([0.1, 4])
            .on('zoom', (event) => {
                this.transform = event.transform;
                this.g.attr('transform', this.transform);
                this.updateVisibility();
            });
            
        this.svg.call(this.zoom);
        
        // Initialize forces with optimized parameters
        this.initializeForces();
        
        // Add WebGL renderer for large graphs
        this.initWebGL();
        
        // Add worker for layout calculations
        this.layoutWorker = new Worker('/static/js/layout-worker.js');
        this.layoutWorker.onmessage = (e) => this.updateLayout(e.data);
        
        // Performance monitoring
        this.fps = 0;
        this.lastFrameTime = performance.now();
        this.monitorPerformance();
    }
    
    initializeForces() {
        this.simulation = d3.forceSimulation()
            .force('link', d3.forceLink().id(d => d.id).distance(50))
            .force('charge', d3.forceManyBody().strength(-100).distanceMax(200))
            .force('center', d3.forceCenter(this.width / 2, this.height / 2))
            .force('collision', d3.forceCollide().radius(30))
            .alphaDecay(0.02) // Slower decay for smoother animation
            .velocityDecay(0.4); // Reduced for more stable movement
            
        // Optimize force calculation frequency
        this.simulation.on('tick', () => {
            if (performance.now() - this.lastFrameTime > 16.7) { // ~60fps
                this.updatePositions();
                this.lastFrameTime = performance.now();
            }
        });
    }
    
    initWebGL() {
        try {
            this.renderer = new THREE.WebGLRenderer({ antialias: true });
            this.renderer.setSize(this.width, this.height);
            this.container.node().appendChild(this.renderer.domElement);
            this.renderer.domElement.style.display = 'none';
            
            this.scene = new THREE.Scene();
            this.camera = new THREE.PerspectiveCamera(75, this.width / this.height, 0.1, 1000);
            this.camera.position.z = 500;
            
            // Create geometry for nodes and links
            this.nodeGeometry = new THREE.SphereGeometry(5, 16, 16);
            this.nodeMaterial = new THREE.MeshBasicMaterial({ vertexColors: true });
            this.linkMaterial = new THREE.LineBasicMaterial({ vertexColors: true });
            
            this.webglEnabled = true;
        } catch (e) {
            console.warn('WebGL initialization failed, falling back to SVG');
            this.webglEnabled = false;
        }
    }
    
    async loadData() {
        try {
            const response = await fetch('/api/graph');
            const data = await response.json();
            this.data = this.optimizeData(data);
            this.render();
        } catch (error) {
            console.error('Error loading graph data:', error);
        }
    }
    
    optimizeData(data) {
        // Create indexed data structures for faster lookups
        const nodeIndex = new Map(data.nodes.map(n => [n.id, n]));
        const linkIndex = new Map();
        
        data.links.forEach(l => {
            const key = `${l.source}-${l.target}`;
            linkIndex.set(key, l);
        });
        
        // Pre-calculate node degrees
        data.nodes.forEach(n => {
            n.inDegree = data.links.filter(l => l.target === n.id).length;
            n.outDegree = data.links.filter(l => l.source === n.id).length;
        });
        
        return {
            nodes: data.nodes,
            links: data.links,
            nodeIndex,
            linkIndex
        };
    }
    
    updateVisibility() {
        if (!this.data) return;
        
        // Calculate visible area
        const visibleArea = {
            x1: -this.transform.x / this.transform.k,
            y1: -this.transform.y / this.transform.k,
            x2: (this.width - this.transform.x) / this.transform.k,
            y2: (this.height - this.transform.y) / this.transform.k
        };
        
        // Update node visibility
        this.g.selectAll('.node')
            .style('display', d => {
                const visible = d.x >= visibleArea.x1 && d.x <= visibleArea.x2 &&
                              d.y >= visibleArea.y1 && d.y <= visibleArea.y2;
                return visible ? null : 'none';
            });
            
        // Update link visibility
        this.g.selectAll('.link')
            .style('display', d => {
                const sourceVisible = d.source.x >= visibleArea.x1 && d.source.x <= visibleArea.x2 &&
                                   d.source.y >= visibleArea.y1 && d.source.y <= visibleArea.y2;
                const targetVisible = d.target.x >= visibleArea.x1 && d.target.x <= visibleArea.x2 &&
                                   d.target.y >= visibleArea.y1 && d.target.y <= visibleArea.y2;
                return (sourceVisible || targetVisible) ? null : 'none';
            });
    }
    
    updatePositions() {
        if (!this.data) return;
        
        if (this.webglEnabled && this.data.nodes.length > 1000) {
            this.updateWebGLPositions();
        } else {
            this.updateSVGPositions();
        }
    }
    
    updateSVGPositions() {
        this.g.selectAll('.node')
            .attr('transform', d => `translate(${d.x},${d.y})`);
            
        this.g.selectAll('.link')
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);
    }
    
    updateWebGLPositions() {
        this.nodeObjects.forEach((obj, i) => {
            const node = this.data.nodes[i];
            obj.position.set(node.x, node.y, 0);
        });
        
        this.linkObjects.forEach((obj, i) => {
            const link = this.data.links[i];
            const positions = obj.geometry.attributes.position;
            positions.array[0] = link.source.x;
            positions.array[1] = link.source.y;
            positions.array[3] = link.target.x;
            positions.array[4] = link.target.y;
            positions.needsUpdate = true;
        });
        
        this.renderer.render(this.scene, this.camera);
    }
    
    setLayout(layout) {
        this.currentLayout = layout;
        this.layoutWorker.postMessage({
            type: layout,
            nodes: this.data.nodes,
            links: this.data.links
        });
    }
    
    updateLayout(positions) {
        this.data.nodes.forEach((node, i) => {
            node.x = positions[i].x;
            node.y = positions[i].y;
        });
        this.updatePositions();
    }
    
    monitorPerformance() {
        let frameCount = 0;
        let lastTime = performance.now();
        
        const updateFPS = () => {
            const currentTime = performance.now();
            const elapsed = currentTime - lastTime;
            
            if (elapsed >= 1000) {
                this.fps = Math.round((frameCount * 1000) / elapsed);
                frameCount = 0;
                lastTime = currentTime;
                
                // Adjust rendering method based on performance
                if (this.fps < 30 && this.data.nodes.length > 1000) {
                    this.enableWebGL();
                }
            }
            
            frameCount++;
            requestAnimationFrame(updateFPS);
        };
        
        updateFPS();
    }
    
    enableWebGL() {
        if (!this.webglEnabled) return;
        
        this.svg.style('display', 'none');
        this.renderer.domElement.style.display = null;
        
        // Convert SVG nodes to WebGL objects
        this.nodeObjects = this.data.nodes.map(node => {
            const mesh = new THREE.Mesh(this.nodeGeometry, this.nodeMaterial);
            mesh.position.set(node.x, node.y, 0);
            this.scene.add(mesh);
            return mesh;
        });
        
        // Convert SVG links to WebGL lines
        this.linkObjects = this.data.links.map(link => {
            const geometry = new THREE.BufferGeometry();
            const positions = new Float32Array([
                link.source.x, link.source.y, 0,
                link.target.x, link.target.y, 0
            ]);
            geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
            const line = new THREE.Line(geometry, this.linkMaterial);
            this.scene.add(line);
            return line;
        });
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
            
        if (this.webglEnabled) {
            this.renderer.setSize(this.width, this.height);
            this.camera.aspect = this.width / this.height;
            this.camera.updateProjectionMatrix();
        }
        
        this.simulation.force('center', d3.forceCenter(this.width / 2, this.height / 2));
        this.simulation.alpha(0.3).restart();
    }
} 