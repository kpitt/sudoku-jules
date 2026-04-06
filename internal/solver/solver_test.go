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

func TestSolver_RemoveCellCandidate(t *testing.T) {
	t.Run("Standard Elimination", func(t *testing.T) {
		p, _ := puzzle.FromString(".................................................................................")
		s := NewSolver(p, nil)

		idx := 0
		val := 5
		r, c := idx/9, idx%9
		box, boxLoc := getBoxLoc(r, c)

		// Before removal, candidate should be present
		if !p.Cell(idx).Candidates.Contains(val) {
			t.Errorf("Expected cell %d to contain candidate %d", idx, val)
		}
		if !s.rows[r].Unsolved[val].Contains(c) {
			t.Errorf("Expected row %d to contain candidate %d at loc %d", r, val, c)
		}
		if !s.columns[c].Unsolved[val].Contains(r) {
			t.Errorf("Expected col %d to contain candidate %d at loc %d", c, val, r)
		}
		if !s.boxes[box].Unsolved[val].Contains(boxLoc) {
			t.Errorf("Expected box %d to contain candidate %d at loc %d", box, val, boxLoc)
		}

		s.removeCellCandidate(idx, val)

		// After removal, candidate should not be present
		if p.Cell(idx).Candidates.Contains(val) {
			t.Errorf("Expected cell %d to not contain candidate %d", idx, val)
		}
		if s.rows[r].Unsolved[val].Contains(c) {
			t.Errorf("Expected row %d to not contain candidate %d at loc %d", r, val, c)
		}
		if s.columns[c].Unsolved[val].Contains(r) {
			t.Errorf("Expected col %d to not contain candidate %d at loc %d", c, val, r)
		}
		if s.boxes[box].Unsolved[val].Contains(boxLoc) {
			t.Errorf("Expected box %d to not contain candidate %d at loc %d", box, val, boxLoc)
		}
	})

	t.Run("Naked Single Trigger", func(t *testing.T) {
		p, _ := puzzle.FromString(".................................................................................")
		s := NewSolver(p, nil)

		idx := 0

		// Remove all candidates except 1 and 2
		for val := 3; val <= 9; val++ {
			s.removeCellCandidate(idx, val)
		}

		// Currently, cell has candidates 1 and 2. solution len is 0.
		initialSolLen := len(s.solution)

		// Remove candidate 1. This should leave only candidate 2, triggering Naked Single.
		s.removeCellCandidate(idx, 1)

		if len(s.solution) != initialSolLen+1 {
			t.Fatalf("Expected solution to have 1 new step, got %d", len(s.solution)-initialSolLen)
		}

		lastStep := s.solution[len(s.solution)-1]
		if lastStep.technique != kindNakedSingle {
			t.Errorf("Expected Naked Single step, got %v", lastStep.technique)
		}
		if len(lastStep.indices) != 1 || lastStep.indices[0] != idx {
			t.Errorf("Expected Naked Single step to have index %d, got %v", idx, lastStep.indices)
		}
		if len(lastStep.values) != 1 || lastStep.values[0] != 2 {
			t.Errorf("Expected Naked Single step to have value 2, got %v", lastStep.values)
		}
	})
}
