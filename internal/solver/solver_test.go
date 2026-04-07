package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestCheckPuzzle(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   bool
		wantErr     error
	}{
		{
			name:      "Valid puzzle",
			input:     "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79",
			wantValid: true,
			wantErr:   nil,
		},
		{
			name:      "Multiple solutions",
			// Remove enough givens so it has multiple solutions
			input:     ".................................................................................",
			wantValid: false,
			wantErr:   puzzle.ErrMultipleSolutions,
		},
		{
			name:      "No solution",
			// A puzzle with an obvious contradiction. But we can't use FromString to parse invalid ones,
			// because it throws error before. So we'll construct it manually.
			input:     "no_solution_manual",
			wantValid: false,
			wantErr:   puzzle.ErrNoSolution,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p *puzzle.Puzzle
			if tt.input == "no_solution_manual" {
				// Let's create a puzzle that is logically impossible but syntactically valid.
				// We can just construct a completely blocked puzzle state where Dancing Links will fail.
				// For example, place given values such that one column cannot have a '1' anywhere.
				p, _ = puzzle.FromString(".................................................................................")
				p.GivenValue(0, 1) // r1c1 = 1, so r1c2 cannot be 1.
				p.GivenValue(10, 1) // r2c2 = 1, r2c2 cannot be 1, wait, r2c2 is index 1*9+1=10
				p.GivenValue(20, 1) // r3c3 = 1, index 2*9+2=20
				// This still has many solutions. We want 0 solutions.
				// Wait, let's just use the valid string but clear a candidate? DancingLinks uses the puzzle candidates!
				p, _ = puzzle.FromString("53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79")

				// Apply deductive steps to reach a dead end or just modify the puzzle directly
				// by setting given values that are mathematically impossible.
				// Like placing 1,2,3,4,5,6,7,8 in a row, and the remaining cell must be 9.
				// But we make sure 9 is already in its column.
				// Let's manually set a cell's candidates to 0
				p.Cells[0].Candidates.Clear()
				p.Cells[0].GivenValue(0) // It's not a given value? No, just clear candidates.

				// Actually, DancingLinks converts puzzle candidates to matrix rows. If a cell has 0 candidates,
				// then no rows will be created for that cell. The column for that cell will be empty,
				// and DancingLinks will find 0 solutions!
				// BUT Wait! In `dancing_links.go`, the cell constraint columns are only added if the cell is unsolved.
				// Wait, if candidates is empty, then size of column is 0.
				// Let's try to clear candidates from an UNSOLVED cell.
				// Index 2 is '.' in the string (r1c3). Let's clear its candidates.
				p.Cells[2].Candidates.Clear()
			} else {
				var err error
				p, err = puzzle.FromString(tt.input)
				if err != nil {
					t.Fatalf("puzzle.FromString failed: %v", err)
				}
			}

			valid, err := CheckPuzzle(p)
			if valid != tt.wantValid {
				t.Errorf("CheckPuzzle() valid = %v, want %v", valid, tt.wantValid)
			}
			if err != tt.wantErr {
				t.Errorf("CheckPuzzle() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

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
