package solver

import (
	"reflect"
	"testing"

	"github.com/kpitt/sudoku/internal/bitset"
	"github.com/kpitt/sudoku/internal/puzzle"
)

func TestNewHouse(t *testing.T) {
	kind := kindRow
	index := 5
	h := NewHouse(kind, index)

	if h.Kind != kind {
		t.Errorf("Expected Kind %v, got %v", kind, h.Kind)
	}
	if h.Index != index {
		t.Errorf("Expected Index %d, got %d", index, h.Index)
	}

	for i := 1; i <= 9; i++ {
		expected := bitset.BitSet16(allLocBits)
		if h.Unsolved[i] != expected {
			t.Errorf("Expected Unsolved[%d] to be %v, got %v", i, expected, h.Unsolved[i])
		}
	}
}

func TestHouse_NumUnsolved(t *testing.T) {
	h := NewHouse(kindRow, 0)

	// Initially, all 9 digits should have possible locations
	if got := h.NumUnsolved(); got != 9 {
		t.Errorf("NewHouse().NumUnsolved() = %v, want %v", got, 9)
	}

	// Clear candidates for digit 1
	h.Unsolved[1].Clear()
	if got := h.NumUnsolved(); got != 8 {
		t.Errorf("After clearing 1 digit, NumUnsolved() = %v, want %v", got, 8)
	}

	// Clear remaining candidates
	for i := 2; i <= 9; i++ {
		h.Unsolved[i].Clear()
	}
	if got := h.NumUnsolved(); got != 0 {
		t.Errorf("After clearing all digits, NumUnsolved() = %v, want %v", got, 0)
	}
}

func TestHouse_UnsolvedDigits(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*House)
		wantDigits []int
		wantNum    int
	}{
		{
			name: "all unsolved",
			setup: func(h *House) {
				// NewHouse already has all 1-9 unsolved by default.
			},
			wantDigits: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantNum:    9,
		},
		{
			name: "some solved",
			setup: func(h *House) {
				h.Unsolved[1].Clear()
				h.Unsolved[4].Clear()
				h.Unsolved[9].Clear()
			},
			wantDigits: []int{2, 3, 5, 6, 7, 8},
			wantNum:    6,
		},
		{
			name: "all solved",
			setup: func(h *House) {
				for i := 1; i <= 9; i++ {
					h.Unsolved[i].Clear()
				}
			},
			wantDigits: []int{},
			wantNum:    0,
		},
		{
			name: "only one unsolved",
			setup: func(h *House) {
				for i := 1; i <= 9; i++ {
					if i != 5 {
						h.Unsolved[i].Clear()
					}
				}
			},
			wantDigits: []int{5},
			wantNum:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHouse(kindRow, 0)
			tt.setup(h)

			if got := h.UnsolvedDigits(); !reflect.DeepEqual(got, tt.wantDigits) {
				t.Errorf("House.UnsolvedDigits() = %v, want %v", got, tt.wantDigits)
			}
			if got := h.NumUnsolved(); got != tt.wantNum {
				t.Errorf("House.NumUnsolved() = %v, want %v", got, tt.wantNum)
			}
		})
	}
}

func TestHouse_UnsolvedStatistics(t *testing.T) {
	h := NewHouse(kindRow, 0)

	// Initially all 9 digits are unsolved
	if got := h.NumUnsolved(); got != 9 {
		t.Errorf("NumUnsolved() = %d, want 9", got)
	}

	digits := h.UnsolvedDigits()
	if len(digits) != 9 {
		t.Errorf("UnsolvedDigits() length = %d, want 9", len(digits))
	}

	// Remove one digit entirely
	h.RemoveCandidateValue(1, 0)
	if got := h.NumUnsolved(); got != 8 {
		t.Errorf("NumUnsolved() = %d, want 8", got)
	}

	// Check NumLocations and Locations for a specific value
	val := 2
	// After RemoveCandidateValue(1, 0), location 0 is removed for all values.
	// Initial bits: 0111111111 (9 locations: 0-8)
	// After removing loc 0: 0111111110 (8 locations: 1-8)
	if got := h.NumLocations(val); got != 8 {
		t.Errorf("NumLocations(%d) = %d, want 8", val, got)
	}

	locs := h.Locations(val)
	if locs.Size() != 8 {
		t.Errorf("Locations(%d).Size() = %d, want 8", val, locs.Size())
	}
	if locs.Contains(0) {
		t.Errorf("Locations(%d) should NOT contain 0", val)
	}
}

func TestHouse_Shared(t *testing.T) {
	p := puzzle.NewPuzzle()
	// Box 0: (0,0)-(2,2)
	// (0,0) index 0, loc 0 in Box 0
	// (0,1) index 1, loc 1 in Box 0
	// (0,2) index 2, loc 2 in Box 0
	// (1,0) index 9, loc 3 in Box 0
	// (1,1) index 10, loc 4 in Box 0
	// (1,2) index 11, loc 5 in Box 0
	// (2,0) index 18, loc 6 in Box 0
	// (2,1) index 19, loc 7 in Box 0
	// (2,2) index 20, loc 8 in Box 0

	h := NewHouse(kindBox, 0)
	for i := 0; i < 9; i++ {
		r, c := (i / 3), (i % 3)
		h.Cells[i] = p.Get(r, c)
	}

	t.Run("sharedRow", func(t *testing.T) {
		// Locs 0, 1, 2 are all in Row 0
		locs := bitset.FromValues16(0, 1, 2)
		row, ok := h.sharedRow(locs)
		if !ok || row != 0 {
			t.Errorf("sharedRow(0,1,2) = (%d, %v), want (0, true)", row, ok)
		}

		// Locs 0, 3 are in Row 0 and Row 1
		locs = bitset.FromValues16(0, 3)
		_, ok = h.sharedRow(locs)
		if ok {
			t.Errorf("sharedRow(0,3) should be false")
		}
	})

	t.Run("sharedCol", func(t *testing.T) {
		// Locs 0, 3, 6 are all in Col 0
		locs := bitset.FromValues16(0, 3, 6)
		col, ok := h.sharedCol(locs)
		if !ok || col != 0 {
			t.Errorf("sharedCol(0,3,6) = (%d, %v), want (0, true)", col, ok)
		}

		// Locs 0, 1 are in Col 0 and Col 1
		locs = bitset.FromValues16(0, 1)
		_, ok = h.sharedCol(locs)
		if ok {
			t.Errorf("sharedCol(0,1) should be false")
		}
	})

	t.Run("sharedBox", func(t *testing.T) {
		// All locs in h (Box 0) should be in Box 0
		locs := bitset.FromValues16(0, 1, 4, 8)
		box, ok := h.sharedBox(locs)
		if !ok || box != 0 {
			t.Errorf("sharedBox(0,1,4,8) = (%d, %v), want (0, true)", box, ok)
		}

		// Test with mixed boxes (needs a different setup or fake cells)
		// but since h.Cells only contains cells from Box 0, sharedBox should always return 0, true if locs are valid.
		// Let's manually put a cell from another box.
		h.Cells[8] = p.Get(0, 3) // Col 3 is Box 1
		locs = bitset.FromValues16(0, 8)
		_, ok = h.sharedBox(locs)
		if ok {
			t.Errorf("sharedBox(0,8) with mixed cells should be false")
		}
	})
}

func TestHouse_TransformLocs(t *testing.T) {
	p := puzzle.NewPuzzle()
	h := NewHouse(kindRow, 0) // Let's say Row 0
	for i := 0; i < 9; i++ {
		h.Cells[i] = p.Get(0, i)
	}

	locs := []int{0, 4, 8}

	t.Run("cellsFromLocs", func(t *testing.T) {
		cells := h.cellsFromLocs(locs)
		if len(cells) != 3 {
			t.Fatalf("cellsFromLocs len = %d, want 3", len(cells))
		}
		if cells[0] != h.Cells[0] || cells[1] != h.Cells[4] || cells[2] != h.Cells[8] {
			t.Errorf("cellsFromLocs returned wrong cells")
		}
	})

	t.Run("indexesFromLocs", func(t *testing.T) {
		indexes := h.indexesFromLocs(locs)
		if len(indexes) != 3 {
			t.Fatalf("indexesFromLocs len = %d, want 3", len(indexes))
		}
		// Row 0, cols 0, 4, 8 should have indices 0, 4, 8
		if indexes[0] != 0 || indexes[1] != 4 || indexes[2] != 8 {
			t.Errorf("indexesFromLocs returned wrong indexes: %v", indexes)
		}
	})
}

func TestHouse_RemoveCandidateLoc(t *testing.T) {
	h := NewHouse(kindRow, 0)
	val := 1
	loc := 3

	if !h.Unsolved[val].Contains(loc) {
		t.Fatalf("Expected Unsolved[%d] to contain location %d after initialization", val, loc)
	}

	h.RemoveCandidateLoc(val, loc)

	if h.Unsolved[val].Contains(loc) {
		t.Errorf("Expected Unsolved[%d] to NOT contain location %d after removal", val, loc)
	}
}

func TestHouse_RemoveCandidateValue(t *testing.T) {
	h := NewHouse(kindRow, 0)
	val := 5
	loc := 4

	h.RemoveCandidateValue(val, loc)

	// val should be cleared from unsolved
	if !h.Unsolved[val].Empty() {
		t.Errorf("Expected Unsolved[%d] to be empty after RemoveCandidateValue, got %v", val, h.Unsolved[val])
	}

	// loc should be removed from all other values
	for i := 1; i <= 9; i++ {
		if i == val {
			continue
		}
		if h.Unsolved[i].Contains(loc) {
			t.Errorf("Expected Unsolved[%d] to NOT contain location %d after RemoveCandidateValue", i, loc)
		}
	}
}

func TestHouse_Name(t *testing.T) {
	tests := []struct {
		kind houseKind
		want string
	}{
		{kindRow, "Row"},
		{kindColumn, "Column"},
		{kindBox, "Box"},
	}

	for _, tt := range tests {
		h := NewHouse(tt.kind, 0)
		if got := h.Name(); got != tt.want {
			t.Errorf("House.Name() for kind %v = %v, want %v", tt.kind, got, tt.want)
		}
	}
}
