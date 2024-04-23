document.addEventListener('DOMContentLoaded', () => {
    const graph = new DependencyGraph('#graph');
    graph.loadData();

    // Layout selector
    const layoutSelect = document.getElementById('layout');
    layoutSelect.addEventListener('change', (e) => {
        graph.setLayout(e.target.value);
    });

    // Toggle stats panel
    const toggleStats = document.getElementById('toggleStats');
    const sidebar = document.getElementById('sidebar');
    const statsPanel = document.getElementById('stats');
    const securityPanel = document.getElementById('security');

    toggleStats.addEventListener('click', () => {
        sidebar.classList.toggle('hidden');
        statsPanel.classList.remove('hidden');
        securityPanel.classList.add('hidden');
        toggleStats.textContent = sidebar.classList.contains('hidden') ? 'Show Stats' : 'Hide Stats';
    });

    // Toggle security panel
    const toggleSecurity = document.getElementById('toggleSecurity');
    toggleSecurity.addEventListener('click', async () => {
        sidebar.classList.remove('hidden');
        statsPanel.classList.add('hidden');
        securityPanel.classList.remove('hidden');

        // Load security scan results
        try {
            const response = await fetch('/api/security');
            const results = await response.json();
            displaySecurityResults(results);
        } catch (error) {
            console.error('Error loading security scan results:', error);
        }
    });

    // Handle node selection
    window.addEventListener('nodeSelected', (e) => {
        const node = e.detail;
        displayNodeDetails(node);
    });

    // Handle window resize
    let resizeTimeout;
    window.addEventListener('resize', () => {
        clearTimeout(resizeTimeout);
        resizeTimeout = setTimeout(() => {
            graph.resize();
        }, 250);
    });
});

function displaySecurityResults(results) {
    const container = document.getElementById('scanResults');
    container.innerHTML = '';

    if (results.length === 0) {
        container.innerHTML = '<p>No security issues found.</p>';
        return;
    }

    const table = document.createElement('table');
    table.innerHTML = `
        <thead>
            <tr>
                <th>Module</th>
                <th>Version</th>
                <th>Risk Level</th>
                <th>Fix</th>
            </tr>
        </thead>
        <tbody>
        </tbody>
    `;

    const tbody = table.querySelector('tbody');
    results.forEach(result => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${result.module}</td>
            <td>${result.version}</td>
            <td class="risk-${result.riskLevel.toLowerCase()}">${result.riskLevel}</td>
            <td>${result.recommendedFix}</td>
        `;
        tbody.appendChild(row);
    });

    container.appendChild(table);
}

function displayNodeDetails(node) {
    const container = document.getElementById('nodeInfo');
    const detailsPanel = document.getElementById('details');
    const sidebar = document.getElementById('sidebar');

    sidebar.classList.remove('hidden');
    detailsPanel.classList.remove('hidden');

    container.innerHTML = `
        <dl>
            <dt>ID:</dt>
            <dd>${node.id}</dd>
            <dt>Label:</dt>
            <dd>${node.label}</dd>
            <dt>Type:</dt>
            <dd>${node.type}</dd>
        </dl>
    `;

    // Load additional details based on node type
    if (node.type === 'module') {
        loadModuleDetails(node.id);
    } else if (node.type === 'repository') {
        loadRepositoryDetails(node.id);
    }
}

async function loadModuleDetails(moduleId) {
    try {
        const response = await fetch(`/api/impact?module=${moduleId}`);
        const impact = await response.json();
        
        const container = document.getElementById('nodeInfo');
        const impactDetails = document.createElement('div');
        impactDetails.innerHTML = `
            <h3>Impact Analysis</h3>
            <dl>
                <dt>Impact Score:</dt>
                <dd>${impact.impactScore.toFixed(2)}</dd>
                <dt>Breaking Changes:</dt>
                <dd>${impact.breakingChanges ? 'Yes' : 'No'}</dd>
                <dt>Affected Repositories:</dt>
                <dd>${impact.affectedRepos.length}</dd>
            </dl>
        `;
        
        container.appendChild(impactDetails);
    } catch (error) {
        console.error('Error loading module details:', error);
    }
}

async function loadRepositoryDetails(repoId) {
    try {
        const response = await fetch(`/api/chains?repo=${repoId}`);
        const chains = await response.json();
        
        const container = document.getElementById('nodeInfo');
        const chainsDetails = document.createElement('div');
        chainsDetails.innerHTML = `
            <h3>Dependency Chains</h3>
            <ul>
                ${chains.map(chain => `
                    <li>
                        Length: ${chain.length}
                        ${chain.circular ? ' (Circular)' : ''}
                        <br>
                        Path: ${chain.path.join(' → ')}
                    </li>
                `).join('')}
            </ul>
        `;
        
        container.appendChild(chainsDetails);
    } catch (error) {
        console.error('Error loading repository details:', error);
    }
} 