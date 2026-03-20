package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindForcingNet(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	// Setup a simple Forcing Net
	// Cell (0,0) has candidates {1, 2}
	// If (0,0)=1 => (0,5)=3
	// If (0,0)=2 => (0,5)=3
	
	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// (0,0) = {1, 2}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)

	// (0,1) = {1, 3}
	// (0,0)=1 => (0,1) cannot be 1 => (0,1) must be 3
	setCandidate(0, 1, 1)
	setCandidate(0, 1, 3)

	// (0,2) = {2, 3}
	// (0,0)=2 => (0,2) cannot be 2 => (0,2) must be 3
	setCandidate(0, 2, 2)
	setCandidate(0, 2, 3)

	// (0,5) = {3, 4}
	// We want (0,1)=3 => (0,5)=3? No, that's not how it works.
	// We need (0,1)=3 => (0,5) cannot be something?
	
	// Let's use:
	// If (0,0)=1 => (0,1)=3
	// If (0,0)=2 => (0,2)=3
	// (0,1) and (0,2) both see (0,5).
	// If (0,1)=3 => (0,5) cannot be 3.
	// If (0,2)=3 => (0,5) cannot be 3.
	// So (0,5) cannot be 3.
	
	setCandidate(0, 5, 3)
	setCandidate(0, 5, 4)

	// Need to make sure the implications are added to the graph.
	// (0,0)=1 => (0,1)=3: This is a bivalue cell implication.
	// (0,0)=2 => (0,2)=3: This is a bivalue cell implication.
	// (0,1)=3 => (0,5)<>3: This is a weak link (same row).
	// (0,2)=3 => (0,5)<>3: This is a weak link (same row).

	found := s.findForcingNets()
	if !found {
		t.Errorf("Forcing Net should have been found")
	}

	if p.Get(0, 5).HasCandidate(3) {
		t.Errorf("Candidate 3 should have been eliminated from (0,5)")
	}
}
