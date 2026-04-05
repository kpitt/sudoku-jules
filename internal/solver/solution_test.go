package solver

import (
	"testing"

	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestFormatPlacedValue(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		indices []int
		want    string
	}{
		{
			name:    "happy path r1c1",
			values:  []int{5},
			indices: []int{0},
			want:    "r1c1=5",
		},
		{
			name:    "happy path r9c9",
			values:  []int{9},
			indices: []int{80},
			want:    "r9c9=9",
		},
		{
			name:    "missing values",
			values:  nil,
			indices: []int{0},
			want:    "",
		},
		{
			name:    "missing indices",
			values:  []int{5},
			indices: nil,
			want:    "",
		},
		{
			name:    "missing both",
			values:  nil,
			indices: nil,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			step := &SolutionStep{
				values:  tt.values,
				indices: tt.indices,
			}
			if got := step.formatPlacedValue(); got != tt.want {
				t.Errorf("SolutionStep.formatPlacedValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatEmptyRectangle(t *testing.T) {
	// Initialize a mock SolutionStep for Empty Rectangle
	step := NewStep(kindEmptyRectangle).
		WithValues(5).
		WithHouse(&House{Kind: kindBox, Index: 0}). // Box 1
		WithIndices(0, 1)                           // e.g. r1c1, r1c2

	step.DeleteCandidate(2, 5) // e.g. r1c3

	expected := "5 in b1 (r1c12) => r1c3<>5"
	actual := step.formatEmptyRectangle()

	if actual != expected {
		t.Errorf("Expected %q, got %q", expected, actual)
	}
}

func TestSolverFormatStep(t *testing.T) {
	p := puzzle.NewPuzzle()
	s := NewSolver(p, nil)

	tests := []struct {
		name     string
		step     *SolutionStep
		expected string
	}{
		{
			name: "Naked Single",
			step: NewStep(kindNakedSingle).
				WithPlacedValue(0, 5),
			expected: "Naked Single: r1c1=5",
		},
		{
			name: "Hidden Single",
			step: NewStep(kindHiddenSingle).
				WithPlacedValue(0, 5).
				WithHouse(NewHouse(kindRow, 0)),
			expected: "Hidden Single: 5 in r1 => r1c1=5",
		},
		{
			name: "Naked Pair",
			step: func() *SolutionStep {
				step := NewStep(kindNakedPair).
					WithIndices(0, 1).
					WithValues(1, 2)
				step.DeleteCandidate(2, 1)
				return step
			}(),
			expected: "Naked Pair: 1,2 in r1c12 => r1c3<>1",
		},
		{
			name: "Locked Candidates Pointing",
			step: func() *SolutionStep {
				step := NewStep(kindLockedCandidatesPointing).
					WithHouse(NewHouse(kindBox, 0)).
					WithValues(3)
				step.DeleteCandidate(3, 3)
				return step
			}(),
			expected: "Locked Candidates Type 1 (Pointing): 3 in b1 => r1c4<>3",
		},
		{
			name: "X-Wing",
			step: func() *SolutionStep {
				base1 := NewHouse(kindRow, 0)
				base2 := NewHouse(kindRow, 8)
				cover1 := NewHouse(kindColumn, 0)
				cover2 := NewHouse(kindColumn, 8)
				step := NewStep(kindXWing).
					WithValues(4).
					WithBases(base1, base2).
					WithCovers(cover1, cover2)
				step.DeleteCandidate(9, 4) // r2c1<>4
				return step
			}(),
			expected: "X-Wing: 4 r19 c19 => r2c1<>4",
		},
		{
			name: "XY-Wing",
			step: func() *SolutionStep {
				step := NewStep(kindXYWing).
					WithIndices(0, 2, 18). // r1c1, r1c3, r3c1
					WithValues(1, 2, 3)    // Z is 3, eliminated val
				step.DeleteCandidate(20, 3) // r3c3<>3
				return step
			}(),
			expected: "XY-Wing: 1/2/3 in r1c13,r3c1 => r3c3<>3",
		},
		{
			name: "XYZ-Wing",
			step: func() *SolutionStep {
				step := NewStep(kindXYZWing).
					WithIndices(0, 2, 18).
					WithValues(1, 2, 3)
				step.DeleteCandidate(20, 3)
				return step
			}(),
			expected: "XYZ-Wing: 1/2/3 in r1c13,r3c1 => r3c3<>3",
		},
		{
			name: "Skyscraper",
			step: func() *SolutionStep {
				step := NewStep(kindSkyscraper).
					WithValues(5).
					WithIndices(0, 8, 2, 6) // tops (0,8), connectors (2,6)
				step.DeleteCandidate(4, 5) // r1c5<>5
				return step
			}(),
			expected: "Skyscraper: 5 in r1c19 (connected by r1c37) => r1c5<>5",
		},
		{
			name: "Unique Rectangle Type 1",
			step: func() *SolutionStep {
				step := NewStep(kindUniqueRectangle1).
					WithIndices(0, 1, 9, 10). // r12c12
					WithValues(6, 7)
				step.DeleteCandidate(10, 6) // r2c2<>6
				return step
			}(),
			expected: "Unique Rectangle Type 1: 6/7 in r12c12 => r2c2<>6",
		},
		{
			name: "Remote Pair (Fallback to Generic)",
			step: func() *SolutionStep {
				step := NewStep(kindRemotePair)
				step.DeleteCandidate(0, 8)
				return step
			}(),
			expected: "Remote Pair: r1c1<>8",
		},
		{
			name: "Empty Rectangle",
			step: func() *SolutionStep {
				step := NewStep(kindEmptyRectangle).
					WithValues(5).
					WithHouse(&House{Kind: kindBox, Index: 0}). // Box 1
					WithIndices(0, 1)                           // e.g. r1c1, r1c2
				step.DeleteCandidate(2, 5) // e.g. r1c3
				return step
			}(),
			expected: "Empty Rectangle: 5 in b1 (r1c12) => r1c3<>5",
		},
		{
			name: "Brute Force",
			step: NewStep(kindBruteForce).
				WithPlacedValue(80, 9),
			expected: "Brute Force: r9c9=9",
		},
		{
			name: "Format Rect Indices with 3 cells (UR corner)",
			step: func() *SolutionStep {
				step := NewStep(kindUniqueRectangle1).
					WithIndices(0, 1, 9). // r1c1, r1c2, r2c1 (missing r2c2)
					WithValues(6, 7)
				step.DeleteCandidate(10, 6)
				return step
			}(),
			expected: "Unique Rectangle Type 1: 6/7 in r12c12 => r2c2<>6",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := s.FormatStep(tc.step)
			if result != tc.expected {
				t.Errorf("FormatStep() mismatch:\nExpected: %s\nGot:      %s", tc.expected, result)
			}
		})
	}
}
