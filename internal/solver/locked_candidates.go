package solver

import (
	"github.com/kpitt/sudoku/internal/bitset"
	"github.com/kpitt/sudoku/internal/puzzle"
)

func (s *Solver) findClaimingTuples() (found bool) {
	for i := range 9 {
		// We only need to check rows and columns for Locked Candidates.
		if s.checkLockedCandidatesForLine(s.rows[i]) ||
			s.checkLockedCandidatesForLine(s.columns[i]) {

			return true
		}
	}
	return false
}

func (s *Solver) checkLockedCandidatesForLine(line *House) (found bool) {
	for val := 1; val <= 9; val++ {
		locs := line.Unsolved[val]
		if locs.Empty() {
			continue
		}

		// If we have more than 3 candidates in a line, then they can't all be
		// in the same box.
		if locs.Size() > 3 {
			continue
		}

		valueSet := bitset.FromValues16(val)
		if box, ok := line.sharedBox(locs); ok {
			cells := line.cellsFromLocs(locs.Values())
			boxCells := transformSlice(cells, func(c *puzzle.Cell) int {
				_, index := getBoxLoc(c.Row, c.Col)
				return index
			})
			locSet := bitset.FromValues16(boxCells...)
			step := NewStep(kindLockedCandidatesClaiming)
			if s.eliminateFromOtherLocs(s.boxes[box], valueSet, locSet, step) {
				s.applyStep(step.
					WithValues(val).
					WithHouse(line))
				return true
			}
		}
	}

	return false
}

func (s *Solver) findPointingTuples() (found bool) {
	for i := range 9 {
		// We only need to check boxes for Pointing Tuples.
		if s.checkPointingTuplesForBox(s.boxes[i]) {
			return true
		}
	}
	return false
}

func (s *Solver) checkPointingTuplesForBox(box *House) (found bool) {
	for val := 1; val <= 9; val++ {
		locs := box.Unsolved[val]
		if locs.Empty() {
			continue
		}

		// If we have more than 3 candidates in a single box, then they can't all
		// be in the same line.
		if locs.Size() > 3 {
			continue
		}

		step := NewStep(kindLockedCandidatesPointing).
			WithValues(val).
			WithHouse(box)
		valueSet := bitset.FromValues16(val)
		cells := box.cellsFromLocs(locs.Values())
		if row, ok := box.sharedRow(locs); ok {
			cols := transformSlice(cells, func(c *puzzle.Cell) int {
				return c.Col
			})
			locSet := bitset.FromValues16(cols...)
			if s.eliminateFromOtherLocs(s.rows[row], valueSet, locSet, step) {
				s.applyStep(step)
				return true
			}
		}
		if col, ok := box.sharedCol(locs); ok {
			rows := transformSlice(cells, func(c *puzzle.Cell) int {
				return c.Row
			})
			locSet := bitset.FromValues16(rows...)
			if s.eliminateFromOtherLocs(s.columns[col], valueSet, locSet, step) {
				s.applyStep(step)
				return true
			}
		}
	}

	return false
}
