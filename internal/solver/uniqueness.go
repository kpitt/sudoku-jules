package solver

import (
	"github.com/kpitt/sudoku/internal/puzzle"
)

func (s *Solver) findAvoidableRectangles() (found bool) {
	// We iterate through all possible pairs of rows and columns to form rectangles.
	// Since we need exactly two boxes, we must check the box condition.
	for r1 := range 8 {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := range 8 {
				for c2 := c1 + 1; c2 < 9; c2++ {
					// The four cells of the rectangle
					idx11 := r1*9 + c1
					idx12 := r1*9 + c2
					idx21 := r2*9 + c1
					idx22 := r2*9 + c2

					cell11 := s.puzzle.Cell(idx11)
					cell12 := s.puzzle.Cell(idx12)
					cell21 := s.puzzle.Cell(idx21)
					cell22 := s.puzzle.Cell(idx22)

					// Must span exactly two boxes
					// They form a rectangle in exactly two boxes if and only if
					// (r1/3 == r2/3 && c1/3 != c2/3) OR (r1/3 != r2/3 && c1/3 == c2/3)
					if (r1/3 == r2/3) == (c1/3 == c2/3) {
						continue
					}

					// Check how many cells are solved
					solvedCount := 0
					if cell11.IsSolved() {
						solvedCount++
					}
					if cell12.IsSolved() {
						solvedCount++
					}
					if cell21.IsSolved() {
						solvedCount++
					}
					if cell22.IsSolved() {
						solvedCount++
					}

					// If 0, 1, or 4 cells are solved, Avoidable Rectangle does not apply here.
					if solvedCount < 2 || solvedCount == 4 {
						continue
					}

					// None of the cells can be a given
					if cell11.IsGiven || cell12.IsGiven || cell21.IsGiven || cell22.IsGiven {
						continue
					}

					// For Type 1 Avoidable Rectangle, exactly 3 cells are solved.
					if solvedCount == 3 {
						found = s.checkAvoidableRectangleType1(cell11, cell12, cell21, cell22) || found
					} else if solvedCount == 2 {
						// For Type 2 Avoidable Rectangle, exactly 2 cells are solved.
						found = s.checkAvoidableRectangleType2(cell11, cell12, cell21, cell22) || found
					}
				}
			}
		}
	}
	return found
}

func (s *Solver) checkAvoidableRectangleType1(c11, c12, c21, c22 *puzzle.Cell) bool {
	var solved [3]*puzzle.Cell
	var unsolved *puzzle.Cell

	cells := []*puzzle.Cell{c11, c12, c21, c22}
	idx := 0
	for _, c := range cells {
		if c.IsSolved() {
			solved[idx] = c
			idx++
		} else {
			unsolved = c
		}
	}

	// Solved cells must have exactly 2 distinct values, structured as X, X, Y
	v1 := solved[0].Value()
	v2 := solved[1].Value()
	v3 := solved[2].Value()

	if v1 != v2 && v1 != v3 && v2 != v3 {
		return false // Values must form an X, X, Y pattern
	}

	// For an avoidable rectangle, the two identical values (X) must be diagonally opposite.
	// We can check if the unsolved cell is diagonally opposite to the Y value.
	// If it is, then the Y value's opposite is unsolved.
	// The diagonally opposite pairs are (c11, c22) and (c12, c21).
	var diag1 *puzzle.Cell
	if unsolved == c11 || unsolved == c22 {
		if unsolved == c11 {
			diag1 = c22
		} else {
			diag1 = c11
		}
	} else {
		if unsolved == c12 {
			diag1 = c21
		} else {
			diag1 = c12
		}
	}

	// The unsolved cell's diagonal counterpart (diag1) must be Y.
	// And the other two cells (forming the other diagonal) must be X.
	// Wait, if unsolved is X, then the deadly pattern is X, Y / Y, X.
	// Let the other two cells (the fully solved diagonal) have value X.
	// Then diag1 must have value Y.
	// Then unsolved cannot be Y.

	// Let's identify the diagonal that is fully solved.
	var fullySolvedDiag [2]*puzzle.Cell
	if unsolved == c11 || unsolved == c22 {
		fullySolvedDiag[0] = c12
		fullySolvedDiag[1] = c21
	} else {
		fullySolvedDiag[0] = c11
		fullySolvedDiag[1] = c22
	}

	// The fully solved diagonal must have the same values for this to be a deadly pattern
	if fullySolvedDiag[0].Value() != fullySolvedDiag[1].Value() {
		return false
	}

	deadlyValue := fullySolvedDiag[0].Value()
	otherValue := diag1.Value()

	// If unsolved takes the otherValue, we get a deadly pattern.
	// So we can eliminate otherValue from unsolved.
	if unsolved.HasCandidate(otherValue) {
		step := NewStep(kindAvoidableRectangle)
		step.DeleteCandidate(unsolved.Index(), otherValue)
		s.applyStep(step.WithIndices(c11.Index(), c12.Index(), c21.Index(), c22.Index()).WithValues(deadlyValue, otherValue))
		return true
	}

	return false
}

func (s *Solver) checkAvoidableRectangleType2(c11, c12, c21, c22 *puzzle.Cell) bool {
	var solved [2]*puzzle.Cell
	var unsolved [2]*puzzle.Cell

	cells := []*puzzle.Cell{c11, c12, c21, c22}
	sIdx, uIdx := 0, 0
	for _, c := range cells {
		if c.IsSolved() {
			solved[sIdx] = c
			sIdx++
		} else {
			unsolved[uIdx] = c
			uIdx++
		}
	}

	// The two solved cells MUST share a row or column (they form the "floor" of the rectangle)
	if solved[0].Row != solved[1].Row && solved[0].Col != solved[1].Col {
		return false
	}

	v1 := solved[0].Value()
	v2 := solved[1].Value()

	// The solved cells must have distinct values
	if v1 == v2 {
		return false
	}

	// Determine which unsolved cell is opposite which solved cell
	// The unsolved cell that shares a row/col with solved[0] (which has value v1)
	// MUST NOT hold candidate v1, but MUST hold candidate v2.
	// Likewise for the other unsolved cell.
	var uOppositeV1, uOppositeV2 *puzzle.Cell
	if unsolved[0].Row != solved[0].Row && unsolved[0].Col != solved[0].Col {
		// unsolved[0] is diagonally opposite to solved[0]
		uOppositeV1 = unsolved[0]
		uOppositeV2 = unsolved[1]
	} else {
		// unsolved[1] is diagonally opposite to solved[0]
		uOppositeV1 = unsolved[1]
		uOppositeV2 = unsolved[0]
	}

	// To form the deadly pattern, the diagonally opposite cells must be identical.
	// So uOppositeV1 MUST hold candidate v1 (and NOT v2, because it shares row/col with solved[1] which is v2)
	// And uOppositeV2 MUST hold candidate v2 (and NOT v1)
	if !uOppositeV1.HasCandidate(v1) || !uOppositeV2.HasCandidate(v2) {
		return false
	}
	if uOppositeV1.HasCandidate(v2) || uOppositeV2.HasCandidate(v1) {
		return false
	}

	// Type 2 AR requires that BOTH unsolved cells contain exactly one identical extra candidate Z.
	// So uOppositeV1 has candidates v1 and z, and uOppositeV2 has candidates v2 and z.
	if unsolved[0].NumCandidates() != 2 || unsolved[1].NumCandidates() != 2 {
		return false
	}

	// Find the extra candidate Z in both
	var z1, z2 int
	for _, cand := range uOppositeV1.CandidateValues() {
		if cand != v1 {
			z1 = cand
			break
		}
	}
	for _, cand := range uOppositeV2.CandidateValues() {
		if cand != v2 {
			z2 = cand
			break
		}
	}

	// The extra candidate must be the same
	if z1 != z2 || z1 == 0 {
		return false
	}

	z := z1

	found := false
	step := NewStep(kindAvoidableRectangle)
	if s.eliminateFromIntersection(unsolved[0].Index(), unsolved[1].Index(), -1, z, step) {
		s.applyStep(step.WithIndices(c11.Index(), c12.Index(), c21.Index(), c22.Index()).WithValues(v1, v2, z))
		found = true
	}

	return found
}

func (s *Solver) findUniqueRectangleType1() (found bool) {
	b := s.puzzle
	// Check each cell with exactly 2 candidate values to see if it is the base
	// corner of a unique rectangle.
	for i := range 81 {
		cell := b.Cell(i)
		if cell.NumCandidates() != 2 {
			continue
		}
		if s.checkUniqueRectangleForCell(cell) {
			return true
		}
	}

	return false
}

func (s *Solver) checkUniqueRectangleForCell(base *puzzle.Cell) (found bool) {
	b := s.puzzle

	// Look for a cell in the same row as base with the same pair of candidates.
	var rowWing *puzzle.Cell
	for c := range 9 {
		if c != base.Col {
			cell := b.Get(base.Row, c)
			if sameCandidates(base, cell) {
				rowWing = cell
				break
			}
		}
	}
	if rowWing == nil {
		return false
	}

	// Look for a cell in the same column as base with the same pair of candidates.
	var colWing *puzzle.Cell
	for r := range 9 {
		if r != base.Row {
			cell := b.Get(r, base.Col)
			if sameCandidates(base, cell) {
				colWing = cell
				break
			}
		}
	}
	if colWing == nil {
		return false
	}

	// The 2 wing cells must be in different boxes, but one of them must be in
	// the same box as the base.
	if rowWing.Box() != colWing.Box() &&
		(rowWing.Box() == base.Box() || colWing.Box() == base.Box()) {

		// These cells form a unique rectangle, so we can eliminate their candidates
		// from the cell at the 4th corner of the rectangle, which will have the
		// same row as the column-wing and the same column as the row-wing.
		step := NewStep(kindUniqueRectangle1)
		if s.eliminateValuesFromCell(colWing.Row, rowWing.Col, base.Candidates, step) {
			s.applyStep(step.
				WithValues(base.Candidates.Values()...).
				WithIndices(base.Index(), rowWing.Index(), colWing.Index()))
			return true
		}
	}

	return false
}
