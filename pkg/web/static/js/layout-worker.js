// Web worker for handling graph layout calculations
importScripts('https://d3js.org/d3.v7.min.js');

self.onmessage = function(e) {
    const { type, nodes, links } = e.data;
    
    switch (type) {
        case 'tree':
            calculateTreeLayout(nodes, links);
            break;
        case 'radial':
            calculateRadialLayout(nodes, links);
            break;
        case 'force':
            calculateForceLayout(nodes, links);
            break;
        default:
            console.warn('Unknown layout type:', type);
    }
};

function calculateTreeLayout(nodes, links) {
    const width = 2000; // Default size, can be adjusted
    const height = 1500;
    
    const treeLayout = d3.tree()
        .size([height - 100, width - 160]);
        
    const root = d3.stratify()
        .id(d => d.id)
        .parentId(d => {
            const link = links.find(l => l.target === d.id);
            return link ? link.source : null;
        })(nodes);
        
    const treeData = treeLayout(root);
    
    // Extract positions
    const positions = nodes.map(node => {
        const treeNode = treeData.descendants().find(n => n.id === node.id);
        return {
            id: node.id,
            x: treeNode ? treeNode.y : node.x, // Swap x and y for horizontal layout
            y: treeNode ? treeNode.x : node.y
        };
    });
    
    self.postMessage(positions);
}

function calculateRadialLayout(nodes, links) {
    const radius = 750; // Default radius, can be adjusted
    
    const radialLayout = d3.tree()
        .size([2 * Math.PI, radius]);
        
    const root = d3.stratify()
        .id(d => d.id)
        .parentId(d => {
            const link = links.find(l => l.target === d.id);
            return link ? link.source : null;
        })(nodes);
        
    const radialData = radialLayout(root);
    
    // Extract positions
    const positions = nodes.map(node => {
        const radialNode = radialData.descendants().find(n => n.id === node.id);
        return {
            id: node.id,
            x: radialNode ? radialNode.x * Math.cos(radialNode.y) + 1000 : node.x,
            y: radialNode ? radialNode.x * Math.sin(radialNode.y) + 750 : node.y
        };
    });
    
    self.postMessage(positions);
}

function calculateForceLayout(nodes, links) {
    const simulation = d3.forceSimulation(nodes)
        .force('link', d3.forceLink(links).id(d => d.id).distance(50))
        .force('charge', d3.forceManyBody().strength(-100).distanceMax(200))
        .force('center', d3.forceCenter(1000, 750))
        .force('collision', d3.forceCollide().radius(30))
        .stop();
        
    // Run simulation
    for (let i = 0; i < 300; ++i) simulation.tick();
    
    // Extract final positions
    const positions = nodes.map(node => ({
        id: node.id,
        x: node.x,
        y: node.y
    }));
    
    self.postMessage(positions);
}

// Helper function for optimized force calculations
function optimizeForceCalculation(nodes, links) {
    // Build node and link indices for faster lookups
    const nodeIndex = new Map(nodes.map(n => [n.id, n]));
    const linksByNode = new Map(nodes.map(n => [n.id, []]));
    
    links.forEach(link => {
        linksByNode.get(link.source).push(link);
        linksByNode.get(link.target).push(link);
    });
    
    return {
        getNode: id => nodeIndex.get(id),
        getNodeLinks: id => linksByNode.get(id) || []
    };
} 