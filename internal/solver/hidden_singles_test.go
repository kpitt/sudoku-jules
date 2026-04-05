package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindHiddenSingles(t *testing.T) {
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

	// Setup: Candidate 1 is only in (0,0) for Row 0
	setCandidate(0, 0, 1)

	// Add other candidates to (0,0) so it's not a Naked Single
	setCandidate(0, 0, 2)
	setCandidate(0, 0, 3)

	// Other cells in Row 0 get candidates 2 and 3, but NOT 1
	setCandidate(0, 1, 2)
	setCandidate(0, 1, 3)
	setCandidate(0, 2, 2)
	setCandidate(0, 2, 3)

	found := s.findHiddenSingles()
	if !found {
		t.Errorf("Hidden Single technique should have been found")
	}

	// Because finding a Hidden Single places the value, the cell should be
	// marked as solved with the value 1.
	if !p.Get(0, 0).IsSolved() {
		t.Errorf("Hidden Single should have placed a value in cell (0,0)")
	} else if p.Get(0, 0).Value() != 1 {
		t.Errorf("Hidden Single should have placed value 1 in cell (0,0)")
	}
}
