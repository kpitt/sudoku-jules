package solver

func (s *Solver) findSkyscraper() (found bool) {
	check := func(baseLines []*House) bool {
		return s.checkSkyscraper(baseLines)
	}
	// We check for Skyscrapers where the base lines are rows (meaning the strong
	// links are horizontal) and then where the base lines are columns (vertical
	// strong links).
	return check(s.rows[:]) || check(s.columns[:])
}

func (s *Solver) checkSkyscraper(baseLines []*House) (found bool) {
	candidates := make([]*House, 0, 9)

	for val := 1; val <= 9; val++ {
		// Find all lines that have exactly 2 candidates for this value.
		candidates = candidates[:0]
		for _, line := range baseLines {
			if line.NumLocations(val) == 2 {
				candidates = append(candidates, line)
			}
		}

		if len(candidates) < 2 {
			continue
		}

		// Check each pair of candidate lines to see if they form a Skyscraper.
		for i := 0; i < len(candidates)-1; i++ {
			base1 := candidates[i]
			for j := i + 1; j < len(candidates); j++ {
				base2 := candidates[j]

				// To form a Skyscraper, the two base lines must share exactly
				// one column (or row, if bases are columns) where the candidate
				// appears. This shared location forms the base of the Skyscraper.
				locs1 := base1.Unsolved[val]
				locs2 := base2.Unsolved[val]

				// Find intersection of locations.
				// Since we know size is 2, we can just check values.
				commonLoc := -1
				top1Loc := -1
				top2Loc := -1

				// Identify common and distinct locations
				for _, l1 := range locs1.Values() {
					if locs2.Contains(l1) {
						commonLoc = l1
					} else {
						top1Loc = l1
					}
				}

				// If we didn't find exactly one common location, or if the
				// locations are identical (2 common locations), then this isn't
				// a Skyscraper.
				// Note: If locs1 and locs2 are identical, commonLoc would be set
				// but top1Loc would remain -1.
				if commonLoc == -1 || top1Loc == -1 {
					continue
				}

				// Find top2Loc (the one in base2 that isn't common)
				for _, l2 := range locs2.Values() {
					if l2 != commonLoc {
						top2Loc = l2
						break
					}
				}

				// The "top" cells are the ends of the strong links that are NOT
				// the shared base.
				top1 := base1.Cells[top1Loc]
				top2 := base2.Cells[top2Loc]

				// The "floor" cells are the shared base cells.
				floor1 := base1.Cells[commonLoc]
				floor2 := base2.Cells[commonLoc]

				// Try to eliminate candidates from any cell that sees both tops.
				step := NewStep(kindSkyscraper).WithValues(val)
				if s.eliminateFromIntersection(top1.Index(), top2.Index(), -1, val, step) {
					s.applyStep(step.WithIndices(top1.Index(), top2.Index(), floor1.Index(), floor2.Index()))
					return true
				}
			}
		}
	}

	return false
}

func (s *Solver) findTwoStringKite() (found bool) {
	candidates := make([]*House, 0, 9)
	for val := 1; val <= 9; val++ {
		// Find rows and columns with exactly 2 candidates.
		rowCandidates := candidates[:0] // Reuse slice
		for _, row := range s.rows {
			if row.NumLocations(val) == 2 {
				rowCandidates = append(rowCandidates, row)
			}
		}
		// We can't reuse the same slice for cols since we iterate nested.
		// So let's just make a new one or manage it carefully.
		// Actually, let's just filter locally in the loop.
		// Or loop rows then loop cols.

		if len(rowCandidates) == 0 {
			continue
		}

		colCandidates := make([]*House, 0, 9)
		for _, col := range s.columns {
			if col.NumLocations(val) == 2 {
				colCandidates = append(colCandidates, col)
			}
		}

		if len(colCandidates) == 0 {
			continue
		}

		// Check each Row against each Col
		for _, row := range rowCandidates {
			for _, col := range colCandidates {
				if s.checkTwoStringKite(val, row, col) {
					return true
				}
			}
		}
	}
	return false
}

func (s *Solver) checkTwoStringKite(val int, row, col *House) (found bool) {
	// Row candidates (2 of them)
	rLocs := row.Unsolved[val].Values()
	// Col candidates (2 of them)
	cLocs := col.Unsolved[val].Values()

	// rLocs are column indices in the row.
	// cLocs are row indices in the column.

	checkPoly := func(rP, rTail, cP, cTail int) bool {
		// rP is column index of the row candidate involved in the connection
		// cP is row index of the col candidate involved in the connection

		rCell := row.Cells[rP]
		cCell := col.Cells[cP]

		// Ensure they are not the same cell (intersection)
		if rCell.Index() == cCell.Index() {
			return false
		}

		// Check if they are in the same box (connection)
		if rCell.Box() == cCell.Box() {
			// Found connection!

			// Tails are the other candidates
			tail1 := row.Cells[rTail]
			tail2 := col.Cells[cTail]

			step := NewStep(kindTwoStringKite).
				WithValues(val).
				WithBases(row, col)

			if s.eliminateFromIntersection(tail1.Index(), tail2.Index(), -1, val, step) {
				s.applyStep(step.
					WithIndices(tail1.Index(), tail2.Index(), rCell.Index(), cCell.Index()))
				// Note: visualization usually highlights tails and connection cells.
				// We pass all 4.
				return true
			}
		}
		return false
	}

	// Try all combinations
	// Pair 1: (row[0], col[0]) connect?
	if checkPoly(rLocs[0], rLocs[1], cLocs[0], cLocs[1]) {
		return true
	}
	// Pair 2: (row[0], col[1]) connect?
	if checkPoly(rLocs[0], rLocs[1], cLocs[1], cLocs[0]) {
		return true
	}
	// Pair 3: (row[1], col[0]) connect?
	if checkPoly(rLocs[1], rLocs[0], cLocs[0], cLocs[1]) {
		return true
	}
	// Pair 4: (row[1], col[1]) connect?
	if checkPoly(rLocs[1], rLocs[0], cLocs[1], cLocs[0]) {
		return true
	}

	return false
}
