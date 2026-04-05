package solver

import (
	"slices"

	"github.com/kpitt/sudoku/internal/bitset"
)

func (s *Solver) findNakedSubsets(size int, kind techniqueKind) (found bool) {
	return slices.ContainsFunc(s.houses[:], func(h *House) bool {
		return s.checkNakedSubsetsForHouse(size, kind, h)
	})
}

func (s *Solver) checkNakedSubsetsForHouse(size int, kind techniqueKind, h *House) (found bool) {
	var locsBuf [9]int
	locs := locsBuf[:0]
	// Collect a list of all locations with no more than `size` candidate values.
	for i, c := range h.Cells {
		numCand := c.NumCandidates()
		if numCand >= 2 && numCand <= size {
			locs = append(locs, i)
		}
	}
	if len(locs) < size {
		// We need at least `size` candidate values to have a subset of that size.
		return false
	}

	// Try combinations of the required size.
	var checkCombinations func(start int, currentIndices []int) bool
	checkCombinations = func(start int, currentIndices []int) bool {
		if len(currentIndices) == size {
			valueSet := bitset.BitSet16(0)
			for _, idx := range currentIndices {
				valueSet = bitset.Union(valueSet, h.Cells[locs[idx]].Candidates)
			}

			if valueSet.Size() == size {
				var locSet bitset.BitSet16
				for _, idx := range currentIndices {
					locSet.Add(locs[idx])
				}

				step := NewStep(kind)
				if s.eliminateFromOtherLocs(h, valueSet, locSet, step) {
					s.applyStep(step.
						WithIndices(h.indexesFromLocs(locSet.Values())...).
						WithValues(valueSet.Values()...).
						WithHouse(h))
					return true
				}
			}
			return false
		}

		for i := start; i < len(locs); i++ {
			if checkCombinations(i+1, append(currentIndices, i)) {
				return true
			}
		}
		return false
	}

	// Avoid allocation for small known sizes by allocating the initial buffer on the stack.
	indices := make([]int, 0, size)
	return checkCombinations(0, indices)
}

func (s *Solver) findNakedPairs() (found bool) {
	return s.findNakedSubsets(2, kindNakedPair)
}

func (s *Solver) findHiddenSubsets(size int, kind techniqueKind) (found bool) {
	return slices.ContainsFunc(s.houses[:], func(h *House) bool {
		return s.checkHiddenSubsetsForHouse(size, kind, h)
	})
}

func (s *Solver) checkHiddenSubsetsForHouse(size int, kind techniqueKind, h *House) (found bool) {
	var valBuf [9]int
	values := valBuf[:0]
	for val := 1; val <= 9; val++ {
		locs := h.Unsolved[val]
		if locs.Empty() {
			continue
		}
		if locs.Size() >= 2 && locs.Size() <= size {
			values = append(values, val)
		}
	}

	if len(values) < size {
		// We need at least `size` candidate values to have a subset of that size.
		return false
	}

	var checkCombinations func(start int, currentIndices []int) bool
	checkCombinations = func(start int, currentIndices []int) bool {
		if len(currentIndices) == size {
			locSet := bitset.BitSet16(0)
			valueSet := bitset.BitSet16(0)
			for _, idx := range currentIndices {
				val := values[idx]
				locSet = bitset.Union(locSet, h.Unsolved[val])
				valueSet.Add(val)
			}

			if locSet.Size() == size {
				step := NewStep(kind)
				if s.eliminateOtherValues(h, valueSet, locSet, step) {
					s.applyStep(step.
						WithIndices(h.indexesFromLocs(locSet.Values())...).
						WithValues(valueSet.Values()...).
						WithHouse(h))
					return true
				}
			}
			return false
		}

		for i := start; i < len(values); i++ {
			if checkCombinations(i+1, append(currentIndices, i)) {
				return true
			}
		}
		return false
	}

	indices := make([]int, 0, size)
	return checkCombinations(0, indices)
}

func (s *Solver) findHiddenPairs() (found bool) {
	return s.findHiddenSubsets(2, kindHiddenPair)
}

func (s *Solver) findNakedTriples() (found bool) {
	return s.findNakedSubsets(3, kindNakedTriple)
}

func (s *Solver) findHiddenTriples() (found bool) {
	return s.findHiddenSubsets(3, kindHiddenTriple)
}

func (s *Solver) findNakedQuadruples() (found bool) {
	return s.findNakedSubsets(4, kindNakedQuadruple)
}

func (s *Solver) findHiddenQuadruples() (found bool) {
	return s.findHiddenSubsets(4, kindHiddenQuadruple)
}
