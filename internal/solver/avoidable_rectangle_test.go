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
