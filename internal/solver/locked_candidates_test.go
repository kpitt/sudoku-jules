package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindPointingTuples(t *testing.T) {
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

	// Setup: Candidate 1 in Box 0 is restricted to Row 0
	setCandidate(0, 0, 1)
	setCandidate(0, 1, 1)

	// Add noise to avoid singles
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 3)

	// Target to eliminate in Row 0, outside Box 0
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4) // Noise

	found := s.findPointingTuples()
	if !found {
		t.Errorf("Pointing Tuples technique should have been found")
	}

	if p.Get(0, 3).HasCandidate(1) {
		t.Errorf("Target (0,3) should have eliminated candidate 1")
	}
}

func TestFindClaimingTuples(t *testing.T) {
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

	// Setup: Candidate 1 in Row 0 is restricted to Box 0
	setCandidate(0, 0, 1)
	setCandidate(0, 1, 1)

	// Add noise to avoid singles
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 3)

	// Target to eliminate in Box 0, outside Row 0
	setCandidate(1, 0, 1)
	setCandidate(1, 0, 4) // Noise

	found := s.findClaimingTuples()
	if !found {
		t.Errorf("Claiming Tuples technique should have been found")
	}

	if p.Get(1, 0).HasCandidate(1) {
		t.Errorf("Target (1,0) should have eliminated candidate 1")
	}
}
