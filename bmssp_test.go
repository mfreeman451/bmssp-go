package bmssp

import (
	"math"
	"testing"
)

func TestBMSSP_PaperExample(t *testing.T) {
	// Build the Figure 1 example graph:
	g := NewGraph()
	edges := []struct {
		u, v NodeID
		w    Dist
	}{
		{0, 1, 2}, {0, 2, 5}, {1, 3, 4}, {2, 3, 1},
		{1, 4, 1}, {3, 5, 3}, {4, 5, 2}, {5, 6, 1},
		{6, 7, 1},
	}
	for _, e := range edges {
		g.AddEdge(e.u, e.v, e.w)
	}

	// Use the new convenience function
	dhat := BMSSPSingleSource(g, 0, 1000)

	expected := []Dist{0, 2, 5, 6, 3, 5, 6, 7}
	for i := 0; i < 8; i++ {
		if math.Abs(float64(dhat[NodeID(i)]-expected[i])) > 1e-9 {
			t.Errorf("node %d: expected dist %.0f, got %v", i, expected[i], dhat[NodeID(i)])
		}
	}
}
