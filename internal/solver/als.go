package solver

import (
	"github.com/kpitt/sudoku/internal/bitset"
)

type als struct {
	cells      []int
	candidates bitset.BitSet16
	house      *House
}

func (s *Solver) findALSs() []als {
	var result []als
	for _, h := range s.houses {
		unsolved := h.UnsolvedCells()
		if len(unsolved) < 2 {
			continue
		}

		// Find subsets of size N from unsolved
		for n := 1; n <= 5 && n < len(unsolved); n++ {
			s.forEachSubset(unsolved, n, func(subset []int) {
				candidates := bitset.BitSet16(0)
				for _, cellIdx := range subset {
					candidates.Union(s.puzzle.Cell(cellIdx).Candidates)
				}
				if candidates.Size() == n+1 {
					result = append(result, als{
						cells:      subset,
						candidates: candidates,
						house:      h,
					})
				}
			})
		}
	}
	return result
}

func (s *Solver) forEachSubset(elements []int, n int, f func([]int)) {
	if n <= 0 || n > len(elements) {
		return
	}
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	for {
		subset := make([]int, n)
		for i, idx := range indices {
			subset[i] = elements[idx]
		}
		f(subset)

		// Find next combination
		i := n - 1
		for i >= 0 && indices[i] == i+len(elements)-n {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < n; j++ {
			indices[j] = indices[j-1] + 1
		}
	}
}

func (s *Solver) findALSXZ() bool {
	alsList := s.findALSs()
	if len(alsList) < 2 {
		return false
	}

	found := false
	for i := 0; i < len(alsList); i++ {
		for j := i + 1; j < len(alsList); j++ {
			a, b := alsList[i], alsList[j]

			// ALSs must be distinct (though they can share cells)
			if s.alsEqual(a, b) {
				continue
			}

			// Find common candidates
			common := a.candidates.Intersection(b.candidates)
			if common.Empty() {
				continue
			}

			commonVals := common.Values()
			// Restricted Common candidate x:
			// All x in A see all x in B.
			var restrictedCommons []int
			for _, x := range commonVals {
				if s.isRestrictedCommon(a, b, x) {
					restrictedCommons = append(restrictedCommons, x)
				}
			}

			if len(restrictedCommons) == 0 {
				continue
			}

			// Case 1: Singly linked ALS-XZ
			// One restricted common x, and another common candidate z.
			if len(restrictedCommons) == 1 {
				x := restrictedCommons[0]
				for _, z := range commonVals {
					if z == x {
						continue
					}

					// Target cells see all z in A AND all z in B.
					step := NewStep(kindALSXZ)
					targets := s.findALSXZTargets(a, b, z)
					for _, t := range targets {
						if s.puzzle.Cell(t).HasCandidate(z) {
							step.DeleteCandidate(t, z)
						}
					}

					if len(step.deletedCandidates) > 0 {
						s.applyStep(step.WithValues(x, z).WithIndices(append(a.cells, b.cells...)...))
						found = true
					}
				}
			}

			// Case 2: Doubly linked ALS-XZ
			// Two restricted commons x and y.
			// Then ALL other candidates common to A and B can be eliminated from
			// cells seeing all instances in both.
			// Wait, actually doubly linked is stronger:
			// Candidates in A that are not x or y and are in B...
			if len(restrictedCommons) >= 2 {
				// For any candidate z in A or B that is not one of the restricted commons:
				// If it's common to both, it can be eliminated from peers?
				// Actually, even stronger: any z in A seeing all z in B...
				// Standard Doubly Linked ALS-XZ:
				// If x and y are both restricted commons, then:
				// 1. Any other common candidate z can be eliminated from common peers.
				// 2. Any z restricted to A can be eliminated from its peers in A? No.
				
				for _, x := range restrictedCommons {
					for _, y := range restrictedCommons {
						if x >= y {
							continue
						}
						// x and y are restricted commons.
						for _, z := range commonVals {
							if z == x || z == y {
								continue
							}
							targets := s.findALSXZTargets(a, b, z)
							step := NewStep(kindALSXZ)
							for _, t := range targets {
								if s.puzzle.Cell(t).HasCandidate(z) {
									step.DeleteCandidate(t, z)
								}
							}
							if len(step.deletedCandidates) > 0 {
								s.applyStep(step.WithValues(x, y, z).WithIndices(append(a.cells, b.cells...)...))
								found = true
							}
						}
					}
				}
			}
		}
	}

	return found
}

func (s *Solver) alsEqual(a, b als) bool {
	if len(a.cells) != len(b.cells) {
		return false
	}
	for i := range a.cells {
		if a.cells[i] != b.cells[i] {
			return false
		}
	}
	return true
}

func (s *Solver) isRestrictedCommon(a, b als, x int) bool {
	// All x in A must see all x in B.
	aCellsWithX := s.cellsWithCandidate(a.cells, x)
	bCellsWithX := s.cellsWithCandidate(b.cells, x)
	
	for _, ac := range aCellsWithX {
		for _, bc := range bCellsWithX {
			if ac == bc { 
				// If they share a cell with candidate X, then X cannot be a 
				// restricted common candidate because that cell doesn't see itself.
				return false 
			}
			if !s.sees(ac, bc) {
				return false
			}
		}
	}
	return true
}

func (s *Solver) cellsWithCandidate(cells []int, val int) []int {
	var res []int
	for _, c := range cells {
		if s.puzzle.Cell(c).HasCandidate(val) {
			res = append(res, c)
		}
	}
	return res
}

func (s *Solver) findALSXZTargets(a, b als, z int) []int {
	aCellsWithZ := s.cellsWithCandidate(a.cells, z)
	bCellsWithZ := s.cellsWithCandidate(b.cells, z)

	if len(aCellsWithZ) == 0 || len(bCellsWithZ) == 0 {
		return nil
	}

	// Target cell must see all in aCellsWithZ AND all in bCellsWithZ.
	// And target cell must not be in A or B? 
	// Actually, target cell can be in A or B if it's not one of the z-cells?
	// But usually we just check all cells in the puzzle.
	
	var targets []int
	for i := 0; i < 81; i++ {
		// Skip cells in the ALSs? 
		// Actually, standard rule: if it's in A and sees all z in B, it's already accounted for.
		// But let's just check all cells for simplicity.
		inALS := false
		for _, c := range a.cells { if c == i { inALS = true; break } }
		if !inALS {
			for _, c := range b.cells { if c == i { inALS = true; break } }
		}
		if inALS {
			// A cell in ALS A that has z cannot be eliminated unless it also sees all z in A? 
			// No, a cell in A that HAS z is part of the ALS. It can't be a target.
			continue
		}

		seesAllA := true
		for _, ac := range aCellsWithZ {
			if !s.sees(i, ac) {
				seesAllA = false
				break
			}
		}
		if !seesAllA { continue }

		seesAllB := true
		for _, bc := range bCellsWithZ {
			if !s.sees(i, bc) {
				seesAllB = false
				break
			}
		}
		if seesAllB {
			targets = append(targets, i)
		}
	}
	return targets
}

func (s *Solver) findALSXYWing() bool {
	alsList := s.findALSs()
	if len(alsList) < 3 {
		return false
	}

	type connection struct {
		otherIdx int
		rc       int
	}
	connections := make([][]connection, len(alsList))
	for i := 0; i < len(alsList); i++ {
		for j := 0; j < len(alsList); j++ {
			if i == j {
				continue
			}
			common := alsList[i].candidates.Intersection(alsList[j].candidates)
			for _, x := range common.Values() {
				if s.isRestrictedCommon(alsList[i], alsList[j], x) {
					connections[i] = append(connections[i], connection{j, x})
				}
			}
		}
	}

	found := false
	for j := 0; j < len(alsList); j++ {
		b := alsList[j]
		conns := connections[j]
		if len(conns) < 2 {
			continue
		}

		for idxA := 0; idxA < len(conns); idxA++ {
			a := alsList[conns[idxA].otherIdx]
			x := conns[idxA].rc

			for idxC := idxA + 1; idxC < len(conns); idxC++ {
				c := alsList[conns[idxC].otherIdx]
				y := conns[idxC].rc
				if x == y {
					continue
				}

				// A and C must share candidate Z
				commonAC := a.candidates.Intersection(c.candidates)
				for _, z := range commonAC.Values() {
					if z == x || z == y {
						continue
					}

					// Target Z cannot be in B
					if b.candidates.Contains(z) {
						continue
					}

					targets := s.findALSXZTargets(a, c, z)
					step := NewStep(kindALSXYWing)
					for _, t := range targets {
						if s.puzzle.Cell(t).HasCandidate(z) {
							step.DeleteCandidate(t, z)
						}
					}

					if len(step.deletedCandidates) > 0 {
						s.applyStep(step.WithValues(x, y, z).WithIndices(append(append(a.cells, b.cells...), c.cells...)...))
						found = true
					}
				}
			}
		}
	}
	return found
}
