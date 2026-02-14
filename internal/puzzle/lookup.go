package puzzle

// PeerLookup is a pre-computed table of the peers for each cell index.
// A peer is a cell that "sees" the cell at the given index, and therefore
// cannot contain the same value.
// Each cell has 20 peers: 8 cells in the same row, 8 cells in the same column,
// and 4 additional cells in the same box.
var PeerLookup [81][20]int

// HouseLookup is a pre-computed table of cell indices for each house.
// There are 27 houses:
// - Indices 0-8: Rows 0-8
// - Indices 9-17: Columns 0-8
// - Indices 18-26: Boxes 0-8
var HouseLookup [27][9]int

func init() {
	InitPeerLookup()
	InitHouseLookup()
}

// InitPeerLookup initializes the PeerLookup table.
func InitPeerLookup() {
	for i := range 81 {
		r, c := i/9, i%9
		boxStartR, boxStartC := (r/3)*3, (c/3)*3
		peerMap := make(map[int]bool)

		// Add Row and Column peers
		for k := range 9 {
			if k != c {
				peerMap[r*9+k] = true // Row
			}
			if k != r {
				peerMap[k*9+c] = true // Col
			}
		}

		// Add Box peers
		for br := range 3 {
			for bc := range 3 {
				peerIdx := (boxStartR+br)*9 + (boxStartC + bc)
				if peerIdx != i {
					peerMap[peerIdx] = true
				}
			}
		}

		idx := 0
		for p := range peerMap {
			PeerLookup[i][idx] = p
			idx++
		}
	}
}

// InitHouseLookup initializes the HouseLookup table.
func InitHouseLookup() {
	var houseIdx int

	// 1. Rows
	for r := range 9 {
		for c := range 9 {
			HouseLookup[houseIdx][c] = r*9 + c
		}
		houseIdx++
	}

	// 2. Columns
	for c := range 9 {
		for r := range 9 {
			HouseLookup[houseIdx][r] = r*9 + c
		}
		houseIdx++
	}

	// 3. Boxes
	for br := range 3 {
		for bc := range 3 {
			// Internal box index
			for i := range 9 {
				r := br*3 + i/3
				c := bc*3 + i%3
				HouseLookup[houseIdx][i] = r*9 + c
			}
			houseIdx++
		}
	}
}

// GetPeers returns the peers of the cell at the given index.
// The `PeerLookup` table is a global array, so we can return a stable pointer
// to avoid copying the peer array.
func GetPeers(idx int) *[20]int {
	return &PeerLookup[idx]
}

// GetHouse returns the cells in the given house.
func GetHouse(idx int) *[9]int {
	return &HouseLookup[idx]
}
