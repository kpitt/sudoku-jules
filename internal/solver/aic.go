package solver

type aicNode int

func (n aicNode) literal() int { return int(n / 2) }
func (n aicNode) isTrue() bool { return n%2 == 1 }

func node(cell, val int, isTrue bool) aicNode {
	lit := cell*9 + (val - 1)
	if isTrue {
		return aicNode(lit*2 + 1)
	}
	return aicNode(lit*2 + 0)
}

func (s *Solver) findAICs() bool {
	return s.findAICsInternal(kindAIC)
}

func (s *Solver) findNiceLoops() bool {
	return s.findAICsInternal(kindNiceLoop)
}

func (s *Solver) findAICsInternal(kind techniqueKind) bool {
	numLiterals := 81 * 9
	numNodes := numLiterals * 2
	adj := make([][]aicNode, numNodes)

	addImplication := func(u, v aicNode) {
		// Avoid duplicates
		for _, existing := range adj[u] {
			if existing == v {
				return
			}
		}
		adj[u] = append(adj[u], v)
	}

	addStrongLink := func(c1, v1, c2, v2 int) {
		// !L1 => L2  and  !L2 => L1
		n1F := node(c1, v1, false)
		n1T := node(c1, v1, true)
		n2F := node(c2, v2, false)
		n2T := node(c2, v2, true)

		addImplication(n1F, n2T)
		addImplication(n2F, n1T)

		// Strong link is also a weak link:
		// L1 => !L2 and L2 => !L1
		addImplication(n1T, n2F)
		addImplication(n2T, n1F)
	}

	addWeakLink := func(c1, v1, c2, v2 int) {
		// L1 => !L2 and L2 => !L1
		n1T := node(c1, v1, true)
		n1F := node(c1, v1, false)
		n2T := node(c2, v2, true)
		n2F := node(c2, v2, false)

		addImplication(n1T, n2F)
		addImplication(n2T, n1F)
	}

	// 1. Build graph
	// Strong links from houses
	for val := 1; val <= 9; val++ {
		for _, h := range s.houses {
			if h.NumLocations(val) == 2 {
				locs := h.Locations(val).Values()
				c1 := h.Cells[locs[0]].Index()
				c2 := h.Cells[locs[1]].Index()
				addStrongLink(c1, val, c2, val)
			}
		}
	}
	// Strong links from cells (bivalue cells)
	for c := 0; c < 81; c++ {
		cell := s.puzzle.Cell(c)
		if cell.NumCandidates() == 2 {
			vals := cell.CandidateValues()
			addStrongLink(c, vals[0], c, vals[1])
		}
	}

	// Weak links from houses
	for val := 1; val <= 9; val++ {
		for _, h := range s.houses {
			locs := h.Locations(val).Values()
			if len(locs) < 2 {
				continue
			}
			for i := 0; i < len(locs); i++ {
				for j := i + 1; j < len(locs); j++ {
					c1 := h.Cells[locs[i]].Index()
					c2 := h.Cells[locs[j]].Index()
					addWeakLink(c1, val, c2, val)
				}
			}
		}
	}
	// Weak links from cells
	for c := 0; c < 81; c++ {
		cell := s.puzzle.Cell(c)
		vals := cell.CandidateValues()
		if len(vals) < 2 {
			continue
		}
		for i := 0; i < len(vals); i++ {
			for j := i + 1; j < len(vals); j++ {
				addWeakLink(c, vals[i], c, vals[j])
			}
		}
	}

	// 2. Compute reachability and check for eliminations
	// We'll use a BFS for each node on demand to save memory, or precompute if needed.
	// Since we want to find THE first elimination, we can iterate over potential deductions.

	found := false
	// Case 1: Discontinuous Nice Loop / Forcing Chain (False => True)
	for i := 0; i < numLiterals; i++ {
		c, v := i/9, i%9+1
		if !s.puzzle.Cell(c).HasCandidate(v) {
			continue
		}

		// If False(c, v) => True(c, v), then (c, v) MUST be True.
		if s.isReachable(adj, node(c, v, false), node(c, v, true)) {
			step := NewStep(kind).WithPlacedValue(c, v)
			s.applyStep(step)
			found = true
			continue // Already handled this literal
		}

		// If True(c, v) => False(c, v), then (c, v) MUST be False.
		if s.isReachable(adj, node(c, v, true), node(c, v, false)) {
			step := NewStep(kind).WithValues(v)
			step.DeleteCandidate(c, v)
			s.applyStep(step)
			found = true
		}
	}

	// Case 2: AIC / Nice Loop (L1 OR L2)
	// For each pair of literals L1, L2 such that !L1 => L2
	for l1 := 0; l1 < numLiterals; l1++ {
		c1, v1 := l1/9, l1%9+1
		if !s.puzzle.Cell(c1).HasCandidate(v1) {
			continue
		}

		reachableFromL1F := s.getReachable(adj, node(c1, v1, false))

		for l2 := l1 + 1; l2 < numLiterals; l2++ {
			c2, v2 := l2/9, l2%9+1
			if !s.puzzle.Cell(c2).HasCandidate(v2) {
				continue
			}

			if reachableFromL1F[node(c2, v2, true)] {
				// We have L1 OR L2.
				// Find X such that L1 => !X AND L2 => !X.
				reachableFromL1T := s.getReachable(adj, node(c1, v1, true))
				reachableFromL2T := s.getReachable(adj, node(c2, v2, true))

				step := NewStep(kind)
				for x := 0; x < numLiterals; x++ {
					cx, vx := x/9, x%9+1
					if !s.puzzle.Cell(cx).HasCandidate(vx) {
						continue
					}

					// Check if True(L1) => False(X) AND True(L2) => False(X)
					if reachableFromL1T[node(cx, vx, false)] && reachableFromL2T[node(cx, vx, false)] {
						step.DeleteCandidate(cx, vx)
					}
				}

				if len(step.deletedCandidates) > 0 {
					s.applyStep(step.WithValues(v1, v2).WithIndices(c1, c2))
					found = true
				}
			}
		}
	}

	return found
}

func (s *Solver) isReachable(adj [][]aicNode, start, end aicNode) bool {
	visited := make([]bool, len(adj))
	queue := []aicNode{start}
	visited[start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			return true
		}

		for _, next := range adj[curr] {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

func (s *Solver) getReachable(adj [][]aicNode, start aicNode) []bool {
	visited := make([]bool, len(adj))
	queue := []aicNode{start}
	visited[start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, next := range adj[curr] {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return visited
}
