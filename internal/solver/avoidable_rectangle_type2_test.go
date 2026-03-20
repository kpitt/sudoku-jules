package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindAvoidableRectangleType2(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			p.Get(r, c).Candidates.Clear()
			for v := 1; v <= 9; v++ {
				s.rows[r].Unsolved[v].Clear()
				s.columns[c].Unsolved[v].Clear()
				s.boxes[p.Get(r, c).Box()].Unsolved[v].Clear()
			}
		}
	}

	// Avoidable Rectangle Type 2 Setup: {1, 2} with extra 3
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) in Box 0.
	// (0,3) and (1,3) in Box 1.
	
	// Two solved non-given cells in Row 1:
	// (1,0) = 1, (1,3) = 2
	p.Get(1, 0).PlaceValue(1)
	p.Get(1, 3).PlaceValue(2)
	p.Get(1, 0).IsGiven = false
	p.Get(1, 3).IsGiven = false
	
	// Two unsolved cells in Row 0:
	// (0,0) and (0,3) both have {1, 2, 3}
	setCandidate(0, 0, 1); setCandidate(0, 0, 2); setCandidate(0, 0, 3)
	setCandidate(0, 3, 1); setCandidate(0, 3, 2); setCandidate(0, 3, 3)
	
	// Target cell in Row 0: (0,1) sees both (0,0) and (0,3)
	setCandidate(0, 1, 3)

	found := s.findAvoidableRectangles() // Assuming it covers Type 2
	if !found {
		t.Error("Avoidable Rectangle Type 2 should be found")
	}

	if p.Get(0, 1).HasCandidate(3) {
		t.Error("Candidate 3 should be eliminated from (0,1)")
	}
}
