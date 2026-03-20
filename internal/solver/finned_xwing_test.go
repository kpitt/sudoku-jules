package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindFinnedXWing(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 1

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
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

	// Finned X-Wing Setup (Rows as Base)
	// R1: C1, C8
	// R4: C8, C2 (Fin in Box 3)
	// Corner we'd need for perfect X-Wing is R4C1.
	// R4C2 is the fin. Box 3 contains Col 0,1,2 of Rows 3,4,5.
	// Both (4,1) and (4,2) are in Box 3.
	
	setCandidate(1, 1)
	setCandidate(1, 8)
	
	setCandidate(4, 8)
	setCandidate(4, 2) // Fin
	// Note: (4,1) is missing, so it's a Sashimi X-Wing.
	
	// Target to eliminate: (5,1). 
	// (5,1) sees (1,1) via Column 1.
	// (5,1) sees (4,2) via Box 3.
	setCandidate(5, 1)
	p.Get(5, 1).Candidates.Add(2) // Noise

	found := s.findFinnedXWings()
	if !found {
		t.Error("Finned X-Wing should be found")
	}

	if p.Get(5, 1).HasCandidate(1) {
		t.Error("Candidate 1 should be eliminated from (5,1)")
	}
}
