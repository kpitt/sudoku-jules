package solver

// aicNode represents a literal in the implication graph.
// It is encoded as (literalIndex * 2) + (isTrue ? 1 : 0).
type aicNode int

func (n aicNode) negate() aicNode { return aicNode(int(n) ^ 1) }

const (
	numCells         = 81
	numDigits        = 9
	numIntersections = 54 // 27 row-box, 27 col-box
	numCellLits      = numCells * numDigits
	numIntLits       = numIntersections * numDigits
)

// Literal mapping functions
func cellLit(cellIdx, val int) int {
	return cellIdx*9 + (val - 1)
}

func rbLit(rIdx, bIdxInRow, val int) int {
	return numCellLits + (rIdx*3+bIdxInRow)*9 + (val - 1)
}

func cbLit(cIdx, bIdxInCol, val int) int {
	return numCellLits + 27*9 + (cIdx*3+bIdxInCol)*9 + (val - 1)
}

func alsLit(alsIdx, val, numBaseLits int) int {
	return numBaseLits + alsIdx*9 + (val - 1)
}

func node(lit int, isTrue bool) aicNode {
	if isTrue {
		return aicNode(lit*2 + 1)
	}
	return aicNode(lit*2 + 0)
}

type graphBuilder struct {
	s        *Solver
	adj      [][]aicNode
	numLits  int
	baseLits int
}

func newGraphBuilder(s *Solver, numALSLits int) *graphBuilder {
	baseLits := numCellLits + numIntLits
	totalLits := baseLits + numALSLits
	return &graphBuilder{
		s:        s,
		adj:      make([][]aicNode, totalLits*2),
		numLits:  totalLits,
		baseLits: baseLits,
	}
}

func (gb *graphBuilder) addImplication(u, v aicNode) {
	if u == v {
		return
	}
	for _, existing := range gb.adj[u] {
		if existing == v {
			return
		}
	}
	gb.adj[u] = append(gb.adj[u], v)
}

func (gb *graphBuilder) addStrongLink(u, v aicNode) {
	gb.addImplication(u.negate(), v)
	gb.addImplication(v.negate(), u)
	gb.addWeakLink(u, v)
}

func (gb *graphBuilder) addWeakLink(u, v aicNode) {
	gb.addImplication(u, v.negate())
	gb.addImplication(v, u.negate())
}

func (gb *graphBuilder) build(singleDigit int, options chainOptions, alsList []als) {
	// 1. Cell and Grouped links
	if options.useGroupedLinks {
		for c := 0; c < 81; c++ {
			cell := gb.s.puzzle.Cell(c)
			if cell.IsSolved() {
				continue
			}
			r, col := c/9, c%9
			for v := 1; v <= 9; v++ {
				if singleDigit != 0 && v != singleDigit {
					continue
				}
				if !cell.HasCandidate(v) {
					continue
				}

				cN := node(cellLit(c, v), true)
				rbN := node(rbLit(r, col/3, v), true)
				cbN := node(cbLit(col, r/3, v), true)

				gb.addImplication(cN, rbN)
				gb.addImplication(cN, cbN)
				gb.addImplication(rbN.negate(), cN.negate())
				gb.addImplication(cbN.negate(), cN.negate())
			}
		}
	}

	// 2. House links (Strong and Weak)
	for v := 1; v <= 9; v++ {
		if singleDigit != 0 && v != singleDigit {
			continue
		}

		for _, h := range gb.s.houses {
			if h.Unsolved[v].Empty() {
				continue
			}

			if options.useGroupedLinks && h.Kind != kindBox {
				var activeInts []int
				for i := 0; i < 3; i++ {
					var intLit int
					hasCandidate := false
					if h.Kind == kindRow {
						intLit = rbLit(h.Index, i, v)
						for k := 0; k < 3; k++ {
							if gb.s.puzzle.Cell(h.Index*9 + i*3 + k).HasCandidate(v) {
								hasCandidate = true
								break
							}
						}
					} else { // kindColumn
						intLit = cbLit(h.Index, i, v)
						for k := 0; k < 3; k++ {
							if gb.s.puzzle.Cell(k*9*3 + i*9 + h.Index).HasCandidate(v) {
								hasCandidate = true
								break
							}
						}
					}
					if hasCandidate {
						activeInts = append(activeInts, intLit)
					}
				}

				if options.useStrongLinks && len(activeInts) == 2 {
					gb.addStrongLink(node(activeInts[0], true), node(activeInts[1], true))
				} else if options.useWeakLinks && len(activeInts) >= 2 {
					for i := 0; i < len(activeInts); i++ {
						for j := i + 1; j < len(activeInts); j++ {
							gb.addWeakLink(node(activeInts[i], true), node(activeInts[j], true))
						}
					}
				}
			}

			locs := h.Unsolved[v].Values()
			if options.useStrongLinks && len(locs) == 2 {
				gb.addStrongLink(node(cellLit(h.Cells[locs[0]].Index(), v), true), node(cellLit(h.Cells[locs[1]].Index(), v), true))
			} else if options.useWeakLinks && len(locs) >= 2 {
				for i := 0; i < len(locs); i++ {
					for j := i + 1; j < len(locs); j++ {
						gb.addWeakLink(node(cellLit(h.Cells[locs[i]].Index(), v), true), node(cellLit(h.Cells[locs[j]].Index(), v), true))
					}
				}
			}
		}
	}

	// 3. Bivalue cells
	if options.useBivalueCells {
		for c := 0; c < 81; c++ {
			cell := gb.s.puzzle.Cell(c)
			if cell.NumCandidates() == 2 {
				vals := cell.CandidateValues()
				if singleDigit == 0 || (vals[0] == singleDigit || vals[1] == singleDigit) {
					gb.addStrongLink(node(cellLit(c, vals[0]), true), node(cellLit(c, vals[1]), true))
				}
			}
		}
	}

	// 4. Weak links in cell (cannot have two digits in same cell)
	if options.useWeakLinks && singleDigit == 0 {
		for c := 0; c < 81; c++ {
			cell := gb.s.puzzle.Cell(c)
			vals := cell.CandidateValues()
			if len(vals) > 1 {
				for i := 0; i < len(vals); i++ {
					for j := i + 1; j < len(vals); j++ {
						gb.addWeakLink(node(cellLit(c, vals[i]), true), node(cellLit(c, vals[j]), true))
					}
				}
			}
		}
	}

	// 5. Grouped implications (Intersection -> Cell)
	if options.useGroupedLinks {
		for v := 1; v <= 9; v++ {
			if singleDigit != 0 && v != singleDigit {
				continue
			}
			for r := 0; r < 9; r++ {
				for i := 0; i < 3; i++ {
					rbN := node(rbLit(r, i, v), true)
					for cIdx := 0; cIdx < 9; cIdx++ {
						if cIdx/3 == i { continue }
						cellIdx := r*9 + cIdx
						if gb.s.puzzle.Cell(cellIdx).HasCandidate(v) {
							gb.addImplication(rbN, node(cellLit(cellIdx, v), false))
						}
					}
					boxIdx := (r/3)*3 + i
					box := gb.s.boxes[boxIdx]
					for _, cell := range box.Cells {
						if cell.Row == r { continue }
						if cell.HasCandidate(v) {
							gb.addImplication(rbN, node(cellLit(cell.Index(), v), false))
						}
					}
				}
			}
			for c := 0; c < 9; c++ {
				for i := 0; i < 3; i++ {
					cbN := node(cbLit(c, i, v), true)
					for rIdx := 0; rIdx < 9; rIdx++ {
						if rIdx/3 == i { continue }
						cellIdx := rIdx*9 + c
						if gb.s.puzzle.Cell(cellIdx).HasCandidate(v) {
							gb.addImplication(cbN, node(cellLit(cellIdx, v), false))
						}
					}
					boxIdx := (c/3) + i*3
					box := gb.s.boxes[boxIdx]
					for _, cell := range box.Cells {
						if cell.Col == c { continue }
						if cell.HasCandidate(v) {
							gb.addImplication(cbN, node(cellLit(cell.Index(), v), false))
						}
					}
				}
			}
		}
	}

	// 6. ALS links
	if len(alsList) > 0 {
		for i, a := range alsList {
			vals := a.candidates.Values()
			// Internal strong links: not v1 => v2, not v1 => v3, etc.
			for _, v1 := range vals {
				n1 := node(alsLit(i, v1, gb.baseLits), true)
				for _, v2 := range vals {
					if v1 == v2 { continue }
					n2 := node(alsLit(i, v2, gb.baseLits), true)
					gb.addImplication(n1.negate(), n2)
				}
				// Link ALS literals to cell literals
				// Cell(c, v) is TRUE => ALS(i, v) is TRUE
				for _, cIdx := range a.cells {
					if gb.s.puzzle.Cell(cIdx).HasCandidate(v1) {
						gb.addImplication(node(cellLit(cIdx, v1), true), n1)
						gb.addImplication(n1.negate(), node(cellLit(cIdx, v1), false))
					}
				}
			}
		}

		// Inter-ALS weak links (Restricted Commons)
		for i := 0; i < len(alsList); i++ {
			for j := i + 1; j < len(alsList); j++ {
				a, b := alsList[i], alsList[j]
				common := a.candidates.Intersection(b.candidates)
				for _, v := range common.Values() {
					if gb.s.isRestrictedCommon(a, b, v) {
						gb.addWeakLink(node(alsLit(i, v, gb.baseLits), true), node(alsLit(j, v, gb.baseLits), true))
					}
				}
			}
		}
		
		// Interaction between ALS and cells
		for i, a := range alsList {
			for _, v := range a.candidates.Values() {
				aLitN := node(alsLit(i, v, gb.baseLits), true)
				// If ALS(i, v) is TRUE, then any cell that sees ALL v-cells in ALS A cannot be v.
				vCells := gb.s.cellsWithCandidate(a.cells, v)
				for cIdx := 0; cIdx < 81; cIdx++ {
					if !gb.s.puzzle.Cell(cIdx).HasCandidate(v) { continue }
					inALS := false
					for _, ac := range a.cells { if ac == cIdx { inALS = true; break } }
					if inALS { continue }

					seesAll := true
					for _, vc := range vCells {
						if !gb.s.sees(cIdx, vc) {
							seesAll = false
							break
						}
					}
					if seesAll {
						gb.addImplication(aLitN, node(cellLit(cIdx, v), false))
						gb.addImplication(node(cellLit(cIdx, v), true), aLitN.negate())
					}
				}
			}
		}
	}
}

type chainOptions struct {
	useStrongLinks  bool
	useWeakLinks    bool
	useBivalueCells bool
	useGroupedLinks bool
}

func (s *Solver) findAICs() bool {
	return s.findChains(0, kindAIC, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	}, nil)
}

func (s *Solver) findNiceLoops() bool {
	return s.findChains(0, kindNiceLoop, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	}, nil)
}

func (s *Solver) findMultiColoring() bool {
	found := false
	for v := 1; v <= 9; v++ {
		if s.findChains(v, kindMultiColoring, chainOptions{
			useStrongLinks:  true,
			useWeakLinks:    true,
			useBivalueCells: false,
			useGroupedLinks: true,
		}, nil) {
			found = true
		}
	}
	return found
}

func (s *Solver) findXChains() bool {
	found := false
	for v := 1; v <= 9; v++ {
		if s.findChains(v, kindXChain, chainOptions{
			useStrongLinks:  true,
			useWeakLinks:    true,
			useBivalueCells: false,
			useGroupedLinks: true,
		}, nil) {
			found = true
		}
	}
	return found
}

func (s *Solver) findXYChains() bool {
	return s.findChains(0, kindXYChain, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: false,
	}, nil)
}

func (s *Solver) findALSXYChain() bool {
	alsList := s.findALSs()
	return s.findChains(0, kindALSXYChain, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	}, alsList)
}

func (s *Solver) findForcingChains() bool {
	alsList := s.findALSs()
	return s.findChains(0, kindForcingChain, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	}, alsList)
}

func (s *Solver) findForcingNets() bool {
	alsList := s.findALSs()
	numALSLits := len(alsList) * 9
	gb := newGraphBuilder(s, numALSLits)
	gb.build(0, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	}, alsList)

	// 1. Cell Forcing Nets: For each cell, if all its candidates lead to the same implication.
	for c := 0; c < 81; c++ {
		cell := s.puzzle.Cell(c)
		if cell.IsSolved() || cell.NumCandidates() < 2 {
			continue
		}

		candidates := cell.CandidateValues()
		reachableSets := make([]map[aicNode]bool, len(candidates))
		for i, v := range candidates {
			reachableSets[i] = s.getReachableNodes(gb.adj, node(cellLit(c, v), true))
		}

		if s.checkForcingNetImplications(kindForcingNet, reachableSets, []int{c}, intersectionToStep) {
			return true
		}
	}

	// 2. House Forcing Nets: For each house and digit, if all possible locations lead to the same implication.
	for _, h := range s.houses {
		for v := 1; v <= 9; v++ {
			locs := h.Unsolved[v].Values()
			if len(locs) < 2 {
				continue
			}

			reachableSets := make([]map[aicNode]bool, len(locs))
			baseIndices := make([]int, len(locs))
			for i, loc := range locs {
				cIdx := h.Cells[loc].Index()
				baseIndices[i] = cIdx
				reachableSets[i] = s.getReachableNodes(gb.adj, node(cellLit(cIdx, v), true))
			}

			if s.checkForcingNetImplications(kindForcingNet, reachableSets, baseIndices, intersectionToStep) {
				return true
			}
		}
	}

	return false
}

type forcingNetHandler func(s *Solver, kind techniqueKind, intersection map[aicNode]bool, baseIndices []int) bool

func (s *Solver) checkForcingNetImplications(
	kind techniqueKind,
	reachableSets []map[aicNode]bool,
	baseIndices []int,
	handler forcingNetHandler,
) bool {
	if len(reachableSets) == 0 {
		return false
	}

	// Find intersection of all reachableSets
	intersection := make(map[aicNode]bool)
	for n := range reachableSets[0] {
		inAll := true
		for i := 1; i < len(reachableSets); i++ {
			if !reachableSets[i][n] {
				inAll = false
				break
			}
		}
		if inAll {
			intersection[n] = true
		}
	}

	return handler(s, kind, intersection, baseIndices)
}

func intersectionToStep(s *Solver, kind techniqueKind, intersection map[aicNode]bool, baseIndices []int) bool {
	step := NewStep(kind).WithIndices(baseIndices...)
	found := false
	for n := range intersection {
		lit := int(n) / 2
		isTrue := int(n)%2 == 1
		if lit < numCellLits {
			tcIdx, tvIdx := lit/9, lit%9+1
			tCell := s.puzzle.Cell(tcIdx)
			if tCell.IsSolved() {
				continue
			}

			// Avoid trivial implications (the base itself)
			isBase := false
			for _, bIdx := range baseIndices {
				if bIdx == tcIdx {
					isBase = true
					break
				}
			}
			if isBase {
				continue
			}

			if isTrue {
				step.PlaceCandidate(tcIdx, tvIdx)
				found = true
			} else {
				if tCell.HasCandidate(tvIdx) {
					step.DeleteCandidate(tcIdx, tvIdx)
					found = true
				}
			}
		}
	}
	if found {
		s.applyStep(step)
		return true
	}
	return false
}

func (s *Solver) getReachableNodes(adj [][]aicNode, start aicNode) map[aicNode]bool {
	reachable := make(map[aicNode]bool)
	queue := []aicNode{start}
	reachable[start] = true
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, next := range adj[curr] {
			if !reachable[next] {
				reachable[next] = true
				queue = append(queue, next)
			}
		}
	}
	return reachable
}

func (s *Solver) findChains(singleDigit int, kind techniqueKind, options chainOptions, alsList []als) bool {
	numALSLits := len(alsList) * 9
	gb := newGraphBuilder(s, numALSLits)
	gb.build(singleDigit, options, alsList)

	step := NewStep(kind)

	for lit := 0; lit < numCellLits; lit++ {
		cIdx, vIdx := lit/9, lit%9+1
		if !s.puzzle.Cell(cIdx).HasCandidate(vIdx) {
			continue
		}
		if singleDigit != 0 && vIdx != singleDigit {
			continue
		}

		if path := s.getContradictionPath(gb.adj, node(lit, true)); path != nil {
			step.DeleteCandidate(cIdx, vIdx)
			step.WithIndices(extractIndices(path)...)
			if singleDigit != 0 {
				step.WithValues(singleDigit)
			}
			s.applyStep(step)
			return true
		}
		if path := s.getContradictionPath(gb.adj, node(lit, false)); path != nil {
			step.PlaceCandidate(cIdx, vIdx)
			step.WithIndices(extractIndices(path)...)
			if singleDigit != 0 {
				step.WithValues(singleDigit)
			}
			s.applyStep(step)
			return true
		}
	}

	return false
}

func extractIndices(path []aicNode) []int {
	var indices []int
	seen := make(map[int]bool)
	for _, n := range path {
		lit := int(n) / 2
		if lit < numCellLits {
			idx := lit / 9
			if !seen[idx] {
				indices = append(indices, idx)
				seen[idx] = true
			}
		}
	}
	return indices
}

func (s *Solver) getContradictionPath(adj [][]aicNode, start aicNode) []aicNode {
	totalNodes := len(adj)
	visited := make([]aicNode, totalNodes)
	for i := range visited {
		visited[i] = -1
	}
	queue := []aicNode{start}
	visited[start] = start // self as marker
	target := start.negate()

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == target {
			found = true
			break
		}
		for _, next := range adj[curr] {
			if visited[next] == -1 {
				visited[next] = curr
				queue = append(queue, next)
			}
		}
	}

	if found {
		var path []aicNode
		curr := target
		for curr != start {
			path = append(path, curr)
			curr = visited[curr]
		}
		path = append(path, start)
		// Reverse path
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		return path
	}
	return nil
}
