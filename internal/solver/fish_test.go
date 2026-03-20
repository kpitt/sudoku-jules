package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFindFinnedSwordfish(t *testing.T) {
	// Finned Swordfish example from reglib-1.3.txt
	s := ":0311:5:.+6.4+1+7.+83+3+1+4+98..+7.8+7..+6+39147+43.29..+8..+8+3+4+6.2+7..6.+7.34+96..+75+4+8+314+3+1.+9.7..58+7631+4+9+2::557::"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findFinnedSwordfish()
	if !found {
		t.Error("Finned Swordfish should be found")
	}

	// Elimination: 557 => digit 5 at r5c7 (index 42)
	if p.Cell(42).HasCandidate(5) {
		t.Error("Candidate 5 should be eliminated from r5c7")
	}
}

func TestFindFinnedJellyfish(t *testing.T) {
	// Finned Jellyfish example from reglib-1.3.txt
	// Moved eliminations from part 4 to part 5 so they aren't pre-deleted by the reader.
	s := ":0312:7:..2.18.+6..+6..+9+4+8+1.+8+14..+597.658+1....+9+1...+59.+8...+9+8...51.354..1.82.+1+583..64+8+69.+1...::752 754 259 762 299:787::"
	p, err := puzzle.FromHodokuString(s)
	if err != nil {
		t.Fatal(err)
	}
	solver := NewSolver(p, &Options{EnableDebug: true})
	solver.processInitialValues()

	found := solver.findFinnedJellyfish()
	if !found {
		t.Error("Finned Jellyfish should be found")
	}

	// Elimination: 754 => digit 7 at r5c4 (index 39)
	if p.Cell(39).HasCandidate(7) {
		t.Error("Candidate 7 should be eliminated from r5c4")
	}
}


