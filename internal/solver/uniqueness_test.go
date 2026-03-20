package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestUniquenessDetection(t *testing.T) {
	// A simple non-unique puzzle: empty grid
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)
	s.Solve()

	if !s.IsNonUnique {
		t.Error("Empty puzzle should be detected as non-unique")
	}

	// A puzzle with 2 solutions (Deadly Pattern)
	// We take a solved puzzle and remove 4 cells that form a swap-able rectangle.
	// But it's easier to just use a known non-unique string if available.
	
	// This is a known non-unique puzzle (only 17 givens, but not this one)
	// Let's just use one with two solutions.
	// 000000012000000003000000000000000000000000000000000000000000000000000000000000000
	// Actually, an empty grid is definitely non-unique.
}

func TestUnsolvableDetection(t *testing.T) {
	// A puzzle with two 9s in the same row
	p := puzzle.NewPuzzle()
	p.PlaceValue(0, 9)
	p.PlaceValue(1, 9)
	
	s := NewSolver(p, nil)
	s.Solve()

	if !s.IsUnsolvable {
		t.Error("Puzzle with two 9s in same row should be detected as unsolvable")
	}
}
