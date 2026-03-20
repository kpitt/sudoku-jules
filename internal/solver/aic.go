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
	totalLits        = numCellLits + numIntLits
	totalNodes       = totalLits * 2
)

// Literal mapping functions
func cellLit(cellIdx, val int) int {
	return cellIdx*9 + (val - 1)
}

func rbLit(rIdx, bIdxInRow, val int) int {
	// rIdx: 0..8, bIdxInRow: 0..2 (box 0, 1, 2 in that row)
	return numCellLits + (rIdx*3+bIdxInRow)*9 + (val - 1)
}

func cbLit(cIdx, bIdxInCol, val int) int {
	// cIdx: 0..8, bIdxInCol: 0..2 (box 0, 1, 2 in that column)
	return numCellLits + 27*9 + (cIdx*3+bIdxInCol)*9 + (val - 1)
}

func node(lit int, isTrue bool) aicNode {
	if isTrue {
		return aicNode(lit*2 + 1)
	}
	return aicNode(lit*2 + 0)
}

type graphBuilder struct {
	s   *Solver
	adj [][]aicNode
}

func newGraphBuilder(s *Solver) *graphBuilder {
	return &graphBuilder{
		s:   s,
		adj: make([][]aicNode, totalNodes),
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

func (gb *graphBuilder) build(singleDigit int, options chainOptions) {
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
	})
}

func (s *Solver) findNiceLoops() bool {
	return s.findChains(0, kindNiceLoop, chainOptions{
		useStrongLinks:  true,
		useWeakLinks:    true,
		useBivalueCells: true,
		useGroupedLinks: true,
	})
}

func (s *Solver) findXChains() bool {
	found := false
	for v := 1; v <= 9; v++ {
		if s.findChains(v, kindXChain, chainOptions{
			useStrongLinks:  true,
			useWeakLinks:    true,
			useBivalueCells: false,
			useGroupedLinks: true,
		}) {
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
	})
}

func (s *Solver) findChains(singleDigit int, kind techniqueKind, options chainOptions) bool {
	gb := newGraphBuilder(s)
	gb.build(singleDigit, options)

	found := false
	step := NewStep(kind)

	for lit := 0; lit < numCellLits; lit++ {
		cIdx, vIdx := lit/9, lit%9+1
		if !s.puzzle.Cell(cIdx).HasCandidate(vIdx) {
			continue
		}
		if singleDigit != 0 && vIdx != singleDigit {
			continue
		}

		if s.canReachContradiction(gb.adj, node(lit, true)) {
			step.DeleteCandidate(cIdx, vIdx)
			found = true
		}
		if s.canReachContradiction(gb.adj, node(lit, false)) {
			step.PlaceCandidate(cIdx, vIdx)
			found = true
		}
	}

	if found {
		if singleDigit != 0 {
			step.WithValues(singleDigit)
		}
		s.applyStep(step)
	}
	return found
}

func (s *Solver) canReachContradiction(adj [][]aicNode, start aicNode) bool {
	visited := make([]bool, totalNodes)
	queue := []aicNode{start}
	visited[start] = true
	target := start.negate()

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == target {
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
