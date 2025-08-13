package bmssp

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// Graph generators for benchmarking

// generateRandomGraph creates a random directed graph with n nodes and approximately m edges
func generateRandomGraph(n, m int, maxWeight float64, seed int64) *Graph {
	rand.Seed(seed)
	g := NewGraph()
	
	edgeCount := 0
	for edgeCount < m {
		u := NodeID(rand.Intn(n))
		v := NodeID(rand.Intn(n))
		if u != v { // avoid self-loops
			w := Dist(rand.Float64() * maxWeight + 1) // weights between 1 and maxWeight+1
			g.AddEdge(u, v, w)
			edgeCount++
		}
	}
	
	return g
}

// generateGridGraph creates a 2D grid graph (good for testing shortest paths)
func generateGridGraph(width, height int) *Graph {
	g := NewGraph()
	
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			node := NodeID(i*width + j)
			
			// Right edge
			if j < width-1 {
				right := NodeID(i*width + j + 1)
				g.AddEdge(node, right, 1)
			}
			
			// Down edge  
			if i < height-1 {
				down := NodeID((i+1)*width + j)
				g.AddEdge(node, down, 1)
			}
			
			// Left edge
			if j > 0 {
				left := NodeID(i*width + j - 1)
				g.AddEdge(node, left, 1)
			}
			
			// Up edge
			if i > 0 {
				up := NodeID((i-1)*width + j)
				g.AddEdge(node, up, 1)
			}
		}
	}
	
	return g
}

// generateCompleteGraph creates a complete directed graph
func generateCompleteGraph(n int) *Graph {
	g := NewGraph()
	
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				weight := Dist(math.Abs(float64(i-j))) // distance-based weight
				g.AddEdge(NodeID(i), NodeID(j), weight)
			}
		}
	}
	
	return g
}

// Helper function to initialize distance map for all nodes in graph
func initializeDistanceMap(g *Graph, source NodeID) map[NodeID]Dist {
	dhat := make(map[NodeID]Dist)
	
	// Add all nodes that appear as sources
	for u := range g.adj {
		dhat[u] = INF
	}
	
	// Add all nodes that appear as destinations
	for _, edges := range g.adj {
		for _, edge := range edges {
			dhat[edge.V] = INF
		}
	}
	
	dhat[source] = 0
	return dhat
}

// Benchmark Dijkstra on random graphs
func BenchmarkDijkstraRandom100(b *testing.B) {
	g := generateRandomGraph(100, 500, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		DijkstraSingleSource(g, source, dhat)
	}
}

func BenchmarkDijkstraRandom500(b *testing.B) {
	g := generateRandomGraph(500, 2500, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		DijkstraSingleSource(g, source, dhat)
	}
}

func BenchmarkDijkstraRandom1000(b *testing.B) {
	g := generateRandomGraph(1000, 5000, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		DijkstraSingleSource(g, source, dhat)
	}
}

// Benchmark BMSSP on random graphs
func BenchmarkBMSSPRandom100(b *testing.B) {
	g := generateRandomGraph(100, 500, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BMSSP(2, 1000, S, 50, 1, g, dhat)
	}
}

func BenchmarkBMSSPRandom500(b *testing.B) {
	g := generateRandomGraph(500, 2500, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BMSSP(2, 1000, S, 100, 1, g, dhat)
	}
}

func BenchmarkBMSSPRandom1000(b *testing.B) {
	g := generateRandomGraph(1000, 5000, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BMSSP(2, 1000, S, 100, 1, g, dhat)
	}
}

// Benchmark on grid graphs (more structured)
func BenchmarkDijkstraGrid20x20(b *testing.B) {
	g := generateGridGraph(20, 20)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		DijkstraSingleSource(g, source, dhat)
	}
}

func BenchmarkBMSSPGrid20x20(b *testing.B) {
	g := generateGridGraph(20, 20)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BMSSP(2, 1000, S, 50, 1, g, dhat)
	}
}

func BenchmarkDijkstraGrid50x50(b *testing.B) {
	g := generateGridGraph(50, 50)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		DijkstraSingleSource(g, source, dhat)
	}
}

func BenchmarkBMSSPGrid50x50(b *testing.B) {
	g := generateGridGraph(50, 50)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BMSSP(2, 1000, S, 100, 1, g, dhat)
	}
}

// Benchmark BaseCase (bounded Dijkstra) vs full BMSSP
func BenchmarkBaseCaseRandom1000(b *testing.B) {
	g := generateRandomGraph(1000, 5000, 10.0, 42)
	source := NodeID(0)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dhat := initializeDistanceMap(g, source)
		S := NewNodeSet()
		S.Add(source)
		BaseCase(1000, S, 1000, g, dhat)
	}
}

// Test to verify algorithms produce the same results
func TestAlgorithmEquivalence(t *testing.T) {
	sizes := []int{50, 100}
	
	for _, n := range sizes {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			g := generateRandomGraph(n, n*3, 10.0, 42)
			source := NodeID(0)
			
			// Run Dijkstra
			dhatDijkstra := initializeDistanceMap(g, source)
			DijkstraSingleSource(g, source, dhatDijkstra)
			
			// Run BMSSP
			dhatBMSSP := initializeDistanceMap(g, source)
			S := NewNodeSet()
			S.Add(source)
			BMSSP(2, 1000, S, n, 1, g, dhatBMSSP)
			
			// Compare results
			for node := range dhatDijkstra {
				dijkstraDist := dhatDijkstra[node]
				bmsspDist := dhatBMSSP[node]
				
				if math.Abs(float64(dijkstraDist-bmsspDist)) > 1e-9 {
					t.Errorf("Node %d: Dijkstra=%v, BMSSP=%v", node, dijkstraDist, bmsspDist)
				}
			}
		})
	}
}