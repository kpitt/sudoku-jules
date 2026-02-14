package puzzle

import "testing"

func TestPeerLookup(t *testing.T) {
	// Pick a few cells and verify their peers.
	// Index 0 (r0c0) should have peers:
	// Row 0: 1,2,3,4,5,6,7,8
	// Col 0: 9,18,27,36,45,54,63,72
	// Box 0: 10,11,19,20 (since 0,1,2,9,10,11,18,19,20 are in box 0, and 1,2,9,18 are already row/col peers)
	// Wait, let's just count unique peers.
	// Row peers: 8
	// Col peers: 8
	// Box peers: 4 (remaining)
	// Total: 20

	peers := GetPeers(0)
	if len(peers) != 20 {
		t.Errorf("Expected 20 peers, got %d", len(peers))
	}

	expectedPeers := map[int]bool{
		// Row 0
		1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true,
		// Col 0
		9: true, 18: true, 27: true, 36: true, 45: true, 54: true, 63: true, 72: true,
		// Box 0 (excluding row 0 and col 0)
		10: true, 11: true,
		19: true, 20: true,
	}

	for _, p := range peers {
		if !expectedPeers[p] {
			t.Errorf("Unexpected peer %d for cell 0", p)
		}
		delete(expectedPeers, p)
	}
	if len(expectedPeers) > 0 {
		t.Errorf("Missing peers for cell 0: %v", expectedPeers)
	}
}

func TestHouseLookup(t *testing.T) {
	// Helper to check if a slice contains a value
	contains := func(slice *[9]int, val int) bool {
		for _, v := range slice {
			if v == val {
				return true
			}
		}
		return false
	}

	// Test a Row (House 0)
	row0 := GetHouse(0) // Row 0
	for c := 0; c < 9; c++ {
		idx := 0*9 + c
		if !contains(row0, idx) {
			t.Errorf("Row 0 missing cell %d", idx)
		}
	}

	// Test a Column (House 9)
	col0 := GetHouse(9) // Col 0 is 9th house (0-8 are rows)
	for r := 0; r < 9; r++ {
		idx := r*9 + 0
		if !contains(col0, idx) {
			t.Errorf("Col 0 missing cell %d", idx)
		}
	}

	// Test a Box (House 18)
	box0 := GetHouse(18) // Box 0 is 18th house (0-8 rows, 9-17 cols)
	// Box 0 cells: r0-2, c0-2
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			idx := r*9 + c
			if !contains(box0, idx) {
				t.Errorf("Box 0 missing cell %d", idx)
			}
		}
	}
}
