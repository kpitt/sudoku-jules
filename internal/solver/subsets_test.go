package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindNakedPairs(t *testing.T) {
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

	// Naked Pair in Row 0: Cells (0,0) and (0,1) with candidates {1, 2}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2) // Base
	setCandidate(0, 1, 1)
	setCandidate(0, 1, 2)

	// Other cells in Row 0 that have candidates 1 and 2
	setCandidate(0, 2, 1) // Target to eliminate
	setCandidate(0, 3, 2) // Target to eliminate

	// Add other candidates so these cells aren't empty
	setCandidate(0, 2, 3)
	setCandidate(0, 3, 4)

	found := s.findNakedPairs()
	if !found {
		t.Errorf("Naked Pair technique should have been found")
	}

	if p.Get(0, 2).HasCandidate(1) {
		t.Errorf("Target (0,2) should have eliminated candidate 1")
	}
	if p.Get(0, 3).HasCandidate(2) {
		t.Errorf("Target (0,3) should have eliminated candidate 2")
	}
}

func TestFindNakedTriples(t *testing.T) {
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

	// Naked Triple: {1,2}, {2,3}, {1,3}
	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 1)
	setCandidate(0, 2, 3)

	setCandidate(0, 4, 1) // Target
	setCandidate(0, 4, 2) // Target
	setCandidate(0, 4, 4) // Noise

	found := s.findNakedTriples()
	if !found {
		t.Errorf("Naked Triple technique should have been found")
	}

	if p.Get(0, 4).HasCandidate(1) || p.Get(0, 4).HasCandidate(2) {
		t.Errorf("Target (0,4) should have eliminated candidates 1 and 2")
	}
}

func TestFindNakedQuadruples(t *testing.T) {
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

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 3)
	setCandidate(0, 2, 4)
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4)

	setCandidate(0, 5, 1) // Target
	setCandidate(0, 5, 5) // Noise

	found := s.findNakedQuadruples()
	if !found {
		t.Errorf("Naked Quadruple technique should have been found")
	}

	if p.Get(0, 5).HasCandidate(1) {
		t.Errorf("Target (0,5) should have eliminated candidate 1")
	}
}

func TestFindHiddenPairs(t *testing.T) {
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

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 1)
	setCandidate(0, 1, 2)

	setCandidate(0, 0, 3) // Target
	setCandidate(0, 1, 4) // Target

	setCandidate(0, 2, 3)
	setCandidate(0, 3, 4)
	setCandidate(0, 4, 5)

	found := s.findHiddenPairs()
	if !found {
		t.Errorf("Hidden Pair technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(3) {
		t.Errorf("Target (0,0) should have eliminated candidate 3")
	}
	if p.Get(0, 1).HasCandidate(4) {
		t.Errorf("Target (0,1) should have eliminated candidate 4")
	}
}

func TestFindHiddenTriples(t *testing.T) {
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

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 1)
	setCandidate(0, 2, 3)

	setCandidate(0, 0, 4) // Target
	setCandidate(0, 1, 5) // Target
	setCandidate(0, 2, 6) // Target

	setCandidate(0, 3, 4)
	setCandidate(0, 3, 5)
	setCandidate(0, 4, 5)
	setCandidate(0, 4, 6)

	found := s.findHiddenTriples()
	if !found {
		t.Errorf("Hidden Triple technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(4) {
		t.Errorf("Target (0,0) should have eliminated candidate 4")
	}
	if p.Get(0, 1).HasCandidate(5) {
		t.Errorf("Target (0,1) should have eliminated candidate 5")
	}
	if p.Get(0, 2).HasCandidate(6) {
		t.Errorf("Target (0,2) should have eliminated candidate 6")
	}
}

func TestFindHiddenQuadruples(t *testing.T) {
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

	setCandidate(0, 0, 1)
	setCandidate(0, 0, 2)
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 3)
	setCandidate(0, 2, 4)
	setCandidate(0, 3, 1)
	setCandidate(0, 3, 4)

	setCandidate(0, 0, 5) // Target
	setCandidate(0, 1, 6) // Target

	setCandidate(0, 4, 5)
	setCandidate(0, 4, 6)
	setCandidate(0, 5, 7)
	setCandidate(0, 5, 8)

	found := s.findHiddenQuadruples()
	if !found {
		t.Errorf("Hidden Quadruple technique should have been found")
	}

	if p.Get(0, 0).HasCandidate(5) {
		t.Errorf("Target (0,0) should have eliminated candidate 5")
	}
	if p.Get(0, 1).HasCandidate(6) {
		t.Errorf("Target (0,1) should have eliminated candidate 6")
	}
}
