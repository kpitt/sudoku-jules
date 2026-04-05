package solver

import (
	"testing"
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
