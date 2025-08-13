package bmssp_test

import (
	"fmt"
	"math"

	"github.com/mfreeman451/bmssp-go"
)

// ExampleBMSSP demonstrates basic usage of the BMSSP algorithm.
func ExampleBMSSP() {
	// Create a new graph
	g := bmssp.NewGraph()

	// Add edges (source, destination, weight)
	g.AddEdge(0, 1, 2.0)
	g.AddEdge(0, 2, 5.0)
	g.AddEdge(1, 3, 4.0)
	g.AddEdge(2, 3, 1.0)
	g.AddEdge(1, 4, 1.0)
	g.AddEdge(3, 5, 3.0)
	g.AddEdge(4, 5, 2.0)

	// Initialize distance map with infinity for all nodes
	dhat := make(map[bmssp.NodeID]bmssp.Dist)
	for i := 0; i < 6; i++ {
		dhat[bmssp.NodeID(i)] = bmssp.Dist(math.Inf(1))
	}
	dhat[0] = 0 // source node has distance 0

	// Run BMSSP algorithm
	dhat = bmssp.BMSSPSingleSource(g, 0, 1000)

	// Count reachable nodes
	reachable := 0
	for _, dist := range dhat {
		if dist != bmssp.Dist(math.Inf(1)) {
			reachable++
		}
	}
	fmt.Printf("Reached %d nodes\n", reachable)

	// Print shortest distances
	fmt.Println("Shortest distances from node 0:")
	for i := 0; i < 6; i++ {
		nodeID := bmssp.NodeID(i)
		dist := dhat[nodeID]
		if dist == bmssp.Dist(math.Inf(1)) {
			fmt.Printf("  Node %d: unreachable\n", i)
		} else {
			fmt.Printf("  Node %d: %.0f\n", i, dist)
		}
	}

	// Output:
	// Reached 6 nodes
	// Shortest distances from node 0:
	//   Node 0: 0
	//   Node 1: 2
	//   Node 2: 5
	//   Node 3: 6
	//   Node 4: 3
	//   Node 5: 5
}

// ExampleDijkstra demonstrates the standard Dijkstra algorithm for comparison.
func ExampleDijkstra() {
	// Create the same graph as above
	g := bmssp.NewGraph()
	g.AddEdge(0, 1, 2.0)
	g.AddEdge(0, 2, 5.0)
	g.AddEdge(1, 3, 4.0)
	g.AddEdge(2, 3, 1.0)

	// Run Dijkstra's algorithm
	distances := bmssp.Dijkstra(g, 0)

	fmt.Println("Dijkstra distances from node 0:")
	for i := 0; i < 4; i++ {
		nodeID := bmssp.NodeID(i)
		if dist, exists := distances[nodeID]; exists {
			fmt.Printf("  Node %d: %.0f\n", i, dist)
		}
	}

	// Output:
	// Dijkstra distances from node 0:
	//   Node 0: 0
	//   Node 1: 2
	//   Node 2: 5
	//   Node 3: 6
}

// ExampleBMSSP_gridGraph demonstrates BMSSP on a structured grid graph.
func ExampleBMSSP_gridGraph() {
	// Create a 3x3 grid graph
	g := bmssp.NewGraph()

	// Add edges for a 3x3 grid (nodes 0-8)
	// Grid layout:
	// 0 - 1 - 2
	// |   |   |
	// 3 - 4 - 5
	// |   |   |
	// 6 - 7 - 8

	edges := [][3]int{
		// Horizontal edges
		{0, 1, 1}, {1, 2, 1},
		{3, 4, 1}, {4, 5, 1},
		{6, 7, 1}, {7, 8, 1},
		// Vertical edges
		{0, 3, 1}, {1, 4, 1}, {2, 5, 1},
		{3, 6, 1}, {4, 7, 1}, {5, 8, 1},
	}

	// Add bidirectional edges
	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		g.AddEdge(bmssp.NodeID(u), bmssp.NodeID(v), bmssp.Dist(w))
		g.AddEdge(bmssp.NodeID(v), bmssp.NodeID(u), bmssp.Dist(w))
	}

	// Initialize distances
	dhat := make(map[bmssp.NodeID]bmssp.Dist)
	for i := 0; i < 9; i++ {
		dhat[bmssp.NodeID(i)] = bmssp.Dist(math.Inf(1))
	}
	dhat[0] = 0 // start from top-left corner

	// Run BMSSP - works particularly well on structured graphs
	dhat = bmssp.BMSSPSingleSource(g, 0, 100)

	fmt.Println("Grid distances from corner (0,0):")
	for i := 0; i < 9; i++ {
		nodeID := bmssp.NodeID(i)
		fmt.Printf("  Node %d: %.0f\n", i, dhat[nodeID])
	}

	// Output:
	// Grid distances from corner (0,0):
	//   Node 0: 0
	//   Node 1: 1
	//   Node 2: 2
	//   Node 3: 1
	//   Node 4: 2
	//   Node 5: 3
	//   Node 6: 2
	//   Node 7: 3
	//   Node 8: 4
}

// Example_parameterTuning shows how to choose parameters for different graph sizes.
func Example_parameterTuning() {
	// Parameter guidelines for different graph sizes

	// Small graph parameters
	fmt.Println("Small graph parameters:")
	fmt.Println("  l=1, k=50, t=1")

	// Medium graph (100-1000 nodes): increase exploration
	fmt.Println("Medium graph parameters:")
	fmt.Println("  l=2, k=100, t=1")

	// Large graph (> 1000 nodes): more aggressive exploration
	fmt.Println("Large graph parameters:")
	fmt.Println("  l=2, k=200, t=1")

	// Structured graphs (grids, trees): can use higher k
	fmt.Println("Structured graph parameters:")
	fmt.Println("  l=1, k=200, t=1")

	// Output:
	// Small graph parameters:
	//   l=1, k=50, t=1
	// Medium graph parameters:
	//   l=2, k=100, t=1
	// Large graph parameters:
	//   l=2, k=200, t=1
	// Structured graph parameters:
	//   l=1, k=200, t=1
}
