package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindUniqueRectangleType3(t *testing.T) {
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

	// UR Type 3 Setup: {1, 2} with extra {3, 4}
	// Corners: (0,0), (0,3), (1,0), (1,3)
	// (0,0) and (1,0) in Box 0.
	// (0,3) and (1,3) in Box 1.
	
	// Bivalue cells (Column 0)
	setCandidate(0, 0, 1); setCandidate(0, 0, 2)
	setCandidate(1, 0, 1); setCandidate(1, 0, 2)
	
	// Extra cells (Column 3)
	// (0,3) has extra 3
	setCandidate(0, 3, 1); setCandidate(0, 3, 2); setCandidate(0, 3, 3)
	// (1,3) has extra 4
	setCandidate(1, 3, 1); setCandidate(1, 3, 2); setCandidate(1, 3, 4)
	
	// Naked Pair helper in Column 3: (4,3) has {3, 4}
	setCandidate(4, 3, 3); setCandidate(4, 3, 4)
	
	// Target cell in Column 3: (5,3) has {3}
	setCandidate(5, 3, 3)
	setCandidate(5, 3, 5) // Noise

	found := s.findUniqueRectangleType3()
	if !found {
		t.Error("Unique Rectangle Type 3 should be found")
	}

	if p.Get(5, 3).HasCandidate(3) {
		t.Error("Candidate 3 should be eliminated from (5,3)")
	}
}
