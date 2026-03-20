package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindHiddenRectangle(t *testing.T) {
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

	// Hidden Rectangle Setup: {1, 2}
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) in Box 0.
	// (0,3) and (1,3) in Box 1.
	
	// C4 is bivalue {1, 2} at (1,3)
	setCandidate(1, 3, 1); setCandidate(1, 3, 2)
	
	// C2, C3 contain {1, 2}
	setCandidate(0, 3, 1); setCandidate(0, 3, 2)
	setCandidate(1, 0, 1); setCandidate(1, 0, 2)
	
	// C1 contains {1, 2, 3} at (0,0)
	setCandidate(0, 0, 1); setCandidate(0, 0, 2); setCandidate(0, 0, 3)
	
	// Conjugate pair for digit 1 in Row 0: restricted to (0,0) and (0,3)
	// (It's already restricted because we cleared Row 0)
	// Conjugate pair for digit 1 in Col 0: restricted to (0,0) and (1,0)
	// (It's already restricted)
	
	// Eliminate digit 2 from (0,0)
	found := s.findHiddenRectangle()
	if !found {
		t.Error("Hidden Rectangle should be found")
	}

	if p.Get(0, 0).HasCandidate(2) {
		t.Error("Candidate 2 should be eliminated from (0,0)")
	}
}
