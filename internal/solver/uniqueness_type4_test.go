package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindUniqueRectangleType4(t *testing.T) {
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

	// UR Type 4 Setup: {1, 2}
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) in Box 0.
	// (0,3) and (1,3) in Box 1.
	
	// Bivalue cells: (1,0) and (1,3) [Row 1]
	setCandidate(1, 0, 1); setCandidate(1, 0, 2)
	setCandidate(1, 3, 1); setCandidate(1, 3, 2)
	
	// Extra cells: (0,0) and (0,3) [Row 0]
	// They have extra candidates
	setCandidate(0, 0, 1); setCandidate(0, 0, 2); setCandidate(0, 0, 3)
	setCandidate(0, 3, 1); setCandidate(0, 3, 2); setCandidate(0, 3, 4)
	
	// For Type 4, one of the candidates (say 1) must be restricted to (0,0) and (0,3) in Row 0.
	// So we ensure candidate 1 does not appear elsewhere in Row 0.
	// It already doesn't because we cleared it.
	
	// Run UR Type 4
	found := s.findUniqueRectangleType4()
	if !found {
		t.Error("Unique Rectangle Type 4 should be found")
	}

	// Elimination: Candidate 2 (other than the conjugate pair 1) should be removed from (0,0) and (0,3).
	if p.Get(0, 0).HasCandidate(2) {
		t.Error("Candidate 2 should be eliminated from (0,0)")
	}
	if p.Get(0, 3).HasCandidate(2) {
		t.Error("Candidate 2 should be eliminated from (0,3)")
	}
}
