package graph

import (
	"container/heap"
	"sync"
)

// NodeQueue is a priority queue for graph traversal
type NodeQueue struct {
	nodes []*Node
	cost  map[string]int
	index map[string]int
}

func (q *NodeQueue) Len() int { return len(q.nodes) }

func (q *NodeQueue) Less(i, j int) bool {
	return q.cost[q.nodes[i].ID] < q.cost[q.nodes[j].ID]
}

func (q *NodeQueue) Swap(i, j int) {
	q.nodes[i], q.nodes[j] = q.nodes[j], q.nodes[i]
	q.index[q.nodes[i].ID] = i
	q.index[q.nodes[j].ID] = j
}

func (q *NodeQueue) Push(x interface{}) {
	n := x.(*Node)
	q.index[n.ID] = len(q.nodes)
	q.nodes = append(q.nodes, n)
}

func (q *NodeQueue) Pop() interface{} {
	old := q.nodes
	n := len(old)
	x := old[n-1]
	q.nodes = old[0 : n-1]
	delete(q.index, x.ID)
	return x
}

// OptimizedTraversal performs an optimized traversal of the graph
func (g *Graph) OptimizedTraversal(start string, visitor func(*Node)) {
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process node in parallel if it has multiple children
	var processNode func(nodeID string)
	processNode = func(nodeID string) {
		defer wg.Done()

		mu.Lock()
		if visited[nodeID] {
			mu.Unlock()
			return
		}
		visited[nodeID] = true
		node := g.Nodes[nodeID]
		mu.Unlock()

		if node == nil {
			return
		}

		visitor(node)

		// Get all outgoing edges
		var edges []*Edge
		for _, edge := range g.Edges {
			if edge.Source == nodeID {
				edges = append(edges, edge)
			}
		}

		// Process children in parallel if there are multiple
		if len(edges) > 1 {
			for _, edge := range edges {
				wg.Add(1)
				go processNode(edge.Target)
			}
		} else if len(edges) == 1 {
			// Process single child in the same goroutine
			wg.Add(1)
			processNode(edges[0].Target)
		}
	}

	wg.Add(1)
	go processNode(start)
	wg.Wait()
}

// FindShortestPaths finds the shortest paths from start to all other nodes
func (g *Graph) FindShortestPaths(start string) map[string][]string {
	paths := make(map[string][]string)
	cost := make(map[string]int)
	queue := &NodeQueue{
		cost:  cost,
		index: make(map[string]int),
	}

	// Initialize costs
	for id := range g.Nodes {
		if id == start {
			cost[id] = 0
		} else {
			cost[id] = int(^uint(0) >> 1) // Max int
		}
	}

	// Initialize queue with start node
	heap.Push(queue, g.Nodes[start])
	paths[start] = []string{start}

	for queue.Len() > 0 {
		node := heap.Pop(queue).(*Node)
		nodeCost := cost[node.ID]

		// Process all outgoing edges
		for _, edge := range g.Edges {
			if edge.Source != node.ID {
				continue
			}

			target := edge.Target
			newCost := nodeCost + 1 // Each edge has weight 1

			if newCost < cost[target] {
				cost[target] = newCost
				newPath := make([]string, len(paths[node.ID]))
				copy(newPath, paths[node.ID])
				paths[target] = append(newPath, target)

				// Update queue
				heap.Push(queue, g.Nodes[target])
			}
		}
	}

	return paths
}

// FindStronglyConnectedComponents finds strongly connected components using Kosaraju's algorithm
func (g *Graph) FindStronglyConnectedComponents() [][]string {
	// First DFS pass - get finishing times
	visited := make(map[string]bool)
	finish := make([]string, 0, len(g.Nodes))
	var finishOrder func(string)

	finishOrder = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true

		for _, edge := range g.Edges {
			if edge.Source == nodeID {
				finishOrder(edge.Target)
			}
		}

		finish = append(finish, nodeID)
	}

	for id := range g.Nodes {
		if !visited[id] {
			finishOrder(id)
		}
	}

	// Create transpose graph
	transpose := make(map[string][]string)
	for _, edge := range g.Edges {
		transpose[edge.Target] = append(transpose[edge.Target], edge.Source)
	}

	// Second DFS pass - find SCCs
	visited = make(map[string]bool)
	var components [][]string
	var component []string

	var collect func(string)
	collect = func(nodeID string) {
		if visited[nodeID] {
			return
		}
		visited[nodeID] = true
		component = append(component, nodeID)

		for _, neighbor := range transpose[nodeID] {
			collect(neighbor)
		}
	}

	// Process nodes in reverse finishing order
	for i := len(finish) - 1; i >= 0; i-- {
		if !visited[finish[i]] {
			component = make([]string, 0)
			collect(finish[i])
			components = append(components, component)
		}
	}

	return components
} 