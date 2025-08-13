package bmssp

import (
	"math"
	"sort"
)

// Basic types
type NodeID int
type Dist float64

var INF = Dist(math.Inf(1))

type Edge struct {
	V NodeID
	W Dist
}

type Graph struct{ adj map[NodeID][]Edge }

func NewGraph() *Graph { return &Graph{adj: make(map[NodeID][]Edge)} }
func (g *Graph) AddEdge(u, v NodeID, w Dist) {
	g.adj[u] = append(g.adj[u], Edge{V: v, W: w})
}
func (g *Graph) OutEdges(u NodeID) []Edge { return g.adj[u] }

// NodeSet helpers
type NodeSet map[NodeID]struct{}

func NewNodeSet() NodeSet      { return make(NodeSet) }
func (s NodeSet) Add(n NodeID) { s[n] = struct{}{} }
func (s NodeSet) Len() int     { return len(s) }
func (s NodeSet) Clone() NodeSet {
	out := NewNodeSet()
	for v := range s {
		out.Add(v)
	}
	return out
}

// Algorithm 2 — BaseCase (bounded Dijkstra for singleton S)
func BaseCase(B Dist, S NodeSet, k int, g *Graph, dhat map[NodeID]Dist) (Dist, NodeSet) {
	type item struct {
		v NodeID
		d Dist
	}
	h := make([]item, 0)
	push := func(v NodeID, d Dist) {
		h = append(h, item{v, d})
		i := len(h) - 1
		for i > 0 {
			p := (i - 1) / 2
			if h[p].d <= h[i].d {
				break
			}
			h[p], h[i] = h[i], h[p]
			i = p
		}
	}
	pop := func() (NodeID, Dist, bool) {
		if len(h) == 0 {
			return 0, 0, false
		}
		it := h[0]
		h[0] = h[len(h)-1]
		h = h[:len(h)-1]
		i := 0
		for {
			l, r := 2*i+1, 2*i+2
			small := i
			if l < len(h) && h[l].d < h[small].d {
				small = l
			}
			if r < len(h) && h[r].d < h[small].d {
				small = r
			}
			if small == i {
				break
			}
			h[i], h[small] = h[small], h[i]
			i = small
		}
		return it.v, it.d, true
	}

	// init
	for x := range S {
		push(x, dhat[x])
	}
	Bprime := B
	U0 := NewNodeSet()
	visited := map[NodeID]struct{}{}
	for {
		u, du, ok := pop()
		if !ok || du >= B || U0.Len() > k {
			break
		}
		if _, seen := visited[u]; seen {
			continue
		}
		visited[u] = struct{}{}
		U0.Add(u)
		if du < Bprime {
			Bprime = du
		}
		for _, e := range g.OutEdges(u) {
			v := e.V
			newd := du + e.W
			if newd < dhat[v] && newd < B {
				dhat[v] = newd
				push(v, newd)
			}
		}
	}
	if U0.Len() <= k {
		return B, U0
	}
	Bmax := Dist(-math.MaxFloat64)
	for v := range U0 {
		if dhat[v] > Bmax {
			Bmax = dhat[v]
		}
	}
	U := NewNodeSet()
	for v := range U0 {
		if dhat[v] < Bmax {
			U.Add(v)
		}
	}
	return Bmax, U
}

// Lemma 3.3: bucketed queue structure D
type pair struct {
	v NodeID
	d Dist
}
type DStruct struct {
	M       int
	B       Dist
	D0, D1  [][]pair
	keyBest map[NodeID]Dist
}

func NewDStruct() *DStruct {
	return &DStruct{
		M:       1,
		B:       INF,
		D0:      nil,
		D1:      [][]pair{{}},
		keyBest: map[NodeID]Dist{},
	}
}
func (d *DStruct) Initialize(M int, B Dist) {
	d.M, d.B = M, B
	d.D0 = nil
	d.D1 = [][]pair{{}}
	d.keyBest = map[NodeID]Dist{}
}
func (d *DStruct) Insert(v NodeID, dist Dist) {
	if old, ok := d.keyBest[v]; ok && old <= dist {
		return
	}
	d.keyBest[v] = dist
	last := &d.D1[len(d.D1)-1]
	*last = append(*last, pair{v, dist})
	if len(*last) > d.M {
		// split
		arr := *last
		sort.Slice(arr, func(i, j int) bool { return arr[i].d < arr[j].d })
		mid := len(arr) / 2
		d.D1[len(d.D1)-1] = arr[:mid]
		d.D1 = append(d.D1, arr[mid:])
	}
}
func (d *DStruct) BatchPrepend(list []pair) {
	if len(list) == 0 {
		return
	}
	for i := len(list); i > 0; i -= d.M {
		start := i - d.M
		if start < 0 {
			start = 0
		}
		chunk := list[start:i]
		d.D0 = append([][]pair{chunk}, d.D0...)
		for _, p := range chunk {
			old, ok := d.keyBest[p.v]
			if !ok || p.d < old {
				d.keyBest[p.v] = p.d
			}
		}
	}
}
func (d *DStruct) NonEmpty() bool {
	if len(d.D0) > 0 && len(d.D0[0]) > 0 {
		return true
	}
	for _, block := range d.D1 {
		if len(block) > 0 {
			return true
		}
	}
	return false
}
func (d *DStruct) Pull() (Dist, NodeSet, bool) {
	S := NewNodeSet()
	if len(d.D0) > 0 && len(d.D0[0]) > 0 {
		block := d.D0[0]
		n := len(block)
		for i := n - 1; i >= 0 && len(S) < d.M; i-- {
			S.Add(block[i].v)
		}
		d.D0 = d.D0[1:]
		return block[0].d, S, true
	}
	if len(d.D1) > 0 && len(d.D1[0]) > 0 {
		block := d.D1[0]
		sort.Slice(block, func(i, j int) bool { return block[i].d < block[j].d })
		for i := 0; i < len(block) && i < d.M; i++ {
			S.Add(block[i].v)
		}
		d.D1 = d.D1[1:]
		return block[0].d, S, true
	}
	return d.B, S, false
}

// FindPivots exactly per Algorithm 1
func FindPivots(B Dist, S NodeSet, k int, g *Graph, dhat map[NodeID]Dist) (NodeSet, NodeSet) {
	W := NewNodeSet()
	for x := range S {
		W.Add(x)
	}
	curr := make(NodeSet)
	for x := range S {
		curr.Add(x)
	}
	for i := 0; i < k; i++ {
		next := NewNodeSet()
		for u := range curr {
			for _, e := range g.OutEdges(u) {
				v, w := e.V, e.W
				if dhat[u]+w < dhat[v] {
					dhat[v] = dhat[u] + w
					if dhat[v] < B {
						next.Add(v)
					}
					W.Add(v)
				}
			}
		}
		curr = next
		if W.Len() > k*len(S) {
			P := NewNodeSet()
			for x := range S {
				P.Add(x)
			}
			return P, W
		}
	}
	parent := map[NodeID]NodeID{}
	children := map[NodeID][]NodeID{}
	for u := range W {
		for _, e := range g.OutEdges(u) {
			v := e.V
			if _, ok := W[v]; ok && dhat[v] == dhat[u]+e.W {
				if _, exists := parent[v]; !exists {
					parent[v] = u
					children[u] = append(children[u], v)
				}
			}
		}
	}
	treeSize := map[NodeID]int{}
	var dfs func(NodeID) int
	visited := map[NodeID]bool{}
	dfs = func(u NodeID) int {
		if visited[u] {
			return treeSize[u]
		}
		visited[u] = true
		size := 1
		for _, c := range children[u] {
			size += dfs(c)
		}
		treeSize[u] = size
		return size
	}
	for v := range W {
		root := v
		for {
			if p, ok := parent[root]; ok {
				root = p
			} else {
				break
			}
		}
		dfs(root)
	}
	P := NewNodeSet()
	for v := range W {
		if _, hasParent := parent[v]; !hasParent && treeSize[v] >= k {
			P.Add(v)
		}
	}
	return P, W
}

// BMSSP — exact paper algorithm
func BMSSP(l int, B Dist, S NodeSet, k, t int, g *Graph, dhat map[NodeID]Dist) (Dist, NodeSet) {
	if l == 0 {
		// When l=0, use BaseCase for each source in S
		// The paper assumes S is a singleton for BaseCase
		return BaseCase(B, S, k, g, dhat)
	}
	P, W := FindPivots(B, S, k, g, dhat)
	M := 1 << ((l - 1) * t)
	D := NewDStruct()
	D.Initialize(M, B)
	for x := range P {
		D.Insert(x, dhat[x])
	}

	U := NewNodeSet()
	B0p := B
	if P.Len() > 0 {
		minp := INF
		for x := range P {
			if dhat[x] < minp {
				minp = dhat[x]
			}
		}
		B0p = minp
	}

	limit := 1
	for i := 0; i < 2*l; i++ {
		limit *= k
	}

	for U.Len() < limit && D.NonEmpty() {
		Bi, Si, _ := D.Pull()
		Bp, Ui := BMSSP(l-1, Bi, Si, k, t, g, dhat)
		for u := range Ui {
			U.Add(u)
		}
		K := []pair{}
		for u := range Ui {
			for _, e := range g.OutEdges(u) {
				v, w := e.V, e.W
				newd := dhat[u] + w
				if newd <= dhat[v] {
					dhat[v] = newd
					if newd >= Bi && newd < B {
						D.Insert(v, newd)
					} else if newd >= Bp && newd < Bi {
						K = append(K, pair{v, newd})
					}
				}
			}
		}
		for x := range Si {
			if dhat[x] >= Bp && dhat[x] < Bi {
				K = append(K, pair{x, dhat[x]})
			}
		}
		if len(K) > 0 {
			D.BatchPrepend(K)
		}
		if Bp < B0p {
			B0p = Bp
		}
	}
	Bprime := B0p
	if B < Bprime {
		Bprime = B
	}
	for x := range W {
		if dhat[x] < Bprime {
			U.Add(x)
		}
	}
	return Bprime, U
}
