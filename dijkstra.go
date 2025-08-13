package bmssp

import "container/heap"

// This file implements standard Dijkstra's algorithm for performance comparison
// with the BMSSP algorithm.

type dijkstraItem struct {
	node  NodeID
	dist  Dist
	index int
}

type dijkstraHeap []*dijkstraItem

func (h dijkstraHeap) Len() int           { return len(h) }
func (h dijkstraHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h dijkstraHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *dijkstraHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*dijkstraItem)
	item.index = n
	*h = append(*h, item)
}

func (h *dijkstraHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (h *dijkstraHeap) update(item *dijkstraItem, dist Dist) {
	item.dist = dist
	heap.Fix(h, item.index)
}

// Dijkstra implements the standard Dijkstra's shortest path algorithm.
// This function is provided for performance comparison with BMSSP.
//
// Parameters:
//   - g: input graph
//   - source: source node for shortest path computation
//
// Returns:
//   - map of node IDs to their shortest distances from source
func Dijkstra(g *Graph, source NodeID) map[NodeID]Dist {
	dist := make(map[NodeID]Dist)
	visited := make(map[NodeID]bool)
	items := make(map[NodeID]*dijkstraItem)

	// Initialize all distances to infinity
	for u := range g.adj {
		dist[u] = INF
		items[u] = &dijkstraItem{node: u, dist: INF}
	}

	// Also check all destination nodes from edges
	for _, edges := range g.adj {
		for _, edge := range edges {
			if _, exists := dist[edge.To]; !exists {
				dist[edge.To] = INF
				items[edge.To] = &dijkstraItem{node: edge.To, dist: INF}
			}
		}
	}

	dist[source] = 0
	items[source].dist = 0

	// Create priority queue
	pq := make(dijkstraHeap, 0, len(items))
	for _, item := range items {
		heap.Push(&pq, item)
	}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*dijkstraItem)
		u := item.node

		if visited[u] {
			continue
		}
		visited[u] = true

		// Relax all outgoing edges
		for _, edge := range g.OutEdges(u) {
			v := edge.To
			alt := dist[u] + edge.Weight

			if alt < dist[v] {
				dist[v] = alt
				if !visited[v] && items[v].index >= 0 {
					pq.update(items[v], alt)
				}
			}
		}
	}

	return dist
}

// DijkstraSingleSource computes shortest distances using Dijkstra's algorithm
// with an interface compatible with BMSSP (uses pre-allocated distance map).
//
// Parameters:
//   - g: input graph
//   - source: source node for shortest path computation
//   - dhat: distance map (modified in-place)
func DijkstraSingleSource(g *Graph, source NodeID, dhat map[NodeID]Dist) {
	visited := make(map[NodeID]bool)
	items := make(map[NodeID]*dijkstraItem)

	// Initialize items for all nodes that appear in dhat
	for node := range dhat {
		if node == source {
			dhat[node] = 0
			items[node] = &dijkstraItem{node: node, dist: 0}
		} else if dhat[node] == 0 && node != source {
			dhat[node] = INF
			items[node] = &dijkstraItem{node: node, dist: INF}
		} else {
			items[node] = &dijkstraItem{node: node, dist: dhat[node]}
		}
	}

	// Create priority queue
	pq := make(dijkstraHeap, 0, len(items))
	for _, item := range items {
		heap.Push(&pq, item)
	}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*dijkstraItem)
		u := item.node

		if visited[u] {
			continue
		}
		visited[u] = true

		// Relax all outgoing edges
		for _, edge := range g.OutEdges(u) {
			v := edge.To
			alt := dhat[u] + edge.Weight

			if alt < dhat[v] {
				dhat[v] = alt
				if !visited[v] && items[v] != nil && items[v].index >= 0 {
					pq.update(items[v], alt)
				}
			}
		}
	}
}
