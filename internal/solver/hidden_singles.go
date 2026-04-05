package solver

// findHiddenSingles places the value of any cells that match the "Hidden
// Single" pattern.  A "Hidden Single" is the only cell that contains a
// particular candidate in its house.
func (s *Solver) findHiddenSingles() bool {
	found := false
	for _, h := range s.houses {
		found = s.checkHiddenSinglesForHouse(h) || found
	}
	return found
}

func (s *Solver) checkHiddenSinglesForHouse(h *House) bool {
	for val := 1; val <= 9; val++ {
		locs := h.Unsolved[val]
		if locs.Empty() {
			continue
		}
		if locs.Size() == 1 {
			index := locs.Value()
			cell := h.Cells[index]
			step := NewStep(kindHiddenSingle).
				WithHouse(h).
				WithPlacedValue(cell.Row*9+cell.Col, val)
			s.applyStep(step)
			return true
		}
	}
	return false
}
