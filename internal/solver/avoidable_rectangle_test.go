package solver

import (
	"testing"
	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindAvoidableRectangles(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// An Avoidable Rectangle requires exactly two boxes, two rows, two columns.
	// E.g., r0c0, r0c3, r1c0, r1c3
	// Box 0 (r0c0, r1c0) and Box 1 (r0c3, r1c3)

	// Set cell to not given
	p.Get(0,0).IsGiven = false
	p.PlaceValue(0, 7) // r0c0 = 7

	p.Get(0,3).IsGiven = false
	p.PlaceValue(3, 9) // r0c3 = 9

	p.Get(1,0).IsGiven = false
	p.PlaceValue(9, 9) // r1c0 = 9

	// r1c3 is at index 1*9 + 3 = 12
	c := p.Cell(12)
	c.Candidates.Clear()
	c.Candidates.Add(7)
	c.Candidates.Add(9)
	c.Candidates.Add(3)

	s.rows[1].Unsolved[7].Add(3)
	s.rows[1].Unsolved[9].Add(3)
	s.rows[1].Unsolved[3].Add(3)
	s.columns[3].Unsolved[7].Add(1)
	s.columns[3].Unsolved[9].Add(1)
	s.columns[3].Unsolved[3].Add(1)
	_, boxLoc := getBoxLoc(1, 3)
	s.boxes[c.Box()].Unsolved[7].Add(boxLoc)
	s.boxes[c.Box()].Unsolved[9].Add(boxLoc)
	s.boxes[c.Box()].Unsolved[3].Add(boxLoc)

	found := s.findAvoidableRectangles()
	if !found {
		t.Errorf("Expected to find Avoidable Rectangle Type 1")
	}

	if c.HasCandidate(7) {
		t.Errorf("Expected candidate 7 to be eliminated from r1c3")
	}
}

func TestFindAvoidableRectanglesType2(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// An Avoidable Rectangle Type 2 requires two placed cells on the same row or col
	// Let's use r0c0 and r0c3 (Row 0) with values 7 and 9.
	// The unsolved cells will be r1c0 and r1c3.

	p.Get(0,0).IsGiven = false
	p.PlaceValue(0, 7) // r0c0 = 7

	p.Get(0,3).IsGiven = false
	p.PlaceValue(3, 9) // r0c3 = 9

	// r1c0 is diagonally opposite to r0c3 (which is 9). So it must be able to be 9.
	// Thus r1c0 has candidates 9 and 3.
	c1 := p.Cell(9)
	c1.Candidates.Clear()
	c1.Candidates.Add(9)
	c1.Candidates.Add(3) // Shared extra candidate

	// r1c3 is diagonally opposite to r0c0 (which is 7). So it must be able to be 7.
	// Thus r1c3 has candidates 7 and 3.
	c2 := p.Cell(12)
	c2.Candidates.Clear()
	c2.Candidates.Add(7)
	c2.Candidates.Add(3) // Shared extra candidate

	// Add candidates to houses so the solver's internal lookup works
	s.rows[1].Unsolved[9].Add(0)
	s.rows[1].Unsolved[3].Add(0)
	s.columns[0].Unsolved[9].Add(1)
	s.columns[0].Unsolved[3].Add(1)
	_, boxLoc1 := getBoxLoc(1, 0)
	s.boxes[c1.Box()].Unsolved[9].Add(boxLoc1)
	s.boxes[c1.Box()].Unsolved[3].Add(boxLoc1)

	s.rows[1].Unsolved[7].Add(3)
	s.rows[1].Unsolved[3].Add(3)
	s.columns[3].Unsolved[7].Add(1)
	s.columns[3].Unsolved[3].Add(1)
	_, boxLoc2 := getBoxLoc(1, 3)
	s.boxes[c2.Box()].Unsolved[7].Add(boxLoc2)
	s.boxes[c2.Box()].Unsolved[3].Add(boxLoc2)

    // A cell that sees both c1 (r1c0) and c2 (r1c3) is r1c1
    c_target := p.Cell(10) // r1c1
    c_target.Candidates.Clear()
    c_target.Candidates.Add(3)
    c_target.Candidates.Add(5)

	found := s.findAvoidableRectangles()
	if !found {
		t.Errorf("Expected to find Avoidable Rectangle Type 2")
	}

	if c_target.HasCandidate(3) {
		t.Errorf("Expected candidate 3 to be eliminated from r1c1")
	}
}

func TestFindAvoidableRectanglesType3(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// Avoidable Rectangle Type 3:
	// r0c0=7, r0c3=9 (solved)
	// r1c0, r1c3 (unsolved)
	// r1c0 has candidates {9, 2}
	// r1c3 has candidates {7, 3}
	// Pseudo-cell: {2, 3}
	// Add another cell r1c1 with candidates {2, 3}
	// Together they form a naked pair {2, 3} in row 1.
	// Eliminate 2 and 3 from r1c2.

	p.Get(0, 0).IsGiven = false
	p.PlaceValue(0, 7) // r0c0 = 7

	p.Get(0, 3).IsGiven = false
	p.PlaceValue(3, 9) // r0c3 = 9

	c1 := p.Cell(9) // r1c0 (opposite 9)
	c1.Candidates.Clear()
	c1.Candidates.Add(9)
	c1.Candidates.Add(2)

	c2 := p.Cell(12) // r1c3 (opposite 7)
	c2.Candidates.Clear()
	c2.Candidates.Add(7)
	c2.Candidates.Add(3)

	c3 := p.Cell(10) // r1c1
	c3.Candidates.Clear()
	c3.Candidates.Add(2)
	c3.Candidates.Add(3)

	cTarget := p.Cell(11) // r1c2
	cTarget.Candidates.Clear()
	cTarget.Candidates.Add(2)
	cTarget.Candidates.Add(5)

	// Update solver state
	s.rows[1].Unsolved[9].Add(0)
	s.rows[1].Unsolved[2].Add(0)
	s.rows[1].Unsolved[7].Add(3)
	s.rows[1].Unsolved[3].Add(3)
	s.rows[1].Unsolved[2].Add(1)
	s.rows[1].Unsolved[3].Add(1)
	s.rows[1].Unsolved[2].Add(2)
	s.rows[1].Unsolved[5].Add(2)

	found := s.findAvoidableRectangles()
	if !found {
		t.Errorf("Expected to find Avoidable Rectangle Type 3")
	}

	if cTarget.HasCandidate(2) {
		t.Errorf("Expected candidate 2 to be eliminated from r1c2")
	}
}

func TestFindAvoidableRectanglesType4(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// Avoidable Rectangle Type 4:
	// r0c0=7, r0c3=9 (solved)
	// r1c0, r1c3 (unsolved)
	// r1c0 has candidates {9, 2}
	// r1c3 has candidates {7, 2}
	// Conjugate pair: Value 9 is restricted to r1c0 and r1c3 in row 1?
	// Wait, r1c3 doesn't have 9.
	// Actually, Type 4: v1 is conjugate in a house.
	// r1c0 has {9, 2}, r1c3 has {7, 2}.
	// If 9 is restricted to r1c0 and some other cell? No.
	// If 9 is restricted to r1c0 and r1c3 in row 1.
	// But r1c3 only has {7, 2}. So 9 is a single in r1c0?
	// Let's make it: r1c0 has {9, 7, 2}, r1c3 has {7, 9, 2}.
	// Value 9 is restricted to r1c0 and r1c3 in row 1.
	// Then eliminate 7 from r1c0 and r1c3? No, eliminate the OTHER deadly value.
	// Deadly values are 9 and 7.
	// If 9 is conjugate, eliminate 7.

	p.Get(0, 0).IsGiven = false
	p.PlaceValue(0, 7) // r0c0 = 7

	p.Get(0, 3).IsGiven = false
	p.PlaceValue(3, 9) // r0c3 = 9

	c1 := p.Cell(9) // r1c0 (opposite 9)
	c1.Candidates.Clear()
	c1.Candidates.Add(9)
	c1.Candidates.Add(7)
	c1.Candidates.Add(2)

	c2 := p.Cell(12) // r1c3 (opposite 7)
	c2.Candidates.Clear()
	c2.Candidates.Add(7)
	c2.Candidates.Add(9)
	c2.Candidates.Add(2)

	// Row 1: 9 is restricted to col 0 and 3.
	s.rows[1].Unsolved[9].Clear()
	s.rows[1].Unsolved[9].Add(0)
	s.rows[1].Unsolved[9].Add(3)

	s.rows[1].Unsolved[7].Add(0)
	s.rows[1].Unsolved[7].Add(3)
	s.rows[1].Unsolved[2].Add(0)
	s.rows[1].Unsolved[2].Add(3)

	found := s.findAvoidableRectangles()
	if !found {
		t.Errorf("Expected to find Avoidable Rectangle Type 4")
	}

	if c1.HasCandidate(7) || c2.HasCandidate(7) {
		t.Errorf("Expected candidate 7 to be eliminated from r1c0 and r1c3")
	}
}
