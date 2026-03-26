package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindForcingNet(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// Clear all candidates
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

	// Setup a simple Forcing Net
	// (0,0)={1, 2}
	// (0,1)={1, 3}
	// (1,0)={2, 3}
	// (1,1)={3, 4}
	// If (0,0)=1 => (0,1)=3 => (1,1)=4
	// If (0,0)=2 => (1,0)=3 => (1,1)=4
	// Result: (1,1)=4

	setCandidate(0, 0, 1); setCandidate(0, 0, 2)
	setCandidate(0, 1, 1); setCandidate(0, 1, 3)
	setCandidate(1, 0, 2); setCandidate(1, 0, 3)
	setCandidate(1, 1, 3); setCandidate(1, 1, 4)

	// Add noise to avoid Hidden Singles/Contradictions
	// Fill the rest of the board with all candidates
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if r < 2 && c < 2 {
				continue
			}
			for v := 1; v <= 9; v++ {
				setCandidate(r, c, v)
			}
		}
	}

	found := s.findForcingNets()
	if !found {
		t.Errorf("Forcing Net should have been found")
	}

	if p.Get(1, 1).Value() != 4 && !p.Get(1, 1).HasCandidate(4) {
		// Wait, if it's found, it should either have placed it or eliminated 3.
		// In my setup, (1,1)={3, 4}. If it eliminates 3, it becomes 4 (Naked Single).
	}
	
	if p.Get(1, 1).Value() != 4 {
		t.Errorf("Cell (1,1) should have been solved with value 4, got %d", p.Get(1, 1).Value())
	}
}

func TestFindForcingChainContradiction(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	setCandidate := func(r, c, v int) {
		p.Get(r, c).Candidates.Add(v)
		s.rows[r].Unsolved[v].Add(c)
		s.columns[c].Unsolved[v].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[v].Add(boxLoc)
	}

	// Clear all candidates
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

	// Contradiction: (0,0)=1 => ... => NOT (0,0)=1.
	// (0,0)=1 => (0,5) NOT 1
	// (0,5)={1, 2} => (0,5)=2
	// (0,5) and (5,5) only 2s in Col 5 => (5,5) NOT 2
	// (5,5)={2, 3} => (5,5)=3
	// (5,5) and (5,0) only 3s in Row 5 => (5,0) NOT 3
	// (5,0)={3, 1} => (5,0)=1
	// (5,0)=1 => (0,0) NOT 1 (since they are in Col 0).

	setCandidate(0, 0, 1); setCandidate(0, 0, 4) // (0,0) has 1
	setCandidate(0, 5, 1); setCandidate(0, 5, 2) // (0,5) is bivalue {1, 2}
	setCandidate(5, 5, 2); setCandidate(5, 5, 3) // (5,5) is bivalue {2, 3}
	setCandidate(5, 0, 3); setCandidate(5, 0, 1) // (5,0) is bivalue {3, 1}

	// House Strong Links:
	// Row 0: only two 1s at (0,0) and (0,5)
	s.rows[0].Unsolved[1].Clear()
	s.rows[0].Unsolved[1].Add(0) // (0,0)
	s.rows[0].Unsolved[1].Add(5) // (0,5)

	// Col 5: only two 2s at (0,5) and (5,5)
	s.columns[5].Unsolved[2].Clear()
	s.columns[5].Unsolved[2].Add(0) // (0,5)
	s.columns[5].Unsolved[2].Add(5) // (5,5)
	
	// Row 5: only two 3s at (5,5) and (5,0)
	s.rows[5].Unsolved[3].Clear()
	s.rows[5].Unsolved[3].Add(5) // (5,5)
	s.rows[5].Unsolved[3].Add(0) // (5,0)

	// Col 0: only two 1s at (5,0) and (0,0)
	s.columns[0].Unsolved[1].Clear()
	s.columns[0].Unsolved[1].Add(5) // (5,0)
	s.columns[0].Unsolved[1].Add(0) // (0,0)

	found := s.findForcingChains()
	if !found {
		t.Errorf("Forcing Chain contradiction should have been found")
	}

	if p.Get(0, 0).HasCandidate(1) {
		t.Errorf("Candidate 1 should have been eliminated from (0,0)")
	}
}
