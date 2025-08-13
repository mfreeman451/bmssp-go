package bmssp

import "container/heap"

// Standard Dijkstra's algorithm implementation for comparison

type dijkstraItem struct {
	node NodeID
	dist Dist
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

// Dijkstra implements standard Dijkstra's shortest path algorithm
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
			if _, exists := dist[edge.V]; !exists {
				dist[edge.V] = INF
				items[edge.V] = &dijkstraItem{node: edge.V, dist: INF}
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
			v := edge.V
			alt := dist[u] + edge.W
			
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

// DijkstraSingleSource is a simpler version that just computes distances from one source
// and stores them in the provided dhat map (compatible with BMSSP interface)
func DijkstraSingleSource(g *Graph, source NodeID, dhat map[NodeID]Dist) {
	visited := make(map[NodeID]bool)
	items := make(map[NodeID]*dijkstraItem)
	
	// Initialize items for all nodes that appear in dhat
	for node := range dhat {
		items[node] = &dijkstraItem{node: node, dist: dhat[node]}
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
			v := edge.V
			alt := dhat[u] + edge.W
			
			if alt < dhat[v] {
				dhat[v] = alt
				if !visited[v] && items[v] != nil && items[v].index >= 0 {
					pq.update(items[v], alt)
				}
			}
		}
	}
}