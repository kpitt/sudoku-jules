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
