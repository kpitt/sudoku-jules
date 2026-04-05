package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestSolveBruteForce(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		solved bool
	}{
		{
			name:   "Partially Filled",
			input:  "530070000600195000098000060800060003400803001700020006060000280000419005000080079",
			solved: true,
		},
		{
			name:   "Empty Puzzle",
			input:  ".................................................................................",
			solved: true,
		},
		{
			name:   "Fully Solved",
			input:  "534678912672195348198342567859761423426853791713924856961537284287419635345286179",
			solved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := puzzle.FromString(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse puzzle: %v", err)
			}

			s := NewSolver(p, nil)
			s.SolveBruteForce()

			if p.IsSolved() != tt.solved {
				t.Errorf("Expected IsSolved() to be %v, got %v", tt.solved, p.IsSolved())
			}

			if tt.solved {
				if err := p.ValidateSolution(); err != nil {
					t.Errorf("Puzzle solution is invalid: %v", err)
				}
			}
		})
	}
}
