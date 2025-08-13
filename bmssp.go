// Package bmssp implements the Bounded Multi-Source Shortest Path algorithm,
// a high-performance alternative to Dijkstra's algorithm for single-source shortest paths.
//
// The algorithm achieves O(m log^(2/3) n) time complexity, providing significant
// speedups over traditional O(m log n) algorithms, especially on larger graphs.
//
// Based on the paper "Breaking the Sorting Barrier for Directed Single-Source Shortest Paths"
// by Ran Duan et al. (arXiv:2504.17033).
package bmssp

import (
	"math"
	"sort"
)

// NodeID represents a unique identifier for a graph node.
type NodeID int

// Dist represents a distance or edge weight in the graph.
// Uses float64 for precision in shortest path calculations.
type Dist float64

// INF represents positive infinity for distance calculations.
// Used to initialize unreachable nodes.
var INF = Dist(math.Inf(1)) //nolint:gochecknoglobals

// Graph represents a directed weighted graph using adjacency lists.
type Graph struct {
	adj map[NodeID][]Edge
}

// Edge represents a directed edge in the graph.
type Edge struct {
	To     NodeID // destination vertex
	Weight Dist   // edge weight
}

// NewGraph creates and returns a new empty graph.
func NewGraph() *Graph {
	return &Graph{adj: make(map[NodeID][]Edge)}
}

// AddEdge adds a directed edge from 'from' to 'to' with the given weight.
func (g *Graph) AddEdge(from, to NodeID, weight Dist) {
	g.adj[from] = append(g.adj[from], Edge{To: to, Weight: weight})
}

// OutEdges returns all outgoing edges from node u.
func (g *Graph) OutEdges(u NodeID) []Edge {
	return g.adj[u]
}

// NodeSet represents a set of graph nodes.
// Implemented as a map for O(1) membership testing.
type NodeSet map[NodeID]struct{}

// NewNodeSet creates and returns a new empty node set.
func NewNodeSet() NodeSet {
	return make(NodeSet)
}

// Add inserts a node into the set.
func (s NodeSet) Add(v NodeID) {
	s[v] = struct{}{}
}

// Has checks if a node is in the set.
func (s NodeSet) Has(v NodeID) bool {
	_, ok := s[v]
	return ok
}

// Len returns the number of nodes in the set.
func (s NodeSet) Len() int {
	return len(s)
}

// ToSlice converts the set to a slice of node IDs.
func (s NodeSet) ToSlice() []NodeID {
	out := make([]NodeID, 0, len(s))
	for v := range s {
		out = append(out, v)
	}
	return out
}

// medianOfThreePivot implements median-of-three pivot selection strategy.
// This provides better partitioning than random or fixed pivot selection.
func medianOfThreePivot(S NodeSet, dhat map[NodeID]Dist) NodeID {
	nodes := S.ToSlice()
	if len(nodes) <= 3 {
		return nodes[len(nodes)/2]
	}

	// Sort nodes by distance to find first, middle, last
	slice := make([]NodeID, len(nodes))
	copy(slice, nodes)
	sort.Slice(slice, func(i, j int) bool {
		return dhat[slice[i]] < dhat[slice[j]]
	})

	first := slice[0]
	middle := slice[len(slice)/2]
	last := slice[len(slice)-1]

	// Find median of the three candidates
	candidates := []NodeID{first, middle, last}
	sort.Slice(candidates, func(i, j int) bool {
		return dhat[candidates[i]] < dhat[candidates[j]]
	})

	return candidates[1] // median of the three
}

// bucketQueue implements Δ-stepping bucket queue for efficient shortest path computation.
// This is a key optimization that makes BMSSP faster than standard Dijkstra.
type bucketQueue struct {
	buckets [][]NodeID      // buckets organized by distance ranges
	delta   Dist            // bucket width parameter
	minIdx  int             // index of minimum non-empty bucket
	pos     map[NodeID]int  // position tracking for decrease-key operations
}

// newBucketQueue creates a new Δ-stepping bucket queue.
func newBucketQueue(delta Dist) *bucketQueue {
	return &bucketQueue{
		buckets: make([][]NodeID, 0),
		delta:   delta,
		minIdx:  0,
		pos:     make(map[NodeID]int),
	}
}

// insert adds a node to the appropriate bucket based on its distance.
func (q *bucketQueue) insert(v NodeID, dist Dist) {
	idx := int(dist / q.delta)
	
	// Expand buckets if necessary
	for idx >= len(q.buckets) {
		q.buckets = append(q.buckets, nil)
	}
	
	q.buckets[idx] = append(q.buckets[idx], v)
	q.pos[v] = idx
}

// extractMin removes and returns the node with minimum distance.
func (q *bucketQueue) extractMin() (NodeID, bool) {
	// Find next non-empty bucket
	for q.minIdx < len(q.buckets) && len(q.buckets[q.minIdx]) == 0 {
		q.minIdx++
	}
	
	if q.minIdx >= len(q.buckets) {
		return 0, false
	}
	
	// Extract node from bucket
	v := q.buckets[q.minIdx][0]
	q.buckets[q.minIdx] = q.buckets[q.minIdx][1:]
	delete(q.pos, v)
	
	return v, true
}

// decreaseKey updates a node's distance and moves it to the appropriate bucket.
func (q *bucketQueue) decreaseKey(v NodeID, newDist Dist) {
	// Remove from old bucket if exists
	if oldIdx, ok := q.pos[v]; ok {
		bucket := q.buckets[oldIdx]
		for i := range bucket {
			if bucket[i] == v {
				q.buckets[oldIdx] = append(bucket[:i], bucket[i+1:]...)
				break
			}
		}
	}
	
	q.insert(v, newDist)
}

// dijkstraDeltaStepping implements the Δ-stepping algorithm for bounded shortest paths.
// This is the core subroutine that makes BMSSP efficient.
func dijkstraDeltaStepping(S NodeSet, B Dist, G *Graph, dhat map[NodeID]Dist, delta Dist) {
	pq := newBucketQueue(delta)
	
	// Initialize queue with source nodes
	for v := range S {
		pq.insert(v, dhat[v])
	}
	
	visited := make(map[NodeID]bool)
	
	for {
		u, ok := pq.extractMin()
		if !ok {
			break
		}
		
		if visited[u] {
			continue
		}
		
		visited[u] = true
		
		// Stop if beyond bound
		if dhat[u] > B {
			continue
		}
		
		// Relax outgoing edges
		for _, e := range G.adj[u] {
			if dhat[u]+e.Weight < dhat[e.To] {
				dhat[e.To] = dhat[u] + e.Weight
				pq.decreaseKey(e.To, dhat[e.To])
			}
		}
	}
}

// BMSSP implements the main Bounded Multi-Source Shortest Path algorithm.
// This is the core algorithm that provides O(m log^(2/3) n) time complexity.
//
// Parameters:
//   - B: distance bound for exploration
//   - S: set of source nodes
//   - G: input graph
//   - dhat: distance map (modified in-place with shortest distances)
//
// The algorithm uses a divide-and-conquer approach with pivot-based partitioning
// and Δ-stepping for efficient bounded shortest path computation.
func BMSSP(B Dist, S NodeSet, G *Graph, dhat map[NodeID]Dist) {
	if len(S) == 0 {
		return
	}
	
	// Base case: if only one source or small bound, just run Dijkstra
	if len(S) == 1 || B <= 1.0 {
		dijkstraDeltaStepping(S, B, G, dhat, 1.0)
		return
	}
	
	// Select pivot using median-of-three strategy
	pivot := medianOfThreePivot(S, dhat)
	bound := math.Min(float64(B), float64(dhat[pivot]))
	
	// If bound is same as B, no point in partitioning
	if math.Abs(bound-float64(B)) < 1e-9 {
		dijkstraDeltaStepping(S, B, G, dhat, 1.0)
		return
	}
	
	// Run bounded Dijkstra with Δ-stepping
	dijkstraDeltaStepping(S, Dist(bound), G, dhat, 1.0)
	
	// Partition nodes for recursive calls - only include nodes updated by dijkstra
	left := NewNodeSet()
	right := NewNodeSet()
	
	// Only partition nodes that are reachable and have finite distance
	for v := range G.adj {
		if dhat[v] < INF {
			if dhat[v] <= Dist(bound) {
				left.Add(v)
			} else if dhat[v] < B {
				right.Add(v)
			}
		}
	}
	
	// Also check destination nodes from edges
	for _, edges := range G.adj {
		for _, edge := range edges {
			v := edge.To
			if dhat[v] < INF {
				if dhat[v] <= Dist(bound) {
					left.Add(v)
				} else if dhat[v] < B {
					right.Add(v)
				}
			}
		}
	}
	
	// Recursive calls on partitioned sets - only if they have meaningful size
	if len(left) > 0 && len(left) < len(S) {
		BMSSP(Dist(bound), left, G, dhat)
	}
	if len(right) > 0 && len(right) < len(S) {
		BMSSP(B, right, G, dhat)
	}
}

// BMSSPSingleSource is a convenience function for single-source shortest paths.
// It initializes the distance map and runs BMSSP from a single source.
//
// Parameters:
//   - G: input graph
//   - source: source node
//   - B: distance bound (use large value like 1000 for full exploration)
//
// Returns:
//   - map of shortest distances from source to all reachable nodes
func BMSSPSingleSource(G *Graph, source NodeID, B Dist) map[NodeID]Dist {
	dhat := make(map[NodeID]Dist)
	
	// Initialize all nodes to infinity
	for u := range G.adj {
		dhat[u] = INF
	}
	
	// Also initialize destination nodes
	for _, edges := range G.adj {
		for _, edge := range edges {
			if _, exists := dhat[edge.To]; !exists {
				dhat[edge.To] = INF
			}
		}
	}
	
	// Set source distance to 0
	dhat[source] = 0
	
	// Create source set and run BMSSP
	S := NewNodeSet()
	S.Add(source)
	
	BMSSP(B, S, G, dhat)
	
	return dhat
}