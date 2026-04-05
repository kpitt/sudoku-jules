package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindXWing(t *testing.T) {
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

	// Clear candidates for val 1
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// X-Wing Setup (Rows as Base)
	// R1C1, R1C8
	// R4C1, R4C8
	setCandidate(1, 1)
	setCandidate(1, 8)
	setCandidate(4, 1)
	setCandidate(4, 8)

	// Target to eliminate (Col 1, not in base rows)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 2 candidates, preventing it from being
	// selected as a base line for X-Wing (Size 2).
	setCandidate(0, 2)
	setCandidate(0, 3)

	found := s.findXWings()
	if !found {
		t.Errorf("X-Wing technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}

func TestFindSwordfish(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 2

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 2
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// Swordfish Setup (Rows as Base)
	// R1: C1, C4
	// R4: C4, C7
	// R7: C1, C7
	setCandidate(1, 1)
	setCandidate(1, 4)
	setCandidate(4, 4)
	setCandidate(4, 7)
	setCandidate(7, 1)
	setCandidate(7, 7)

	// Target to eliminate (Col 1, Row 0)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 3 candidates, preventing it from being
	// selected as a base line for Swordfish (Size 3).
	setCandidate(0, 2)
	setCandidate(0, 3)
	setCandidate(0, 5)

	found := s.findSwordfish()
	if !found {
		t.Errorf("Swordfish technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}

func TestFindJellyfish(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	val := 3

	setCandidate := func(r, c int) {
		p.Get(r, c).Candidates.Add(val)
		s.rows[r].Unsolved[val].Add(c)
		s.columns[c].Unsolved[val].Add(r)
		_, boxLoc := getBoxLoc(r, c)
		s.boxes[p.Get(r, c).Box()].Unsolved[val].Add(boxLoc)
	}

	// Clear candidates for val 3
	for r := 0; r < 9; r++ {
		s.rows[r].Unsolved[val].Clear()
		for c := 0; c < 9; c++ {
			p.Get(r, c).RemoveCandidate(val)
		}
	}
	for c := 0; c < 9; c++ {
		s.columns[c].Unsolved[val].Clear()
	}
	for b := 0; b < 9; b++ {
		s.boxes[b].Unsolved[val].Clear()
	}

	// Jellyfish Setup (Rows as Base)
	// R1: C1, C2, C4, C5
	// R2: C1, C2, C4, C5
	// R4: C1, C2, C4, C5
	// R5: C1, C2, C4, C5
	rows := []int{1, 2, 4, 5}
	cols := []int{1, 2, 4, 5}

	for _, r := range rows {
		for _, c := range cols {
			setCandidate(r, c)
		}
	}

	// Target to eliminate (Col 1, Row 0)
	setCandidate(0, 1)
	// Add noise to Row 0 so it has > 4 candidates, preventing it from being
	// selected as a base line for Jellyfish (Size 4).
	setCandidate(0, 2)
	setCandidate(0, 3)
	setCandidate(0, 6)
	setCandidate(0, 7)
	setCandidate(0, 8)

	found := s.findJellyfish()
	if !found {
		t.Errorf("Jellyfish technique should have been found")
	}

	if p.Get(0, 1).HasCandidate(val) {
		t.Errorf("Target (0,1) should have eliminated candidate %d", val)
	}
}
