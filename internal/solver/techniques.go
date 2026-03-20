package solver

import (
	"slices"
	"time"

	"github.com/kpitt/sudoku/internal/bitset"
	"github.com/kpitt/sudoku/internal/puzzle"
)

type techniqueKind int

const (
	// Enum constants for each solving technique, which must be defined in the
	// order that they should be applied.
	kindFullHouse techniqueKind = iota
	kindNakedSingle
	kindHiddenSingle
	kindLockedCandidatesPointing
	kindLockedCandidatesClaiming
	kindNakedPair
	kindNakedTriple
	kindHiddenPair
	kindHiddenTriple
	kindNakedQuadruple
	kindHiddenQuadruple
	kindXWing
	kindSwordfish
	kindJellyfish
	kindRemotePair
	kindBUG
	kindSkyscraper
	kindTwoStringKite
	kindEmptyRectangle
	kindWWing
	kindXYWing
	kindXYZWing
	kindAvoidableRectangle
	kindUniqueRectangle1
	kindUniqueRectangle2
	kindUniqueRectangle3
	kindUniqueRectangle4
	kindHiddenRectangle
	kindFinnedXWing
	kindFinnedSwordfish
	kindFinnedJellyfish
	kindFrankenXWing
	kindFrankenSwordfish
	kindFrankenJellyfish
	kindMutantXWing
	kindMutantSwordfish
	kindMutantJellyfish
	kindSueDeCoq
	kindSimpleColoring
	kindMultiColoring
	kindXChain
	kindXYChain
	kindNiceLoop
	kindAIC
	kindALSXZ
	kindALSXYWing
	kindALSXYChain
	kindForcingChain
	kindForcingNet
	kindBruteForce
)

// A Technique represents a single Sudoku solving technique, represented by a
// display name and a function to check for the technique in the puzzle.
type Technique struct {
	Name  string
	Check func() bool
}

// initTechniques initializes a list of the known solving techniques and
// assigns it to the solver.
// The order must match the order of the techniqueKind constants.
func (s *Solver) initTechniques() {
	// The order of this list must match the order of the techniqueKind constants.
	s.techniques = []Technique{
		{"Full House", s.findFullHouse},
		{"Naked Single", s.findNakedSingles}, // also checked during candidate removal
		{"Hidden Single", s.findHiddenSingles},
		{"Locked Candidates Type 1 (Pointing)", s.findPointingTuples},
		{"Locked Candidates Type 2 (Claiming)", s.findClaimingTuples},
		{"Naked Pair", s.findNakedPairs},
		{"Naked Triple", s.findNakedTriples},
		{"Hidden Pair", s.findHiddenPairs},
		{"Hidden Triple", s.findHiddenTriples},
		{"Naked Quadruple", s.findNakedQuadruples},
		{"Hidden Quadruple", s.findHiddenQuadruples},
		{"X-Wing", s.findXWings},
		{"Swordfish", s.findSwordfish},
		{"Jellyfish", s.findJellyfish},
		{"Remote Pair", s.findRemotePairs},
		{"BUG+1", s.findBUG},
		{"Skyscraper", s.findSkyscraper},
		{"2-String Kite", s.findTwoStringKite},
		{"Empty Rectangle", s.findEmptyRectangle},
		{"W-Wing", s.findWWing},
		{"XY-Wing", s.findXYWings},
		{"XYZ-Wing", s.findXYZWings},
		{"Avoidable Rectangle", s.findAvoidableRectangles},
		{"Unique Rectangle Type 1", s.findUniqueRectangleType1},
		{"Unique Rectangle Type 2", s.findUniqueRectangleType2},
		{"Unique Rectangle Type 3", s.findUniqueRectangleType3},
		{"Unique Rectangle Type 4", s.findUniqueRectangleType4},
		{"Hidden Rectangle", s.findHiddenRectangle},
		{"Finned X-Wing", s.findFinnedXWings},
		{"Finned Swordfish", s.findFinnedSwordfish},
		{"Finned Jellyfish", s.findFinnedJellyfish},
		{"Franken X-Wing", s.findFrankenXWing},
		{"Franken Swordfish", s.findFrankenSwordfish},
		{"Franken Jellyfish", s.findFrankenJellyfish},
		{"Mutant X-Wing", s.findMutantXWing},
		{"Mutant Swordfish", s.findMutantSwordfish},
		{"Mutant Jellyfish", s.findMutantJellyfish},
		{"Sue de Coq", s.findSueDeCoq},
		{"Simple Coloring", s.findSimpleColoring},
		{"Multi-Coloring", s.findMultiColoring},
		{"X-Chain", s.findXChains},
		{"XY-Chain", s.findXYChains},
		{"Nice Loop", s.findNiceLoops},
		{"AIC", s.findAICs},
		{"ALS-XZ", s.findALSXZ},
		{"ALS-XY-Wing", s.findALSXYWing},
		{"ALS-XY-Chain", s.findALSXYChain},
		{"Forcing Chain", s.findForcingChains},
		{"Forcing Net", s.findForcingNets},
		{"Brute Force", nil}, // custom check as last resort
	}
}

func (s *Solver) findFrankenXWing() bool { return s.findGeneralizedFish(2, kindFrankenXWing, false) }
func (s *Solver) findFrankenSwordfish() bool { return s.findGeneralizedFish(3, kindFrankenSwordfish, false) }
func (s *Solver) findFrankenJellyfish() bool { return s.findGeneralizedFish(4, kindFrankenJellyfish, false) }
func (s *Solver) findMutantXWing() bool { return s.findGeneralizedFish(2, kindMutantXWing, true) }
func (s *Solver) findMutantSwordfish() bool { return s.findGeneralizedFish(3, kindMutantSwordfish, true) }
func (s *Solver) findMutantJellyfish() bool { return s.findGeneralizedFish(4, kindMutantJellyfish, true) }

func (s *Solver) findGeneralizedFish(size int, kind techniqueKind, allowMutant bool) bool {
	globalStep := NewStep(kind)
	foundAny := false

	for val := 1; val <= 9; val++ {
		// Identify potential base houses
		var potentialBases []*House
		for _, h := range s.houses {
			if !allowMutant && h.Kind == kindBox {
				continue
			}
			if h.NumLocations(val) >= 2 {
				potentialBases = append(potentialBases, h)
			}
		}

		if len(potentialBases) < size {
			continue
		}

		// Try all combinations of size base houses
		if s.collectGeneralizedFish(val, size, globalStep, potentialBases) {
			foundAny = true
		}
	}

	if foundAny {
		s.applyStep(globalStep)
		return true
	}
	return false
}

func (s *Solver) collectGeneralizedFish(val, size int, step *SolutionStep, potentialBases []*House) bool {
	found := false
	var check func(start int, current []*House) bool
	check = func(start int, current []*House) bool {
		if len(current) == size {
			if s.checkGeneralizedFishForBaseSetAndStep(val, current, step) {
				found = true
				// We don't return true here because we want to find ALL combinations
				// for this value? Actually, one combination is usually enough for one digit.
				// But let's find all to be safe.
			}
			return false // Keep searching
		}
		for i := start; i < len(potentialBases); i++ {
			if check(i+1, append(current, potentialBases[i])) {
				return true
			}
		}
		return false
	}
	check(0, nil)
	return found
}

func (s *Solver) checkGeneralizedFishForBaseSetAndStep(val int, bases []*House, step *SolutionStep) bool {
	// 1. Get all cell indices in base houses
	var sIdxs bitset.BitSet128
	for _, b := range bases {
		for _, cell := range b.Cells {
			if cell.HasCandidate(val) {
				sIdxs.Add(cell.Index())
			}
		}
	}
	if sIdxs.Size() == 0 {
		return false
	}

	// 2. Identify potential cover houses
	potentialCovers := make(map[int]*House)
	for idx := range sIdxs.All() {
		r, c := idx/9, idx%9
		b := s.puzzle.Cell(idx).Box()
		potentialCovers[r] = s.rows[r]
		potentialCovers[9+c] = s.columns[c]
		potentialCovers[18+b] = s.boxes[b]
	}

	var coverList []*House
	for _, h := range potentialCovers {
		coverList = append(coverList, h)
	}

	// 3. Try to find size cover houses that cover all sIdxs
	size := len(bases)
	var checkCover func(start int, current []*House, covered bitset.BitSet128) bool
	checkCover = func(start int, current []*House, covered bitset.BitSet128) bool {
		if covered.Size() == sIdxs.Size() {
			if len(current) != size { return false }

			foundElim := false
			for _, ch := range current {
				for _, cell := range ch.Cells {
					if !sIdxs.Contains(cell.Index()) && cell.HasCandidate(val) {
						step.DeleteCandidate(cell.Index(), val)
						foundElim = true
					}
				}
			}

			if foundElim {
				step.WithValues(val).WithBases(bases...).WithCovers(current...)
				return true
			}
			return false
		}
		
		if len(current) == size { return false }

		for i := start; i < len(coverList); i++ {
			h := coverList[i]
			newCovered := covered
			hasNew := false
			for _, cell := range h.Cells {
				if sIdxs.Contains(cell.Index()) && !covered.Contains(cell.Index()) {
					newCovered.Add(cell.Index())
					hasNew = true
				}
			}
			if !hasNew { continue }

			if checkCover(i+1, append(current, h), newCovered) {
				return true
			}
		}
		return false
	}

	return checkCover(0, nil, bitset.BitSet128{})
}

func sizeOfBases(bases []*House) int { return len(bases) }

func (s *Solver) findFinnedSwordfish() bool {
	return s.findFinnedFishOfSize(3, kindFinnedSwordfish)
}

func (s *Solver) findFinnedJellyfish() bool {
	return s.findFinnedFishOfSize(4, kindFinnedJellyfish)
}

func (s *Solver) findFinnedFishOfSize(fishSize int, fishKind techniqueKind) bool {
	find := func(baseLines, coverLines []*House) bool {
		return s.findFinnedFishInLines(fishSize, fishKind, baseLines, coverLines)
	}
	return find(s.rows, s.columns) || find(s.columns, s.rows)
}

func (s *Solver) findFinnedFishInLines(
	fishSize int,
	fishKind techniqueKind,
	baseLines, coverLines []*House,
) bool {
	// A finned fish of size N has:
	// - N base lines
	// - N cover lines
	// - All candidates in base lines are covered by cover lines, EXCEPT for
	//   one or more "fin" cells.
	// - All fin cells must be in the same box.
	// - The box must also contain the intersection of one base line and one cover line.

	for val := 1; val <= 9; val++ {
		// Identify potential base lines: any line with at least 2 candidates.
		var candidates []*House
		for _, h := range baseLines {
			if h.NumLocations(val) >= 2 {
				candidates = append(candidates, h)
			}
		}

		if len(candidates) < fishSize {
			continue
		}

		// Try all combinations of fishSize base lines.
		if s.checkFinnedFishCombinations(val, fishSize, fishKind, candidates, coverLines) {
			return true
		}
	}
	return false
}

func (s *Solver) checkFinnedFishCombinations(
	val, fishSize int,
	fishKind techniqueKind,
	candidates []*House,
	coverLines []*House,
) bool {
	var check func(start int, current []*House) bool
	check = func(start int, current []*House) bool {
		if len(current) == fishSize {
			return s.checkFinnedFishForBaseSet(val, fishKind, current, coverLines)
		}

		for i := start; i < len(candidates); i++ {
			if check(i+1, append(current, candidates[i])) {
				return true
			}
		}
		return false
	}

	return check(0, nil)
}

func (s *Solver) checkFinnedFishForBaseSet(
	val int,
	fishKind techniqueKind,
	bases []*House,
	coverLines []*House,
) bool {
	// 1. Collect all candidate locations across all base lines.
	// These are indices into coverLines.
	allLocs := bitset.BitSet16(0)
	for _, b := range bases {
		allLocs = bitset.Union(allLocs, b.Unsolved[val])
	}

	// 2. We need to find N cover lines that cover ALMOST all of these locations.
	locs := allLocs.Values()
	if len(locs) < len(bases) {
		return false
	}

	// Try each subset of fishSize locations as the "cover set".
	var checkCoverSet func(start int, currentLocs []int) bool
	checkCoverSet = func(start int, currentLocs []int) bool {
		if len(currentLocs) == len(bases) {
			coverSet := bitset.FromValues16(currentLocs...)
			fins := allLocs.Difference(coverSet)
			if fins.Empty() {
				return false // Perfect fish, skip here.
			}

			// 3. All fins must be in the same box.
			boxIdx := -1
			for _, finLoc := range fins.Values() {
				// We need to find WHICH base line this fin is in.
				for _, b := range bases {
					if b.Unsolved[val].Contains(finLoc) {
						cell := b.Cells[finLoc]
						if boxIdx == -1 {
							boxIdx = cell.Box()
						} else if boxIdx != cell.Box() {
							return false // Fins in multiple boxes
						}
					}
				}
			}

			if boxIdx == -1 {
				return false
			}

			// 4. The box must contain at least one of the "fish" cells that
			// is on a cover line AND in one of the base lines.
			// Actually, the box must contain the intersection of the "finned" base line
			// and the cover line that it "would have" had.
			// More generally, for Finned Fish, the eliminations are in the box
			// and on the cover lines.
			
			// Let's identify the potential eliminations.
			step := NewStep(fishKind).WithValues(val)
			foundElim := false
			box := s.boxes[boxIdx]
			
			// For each cover line in our cover set:
			for _, cIdx := range currentLocs {
				cover := coverLines[cIdx]
				
				// Candidates in this cover line, in this box, that are NOT in any base line.
				for _, cell := range box.Cells {
					// Is this cell on the cover line?
					onCover := false
					if cover.Kind == kindRow && cell.Row == cover.Index {
						onCover = true
					} else if cover.Kind == kindColumn && cell.Col == cover.Index {
						onCover = true
					}
					
					if !onCover {
						continue
					}
					
					// Is it in a base line?
					inBase := false
					for _, b := range bases {
						if b.Kind == kindRow && cell.Row == b.Index {
							inBase = true
							break
						} else if b.Kind == kindColumn && cell.Col == b.Index {
							inBase = true
							break
						}
					}
					
					if inBase {
						continue
					}
					
					if cell.HasCandidate(val) {
						step.DeleteCandidate(cell.Index(), val)
						foundElim = true
					}
				}
			}

			if foundElim {
				covers := transformSlice(currentLocs, func(idx int) *House {
					return coverLines[idx]
				})
				// Collect all fish/fin indices for highlighting.
				indices := bitset.BitSet16(0)
				for _, b := range bases {
					for _, l := range b.Unsolved[val].Values() {
						indices.Add(b.Cells[l].Index())
					}
				}

				s.applyStep(step.
					WithBases(bases...).
					WithCovers(covers...).
					WithIndices(indices.Values()...))
				return true
			}
			return false
		}

		for i := start; i < len(locs); i++ {
			if checkCoverSet(i+1, append(currentLocs, locs[i])) {
				return true
			}
		}
		return false
	}

	return checkCoverSet(0, nil)
}


func (s *Solver) findFullHouse() bool {
	for _, h := range s.houses {
		if h.NumUnsolved() == 1 {
			// Find the unsolved cell
			for _, cell := range h.Cells {
				if !cell.IsSolved() {
					val := cell.CandidateValues()[0]
					step := NewStep(kindFullHouse).
						WithHouse(h).
						WithPlacedValue(cell.Index(), val)
					s.applyStep(step)
					return true
				}
			}
		}
	}
	return false
}

func (s *Solver) findNakedSingles() bool {
	for i := range 81 {
		cell := s.puzzle.Cell(i)
		if !cell.IsSolved() && cell.NumCandidates() == 1 {
			val := cell.CandidateValues()[0]
			step := NewStep(kindNakedSingle).
				WithPlacedValue(i, val)
			s.applyStep(step)
			return true
		}
	}
	return false
}

// ***** IMPORTANT NOTE *****
//
// When processing a check against a set of viable candidates, _always_
// short-circuit the remaining checks after making a change that could
// invalidate the remaining candidates.  Checks should not be combined in a
// single pass unless the checks are completely independent.  If it isn't
// clear whether or not the checks are independent, go ahead and short-ciruit,
// and any additional candidates will get checked in the next solver pass.
//
// For checks that should _not_ short-circuit, be careful when using OR
// expressions to combine the results.  The safest approach is to use an
// accumulator variable and the pattern `found = check(...) || found` for each
// check.  The `||` operator will short-circuit after the first term that
// evalues to `true`, so only the first term is guaranteed to be evaluated.

// findHiddenSingles places the value of any cells that match the "Hidden
// Single" pattern.  A "Hidden Single" is the only cell that contains a
// particular candidate in its house.
func (s *Solver) findHiddenSingles() bool {
	for _, h := range s.houses {
		if s.checkHiddenSinglesForHouse(h) {
			return true
		}
	}
	return false
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

func (s *Solver) findNakedSubsets(size int, kind techniqueKind) (found bool) {
	return slices.ContainsFunc(s.houses, func(h *House) bool {
		return s.checkNakedSubsetsForHouse(size, kind, h)
	})
}

func (s *Solver) checkNakedSubsetsForHouse(size int, kind techniqueKind, h *House) (found bool) {
	var locsBuf [9]int
	locs := locsBuf[:0]
	// Collect a list of all locations with no more than `size` candidate values.
	for i, c := range h.Cells {
		numCand := c.NumCandidates()
		if numCand >= 2 && numCand <= size {
			locs = append(locs, i)
		}
	}
	if len(locs) < size {
		// We need at least `size` candidate values to have a subset of that size.
		return false
	}

	// Try combinations of the required size.
	var checkCombinations func(start int, currentIndices []int) bool
	checkCombinations = func(start int, currentIndices []int) bool {
		if len(currentIndices) == size {
			valueSet := bitset.BitSet16(0)
			for _, idx := range currentIndices {
				valueSet = bitset.Union(valueSet, h.Cells[locs[idx]].Candidates)
			}

			if valueSet.Size() == size {
				var locSet bitset.BitSet16
				for _, idx := range currentIndices {
					locSet.Add(locs[idx])
				}

				step := NewStep(kind)
				if s.eliminateFromOtherLocs(h, valueSet, locSet, step) {
					s.applyStep(step.
						WithIndices(h.indexesFromLocs(locSet.Values())...).
						WithValues(valueSet.Values()...).
						WithHouse(h))
					return true
				}
			}
			return false
		}

		for i := start; i < len(locs); i++ {
			if checkCombinations(i+1, append(currentIndices, i)) {
				return true
			}
		}
		return false
	}

	// Avoid allocation for small known sizes by allocating the initial buffer on the stack.
	indices := make([]int, 0, size)
	return checkCombinations(0, indices)
}

func (s *Solver) findNakedPairs() (found bool) {
	return s.findNakedSubsets(2, kindNakedPair)
}

// eliminateFromOtherLocs removes the candidates listed in values from all
// cells that are not listed in locs.
func (s *Solver) eliminateFromOtherLocs(
	h *House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	found := false
	for l := range 9 {
		if locs.Contains(l) {
			continue
		}
		c := h.Cells[l]

		// Check for intersection first to avoid inner loop/allocations.
		common := c.Candidates.Intersection(values)
		if common.Empty() {
			continue
		}

		for v := range common.All() {
			step.DeleteCandidate(c.Index(), v)
			found = true
		}
	}

	return found
}

// eliminateFromOtherLocsMulti removes the candidates listed in values from all
// cells from each house in houses whose index is not listed in locs.  Returns
// true if at least one candidate was eliminated.
func (s *Solver) eliminateFromOtherLocsMulti(
	houses []*House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	updated := false
	for _, g := range houses {
		updated = s.eliminateFromOtherLocs(g, values, locs, step) || updated
	}

	return updated
}

func (s *Solver) findClaimingTuples() (found bool) {
	for i := range 9 {
		// We only need to check rows and columns for Locked Candidates.
		if s.checkLockedCandidatesForLine(s.rows[i]) ||
			s.checkLockedCandidatesForLine(s.columns[i]) {
			return true
		}
	}
	return false
}

func (s *Solver) checkLockedCandidatesForLine(line *House) (found bool) {
	for val := 1; val <= 9; val++ {
		locs := line.Unsolved[val]
		if locs.Empty() {
			continue
		}

		// If we have more than 3 candidates in a line, then they can't all be
		// in the same box.
		if locs.Size() > 3 {
			continue
		}

		valueSet := bitset.FromValues16(val)
		if box, ok := line.sharedBox(locs); ok {
			cells := line.cellsFromLocs(locs.Values())
			boxCells := transformSlice(cells, func(c *puzzle.Cell) int {
				_, index := getBoxLoc(c.Row, c.Col)
				return index
			})
			locSet := bitset.FromValues16(boxCells...)
			step := NewStep(kindLockedCandidatesClaiming)
			if s.eliminateFromOtherLocs(s.boxes[box], valueSet, locSet, step) {
				s.applyStep(step.
					WithValues(val).
					WithHouse(line))
				return true
			}
		}
	}

	return false
}

func (s *Solver) findPointingTuples() (found bool) {
	for i := range 9 {
		// We only need to check boxes for Pointing Tuples.
		if s.checkPointingTuplesForBox(s.boxes[i]) {
			return true
		}
	}
	return false
}

func (s *Solver) checkPointingTuplesForBox(box *House) (found bool) {
	for val := 1; val <= 9; val++ {
		locs := box.Unsolved[val]
		if locs.Empty() {
			continue
		}

		// If we have more than 3 candidates in a single box, then they can't all
		// be in the same line.
		if locs.Size() > 3 {
			continue
		}

		step := NewStep(kindLockedCandidatesPointing).
			WithValues(val).
			WithHouse(box)
		valueSet := bitset.FromValues16(val)
		cells := box.cellsFromLocs(locs.Values())
		if row, ok := box.sharedRow(locs); ok {
			cols := transformSlice(cells, func(c *puzzle.Cell) int {
				return c.Col
			})
			locSet := bitset.FromValues16(cols...)
			if s.eliminateFromOtherLocs(s.rows[row], valueSet, locSet, step) {
				s.applyStep(step)
				return true
			}
		}
		if col, ok := box.sharedCol(locs); ok {
			rows := transformSlice(cells, func(c *puzzle.Cell) int {
				return c.Row
			})
			locSet := bitset.FromValues16(rows...)
			if s.eliminateFromOtherLocs(s.columns[col], valueSet, locSet, step) {
				s.applyStep(step)
				return true
			}
		}
	}

	return false
}

func (s *Solver) findHiddenSubsets(size int, kind techniqueKind) (found bool) {
	return slices.ContainsFunc(s.houses, func(h *House) bool {
		return s.checkHiddenSubsetsForHouse(size, kind, h)
	})
}

func (s *Solver) checkHiddenSubsetsForHouse(size int, kind techniqueKind, h *House) (found bool) {
	var valBuf [9]int
	values := valBuf[:0]
	for val := 1; val <= 9; val++ {
		locs := h.Unsolved[val]
		if locs.Empty() {
			continue
		}
		if locs.Size() >= 2 && locs.Size() <= size {
			values = append(values, val)
		}
	}

	if len(values) < size {
		// We need at least `size` candidate values to have a subset of that size.
		return false
	}

	var checkCombinations func(start int, currentIndices []int) bool
	checkCombinations = func(start int, currentIndices []int) bool {
		if len(currentIndices) == size {
			locSet := bitset.BitSet16(0)
			valueSet := bitset.BitSet16(0)
			for _, idx := range currentIndices {
				val := values[idx]
				locSet = bitset.Union(locSet, h.Unsolved[val])
				valueSet.Add(val)
			}

			if locSet.Size() == size {
				step := NewStep(kind)
				if s.eliminateOtherValues(h, valueSet, locSet, step) {
					s.applyStep(step.
						WithIndices(h.indexesFromLocs(locSet.Values())...).
						WithValues(valueSet.Values()...).
						WithHouse(h))
					return true
				}
			}
			return false
		}

		for i := start; i < len(values); i++ {
			if checkCombinations(i+1, append(currentIndices, i)) {
				return true
			}
		}
		return false
	}

	indices := make([]int, 0, size)
	return checkCombinations(0, indices)
}

func (s *Solver) findHiddenPairs() (found bool) {
	return s.findHiddenSubsets(2, kindHiddenPair)
}

// eliminateOtherValues removes candidates that are not listed in values from
// the cells in locs.
func (s *Solver) eliminateOtherValues(
	h *House, values ValSet, locs LocSet, step *SolutionStep,
) bool {
	found := false
	for l := range locs.All() {
		c := h.Cells[l]

		// We want to remove any candidates in `c` that are NOT in `values`.
		toRemove := c.Candidates.Difference(values)
		if toRemove.Empty() {
			continue
		}

		for v := range toRemove.All() {
			step.DeleteCandidate(c.Index(), v)
			found = true
		}
	}

	return found
}

func (s *Solver) findNakedTriples() (found bool) {
	return s.findNakedSubsets(3, kindNakedTriple)
}

func (s *Solver) findHiddenTriples() (found bool) {
	return s.findHiddenSubsets(3, kindHiddenTriple)
}

func (s *Solver) findNakedQuadruples() (found bool) {
	return s.findNakedSubsets(4, kindNakedQuadruple)
}

func (s *Solver) findRemotePairs() (found bool) {
	// 1. Group bivalue cells by their candidate pairs.
	groups := make(map[bitset.BitSet16][]int)
	for i := range 81 {
		cell := s.puzzle.Cell(i)
		if cell.NumCandidates() == 2 {
			groups[cell.Candidates] = append(groups[cell.Candidates], i)
		}
	}

	// 2. For each group of at least 4 cells:
	for candidates, indices := range groups {
		if len(indices) < 4 {
			continue
		}

		// 3. Find connected components and their bicoloring.
		visited := make(map[int]bool)
		for _, startIdx := range indices {
			if visited[startIdx] {
				continue
			}

			// Start a new component with BFS
			colors := make(map[int]int) // index -> color (0 or 1)
			queue := []int{startIdx}
			colors[startIdx] = 0
			visited[startIdx] = true

			component := []int{startIdx}
			head := 0
			for head < len(queue) {
				u := queue[head]
				head++

				uColor := colors[u]
				// Neighbors are other indices in this group that see u
				for _, vIdx := range indices {
					if u == vIdx {
						continue
					}
					if s.sees(u, vIdx) {
						if _, ok := colors[vIdx]; !ok {
							colors[vIdx] = 1 - uColor
							visited[vIdx] = true
							queue = append(queue, vIdx)
							component = append(component, vIdx)
						}
					}
				}
			}

			if len(component) < 4 {
				continue
			}

			// 4. For each pair of nodes in the component with different colors:
			// Any cell seeing both can't have candidates X or Y.
			vals := candidates.Values()
			x, y := vals[0], vals[1]

			for i := 0; i < len(component); i++ {
				for j := i + 1; j < len(component); j++ {
					uIdx, vIdx := component[i], component[j]
					if colors[uIdx] == colors[vIdx] {
						continue
					}

					// Different colors -> different values (one is X, other is Y)
					commonPeers := s.commonPeers(uIdx, vIdx)
					step := NewStep(kindRemotePair).
						WithValues(x, y).
						WithIndices(uIdx, vIdx)

					for _, peerIdx := range commonPeers {
						peerCell := s.puzzle.Cell(peerIdx)
						if peerCell.HasCandidate(x) {
							step.DeleteCandidate(peerIdx, x)
						}
						if peerCell.HasCandidate(y) {
							step.DeleteCandidate(peerIdx, y)
						}
					}

					if len(step.deletedCandidates) > 0 {
						s.applyStep(step)
						return true
					}
				}
			}
		}
	}

	return false
}

func (s *Solver) findBUG() (found bool) {
	// BUG + 1 is a state where every unsolved cell has 2 candidates except one,
	// which has 3 candidates.
	var bugCell *puzzle.Cell
	nonBivalueCount := 0
	for i := range 81 {
		cell := s.puzzle.Cell(i)
		if cell.IsSolved() {
			continue
		}
		num := cell.NumCandidates()
		if num == 2 {
			continue
		}
		nonBivalueCount++
		if num == 3 && bugCell == nil {
			bugCell = cell
			continue
		}
	}

	if nonBivalueCount != 1 || bugCell == nil {
		return false
	}

	// In the BUG cell, the value to place is the one that appears 3 times in
	// its row, column, and box houses.
	candidates := bugCell.CandidateValues()
	for _, v := range candidates {
		if s.rows[bugCell.Row].NumLocations(v) == 3 &&
			s.columns[bugCell.Col].NumLocations(v) == 3 &&
			s.boxes[bugCell.Box()].NumLocations(v) == 3 {
			step := NewStep(kindBUG).
				WithPlacedValue(bugCell.Index(), v)
			s.applyStep(step)
			return true
		}
	}

	return false
}

func (s *Solver) findWWing() (found bool) {
	// A W-Wing requires two bivalue cells with the same candidates {x, y}
	// and a strong link on one of those candidates {x or y} in a separate house.
	// If the cells in the strong link each see one of the bivalue cells,
	// then the OTHER candidate can be eliminated from common peers of the
	// bivalue cells.

	// 1. Group bivalue cells by their candidate pairs.
	groups := make(map[bitset.BitSet16][]int)
	for i := range 81 {
		cell := s.puzzle.Cell(i)
		if cell.NumCandidates() == 2 {
			groups[cell.Candidates] = append(groups[cell.Candidates], i)
		}
	}

	for candidates, indices := range groups {
		if len(indices) < 2 {
			continue
		}

		vals := candidates.Values()
		for _, x := range vals {
			y := vals[0]
			if y == x {
				y = vals[1]
			}

			// We look for a strong link on X.
			// 2. Iterate over all houses to find strong links for X.
			for _, h := range s.houses {
				if h.NumLocations(x) == 2 {
					locs := h.Locations(x).Values()
					s1 := h.Cells[locs[0]].Index()
					s2 := h.Cells[locs[1]].Index()

					// 3. For each pair of bivalue cells:
					for i := 0; i < len(indices); i++ {
						for j := i + 1; j < len(indices); j++ {
							w1, w2 := indices[i], indices[j]

							// W1 sees s1 and W2 sees s2 (or vice versa)
							// AND w1, w2 must NOT be s1 or s2.
							if (s.sees(w1, s1) && s.sees(w2, s2)) ||
								(s.sees(w1, s2) && s.sees(w2, s1)) {
								common := s.commonPeers(w1, w2)
								step := NewStep(kindWWing).
									WithValues(x, y).
									WithIndices(w1, w2, s1, s2)

								for _, peerIdx := range common {
									if s.puzzle.Cell(peerIdx).HasCandidate(y) {
										step.DeleteCandidate(peerIdx, y)
									}
								}

								if len(step.deletedCandidates) > 0 {
									s.applyStep(step)
									return true
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}

func (s *Solver) findEmptyRectangle() (found bool) {
	// Empty Rectangle uses a box where a candidate is restricted to a subset of
	// cells that "pincer" onto a single row and column within that box.
	// If this box-structure is linked to a strong link for the same candidate
	// in a different house, an elimination can be made.

	for x := 1; x <= 9; x++ {
		for b := 0; b < 9; b++ {
			box := s.boxes[b]
			num := box.NumLocations(x)
			if num < 2 {
				continue
			}

			// 1. Identify "pincer" at row R and column C in box B
			// All X's in box B must be in row R or column C.
			for rInBox := 0; rInBox < 3; rInBox++ {
				for cInBox := 0; cInBox < 3; cInBox++ {
					r := (b/3)*3 + rInBox
					c := (b%3)*3 + cInBox

					allIn := true
					locs := box.Locations(x).Values()
					for _, loc := range locs {
						cellR := (b/3)*3 + loc/3
						cellC := (b%3)*3 + loc%3
						if cellR != r && cellC != c {
							allIn = false
							break
						}
					}
					if !allIn {
						continue
					}

					// 2. Look for a strong link on X that connects to this box.
					// Case A: Row R' has a strong link for X at (R', C) and (R', C'').
					// Then we can eliminate X from (R, C'').
					for rPrime := 0; rPrime < 9; rPrime++ {
						if rPrime == r {
							continue
						}
						hRP := s.rows[rPrime]
						if hRP.NumLocations(x) == 2 {
							locsRP := hRP.Locations(x).Values()
							c1, c2 := locsRP[0], locsRP[1]
							var cDoublePrime int
							switch c {
							case c1:
								cDoublePrime = c2
							case c2:
								cDoublePrime = c1
							default:
								continue
							}

							target := r*9 + cDoublePrime
							if s.puzzle.Cell(target).HasCandidate(x) {
								step := NewStep(kindEmptyRectangle).
									WithValues(x).
									WithIndices(r*9+c, rPrime*9+c, rPrime*9+cDoublePrime)
								step.DeleteCandidate(target, x)
								s.applyStep(step)
								return true
							}
						}
					}

					// Case B: Column C' has a strong link for X at (C', R) and (C', R'').
					// Then we can eliminate X from (C', R'').
					for cPrime := 0; cPrime < 9; cPrime++ {
						if cPrime == c {
							continue
						}
						hCP := s.columns[cPrime]
						if hCP.NumLocations(x) == 2 {
							locsCP := hCP.Locations(x).Values()
							r1, r2 := locsCP[0], locsCP[1]
							var rDoublePrime int
							switch r {
							case r1:
								rDoublePrime = r2
							case r2:
								rDoublePrime = r1
							default:
								continue
							}

							target := rDoublePrime*9 + c
							if s.puzzle.Cell(target).HasCandidate(x) {
								step := NewStep(kindEmptyRectangle).
									WithValues(x).
									WithIndices(r*9+c, r*9+cPrime, rDoublePrime*9+cPrime)
								step.DeleteCandidate(target, x)
								s.applyStep(step)
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}

func (s *Solver) findXYWings() (found bool) {
	// Collect a list of all cells with exactly 2 candidates.
	var candidates []*puzzle.Cell
	for i := range 81 {
		if s.puzzle.Cell(i).NumCandidates() == 2 {
			candidates = append(candidates, s.puzzle.Cell(i))
		}
	}
	if len(candidates) < 3 {
		// An XY-Wing requires a pivot cell and 2 pincer cells, so we need at
		// least 3 candidates.
		return false
	}

	// Try each candidate as the pivot cell, checking it against all of the other
	// candidates.
	for _, pivot := range candidates {
		if s.checkXYWingsForPivot(pivot, candidates) {
			return true
		}
	}
	return false
}

func (s *Solver) checkXYWingsForPivot(
	pivot *puzzle.Cell, candidates []*puzzle.Cell,
) (found bool) {
	// Get the x and y values.
	values := pivot.CandidateValues()
	x, y := values[0], values[1]

	// Find the candidate cells that can be seen by the pivot cell and have either
	// x or y as a candidate, but not both.  Collect the cells into separate lists
	// for cells that have x but not y and cells that have y but not x.
	var xCells, yCells []*puzzle.Cell
	for _, cell := range candidates {
		if cell.Index() == pivot.Index() || cell.NumCandidates() != 2 || !seesCell(cell, pivot) {
			continue
		}

		if cell.HasCandidate(x) && !cell.HasCandidate(y) {
			xCells = append(xCells, cell)
		} else if !cell.HasCandidate(x) && cell.HasCandidate(y) {
			yCells = append(yCells, cell)
		}
	}

	if len(xCells) == 0 || len(yCells) == 0 {
		// We need at least one candidate cell for each value to have an XY-Wing.
		return false
	}

	// Check each of the x-cells against each of the y-cells to see if they share
	// a common 3rd value z.
	for _, xc := range xCells {
		// Look for a y-cell that also contains z and is not visible from the x-cell.
		for _, yc := range yCells {
			if seesCell(xc, yc) {
				continue
			}

			// The intersection of {x, z} and {y, z} must be {z}.
			// Since x != y, any intersection is the common value z.
			common := xc.Candidates.Intersection(yc.Candidates)
			if common.Empty() {
				// No common candidate, not an XY-Wing
				continue
			}

			z := common.Value()
			step := NewStep(kindXYWing)
			if s.eliminateXYWingCells(z, xc, yc, step) {
				s.applyStep(step.
					WithIndices(pivot.Index(), xc.Index(), yc.Index()).
					WithValues(x, y, z))
				return true
			}
		}
	}

	return false
}

// eliminateXYWingCells removes candidate value z from all cells that see both
// xCell and yCell.  This assumes that xCell and yCell cannot see each other.
func (s *Solver) eliminateXYWingCells(z int, xCell, yCell *puzzle.Cell, step *SolutionStep) bool {
	seesYCell := func(cell *puzzle.Cell) bool {
		return seesCell(cell, yCell)
	}
	removeZs := func(h *House) bool {
		// Find candidate locations for value z in house h, which is assumed to
		// be a house that contains xCell.
		locs := h.Unsolved[z]
		if !locs.Empty() {
			// Select only the cells that also see yCell.
			cells := h.cellsFromLocs(locs.Values())
			cells = filterSlice(cells, seesYCell)
			for _, zCell := range cells {
				step.DeleteCandidate(zCell.Index(), z)
			}
			// Return true if we found any candidates to remove.
			return len(cells) != 0
		}
		return false
	}

	found := removeZs(s.rows[xCell.Row])
	found = removeZs(s.columns[xCell.Col]) || found
	found = removeZs(s.boxes[xCell.Box()]) || found
	return found
}

func (s *Solver) findAvoidableRectangles() (found bool) {
	// Avoidable Rectangle Type 1:
	// 4 cells forming a rectangle in 2 rows, 2 columns, and 2 boxes.
	// 3 corners are solved (non-givens), 1 corner is unsolved.
	// If the solved corners form an X-Y, Y-X pattern, the unsolved corner
	// cannot take the value that would complete the deadly pattern.

	for r1 := 0; r1 < 9; r1++ {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := 0; c1 < 9; c1++ {
				for c2 := c1 + 1; c2 < 9; c2++ {
					cells := [4]*puzzle.Cell{
						s.puzzle.Get(r1, c1), s.puzzle.Get(r1, c2),
						s.puzzle.Get(r2, c1), s.puzzle.Get(r2, c2),
					}

					// Must be in exactly 2 boxes
					boxes := make(map[int]bool)
					for _, c := range cells {
						boxes[c.Box()] = true
					}
					if len(boxes) != 2 {
						continue
					}

					// Check for Type 1: 3 solved non-givens, 1 unsolved.
					solvedCount := 0
					unsolvedIdx := -1
					for i, c := range cells {
						if c.IsSolved() {
							if c.IsGiven {
								solvedCount = -1 // Cannot use givens
								break
							}
							solvedCount++
						} else {
							unsolvedIdx = i
						}
					}

					switch solvedCount {
					case 3:
						unsolved := cells[unsolvedIdx]
						oppositeIdx := unsolvedIdx ^ 3
						rowAdjIdx := unsolvedIdx ^ 1
						colAdjIdx := unsolvedIdx ^ 2

						X := cells[oppositeIdx].Value()
						Y := cells[rowAdjIdx].Value()

						if cells[colAdjIdx].Value() == Y && X != Y {
							// Potential Avoidable Rectangle Type 1.
							// Unsolved cell cannot be X.
							if unsolved.HasCandidate(X) {
								step := NewStep(kindAvoidableRectangle).
									WithValues(X, Y).
									WithIndices(cells[0].Index(), cells[1].Index(), cells[2].Index(), cells[3].Index())
								step.DeleteCandidate(unsolved.Index(), X)
								s.applyStep(step)
								return true
							}
						}
					case 2:
						// Potential Avoidable Rectangle Type 2.
						var sIdxs, uIdxs []int
						for i, c := range cells {
							if c.IsSolved() {
								sIdxs = append(sIdxs, i)
							} else {
								uIdxs = append(uIdxs, i)
							}
						}

						s1, s2 := cells[sIdxs[0]], cells[sIdxs[1]]
						u1, u2 := cells[uIdxs[0]], cells[uIdxs[1]]

						X, Y := s1.Value(), s2.Value()
						if u1.HasCandidate(X) && u1.HasCandidate(Y) && u2.HasCandidate(X) && u2.HasCandidate(Y) {
							xy := bitset.FromValues16(X, Y)
							e1 := u1.Candidates.Difference(xy)
							e2 := u2.Candidates.Difference(xy)
							commonExtra := e1.Intersection(e2)

							if !commonExtra.Empty() {
								foundElim := false
								step := NewStep(kindAvoidableRectangle).
									WithValues(X, Y).
									WithIndices(cells[0].Index(), cells[1].Index(), cells[2].Index(), cells[3].Index())

								commonPeers := s.commonPeers(u1.Index(), u2.Index())
								for v := range commonExtra.All() {
									for _, pIdx := range commonPeers {
										if s.puzzle.Cell(pIdx).HasCandidate(v) {
											step.DeleteCandidate(pIdx, v)
											foundElim = true
										}
									}
								}

								if foundElim {
									s.applyStep(step)
									return true
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}

// findXYZ searches for 3 cells that fit the "XYZ-Wing" pattern.  An XYZ-Wing
// consists of a pivot cell with 3 candidate values x,y,z, and two pincer cells
// that each see the pivot cell but don't see each other.  One pincer must have
// candidate values x,z and the other must have candidate values y,z.  One of
// these cells must have the value z, so z can be eliminated as a candidate for
// any cell that sees all three.  Note that one pincer *MUST* be in the same
// box as the pivot cell in order for it to be possible for any cell to see the
// pivot and both pincers.
func (s *Solver) findXYZWings() (found bool) {
	// Collect a list of all cells with exactly 3 candidates.
	p := s.puzzle
	var candidates []*puzzle.Cell
	for i := range 81 {
		if p.Cell(i).NumCandidates() == 3 {
			candidates = append(candidates, p.Cell(i))
		}
	}

	// Check each candidate as a possible pivot cell for an XYZ-Wing.
	return slices.ContainsFunc(candidates, s.checkXYZWingsForPivot)
}

func (s *Solver) checkXYZWingsForPivot(pivot *puzzle.Cell) (found bool) {
	// Find cells in the same box as the pivot cell which have exactly 2
	// candidates that both appear in the pivot cell.
	box := s.boxes[pivot.Box()]
	var xzCells []*puzzle.Cell
	for _, cell := range box.Cells {
		if cell.NumCandidates() == 2 {
			// The pivot cell can't match here because it has 3 candidates.
			values := cell.CandidateValues()
			if pivot.HasCandidate(values[0]) && pivot.HasCandidate(values[1]) {
				xzCells = append(xzCells, cell)
			}
		}
	}
	if len(xzCells) == 0 {
		// No valid candidates found.
		return false
	}

	for _, xzCell := range xzCells {
		// Find the y value (pivot - xzCell).
		// Pivot::{x,y,z}, xzCell::{x,z} => Difference::{y}
		ySet := pivot.Candidates.Difference(xzCell.Candidates)
		if ySet.Empty() {
			continue
		} // Should not happen if xzCell is valid subset
		y := ySet.Value()

		// Now find a cell in the same row or column as the pivot cell.
		isYZCandidate := func(cell *puzzle.Cell) bool {
			if cell.Box() == pivot.Box() ||
				cell.NumCandidates() != 2 ||
				!cell.HasCandidate(y) {
				return false
			}

			// Verify it shares Z with xzCell
			// yzCell::{y,z}, xzCell::{x,z} => Intersection::{z}
			// (Assuming x != y, which is true)
			return cell.Candidates.Intersects(xzCell.Candidates)
		}

		yzCells := slices.Concat(
			s.rows[pivot.Row].Cells[:],
			s.columns[pivot.Col].Cells[:],
		)
		for _, yzCell := range yzCells {
			step := NewStep(kindXYZWing)
			if isYZCandidate(yzCell) &&
				s.eliminateXYZWingCells(pivot, xzCell, yzCell, step) {
				s.applyStep(step.
					WithIndices(pivot.Index(), xzCell.Index(), yzCell.Index()).
					WithValues(pivot.CandidateValues()...))
				return true
			}
		}
	}

	return false
}

// eliminateXYZWingCells removes candidate value z from any cells that see all
// three of xyzCell, xzCell, and yzCell.  The value x is the one candidate value
// that appears as a candidate in all 3 cells.  This assumes that xzCell and
// yzCell cannot see each other, and that xzCell is in the same box as xyzCell.
func (s *Solver) eliminateXYZWingCells(xyzCell, xzCell, yzCell *puzzle.Cell, step *SolutionStep) bool {
	// The z value is the only common candidate between xzCell and yzCell.
	zSet := xzCell.Candidates.Intersection(yzCell.Candidates)
	if zSet.Empty() {
		return false
	}
	z := zSet.Value()

	// The only cells that could possibly see all three XYZ-Wing cells are the
	// other cells in the same box as xyzCell and xzCell, so we just need to
	// check the candidate locations for value z in that box and select the
	// ones that can see yzCell.
	box := s.boxes[xyzCell.Box()]
	locs := box.Unsolved[z]
	cells := box.cellsFromLocs(locs.Values())
	cells = filterSlice(cells, func(cell *puzzle.Cell) bool {
		return cell.Index() != xyzCell.Index() &&
			cell.Index() != xzCell.Index() &&
			seesCell(cell, yzCell)
	})
	if len(cells) == 0 {
		// No candidates found to eliminate.
		return false
	}

	for _, xCell := range cells {
		step.DeleteCandidate(xCell.Index(), z)
	}
	return true
}

func (s *Solver) findHiddenQuadruples() (found bool) {
	return s.findHiddenSubsets(4, kindHiddenQuadruple)
}

// FISH TECHNIQUES

func (s *Solver) findXWings() (found bool) {
	return s.findFishOfSize(2, kindXWing)
}

func (s *Solver) findSwordfish() (found bool) {
	return s.findFishOfSize(3, kindSwordfish)
}

func (s *Solver) findJellyfish() (found bool) {
	return s.findFishOfSize(4, kindJellyfish)
}

func (s *Solver) findFishOfSize(fishSize int, fishKind techniqueKind) (found bool) {
	find := func(baseLines, coverLines []*House) bool {
		return s.findFishInLines(fishSize, fishKind, baseLines, coverLines)
	}
	return find(s.rows, s.columns) || find(s.columns, s.rows)
}

func (s *Solver) findFishInLines(
	fishSize int,
	fishKind techniqueKind,
	baseLines, coverLines []*House,
) (found bool) {
	for _, base := range baseLines {
		for val := 1; val <= 9; val++ {
			locs := base.Unsolved[val]
			if locs.Empty() {
				continue
			}

			// A fish line must have no more than fishSize candidate locations
			// for a value. We assume that all singles and smaller fish have
			// already been found.
			if locs.Size() > fishSize {
				continue
			}

			if s.checkFishForValue(fishSize, fishKind, val, base, baseLines, coverLines) {
				return true
			}
		}
	}

	return false
}

func (s *Solver) checkFishForValue(
	fishSize int,
	fishKind techniqueKind,
	val int,
	base1 *House,
	baseLines, coverLines []*House,
) (found bool) {
	// Find all base lines other than base1 that have either 2 or 3 candidate
	// locations for val.
	candidates := filterSlice(baseLines, func(b2 *House) bool {
		numLocs := b2.NumLocations(val)
		return b2.Index != base1.Index && numLocs >= 2 && numLocs <= fishSize
	})

	valueSet := bitset.FromValues16(val)
	// Variables for storing search results.
	fishLines := []int{base1.Index}
	var coverLocs LocSet

	// Must forward-declare func so we can call it self-recursively.
	var checkLines func(lines []*House, fishLocs LocSet) bool
	checkLines = func(lines []*House, fishLocs LocSet) bool {
		for i, line := range lines {
			locs := bitset.Union(fishLocs, line.Unsolved[val])
			if locs.Size() > fishSize {
				// Too many locations, so this line can't be part of the fish.
				continue
			}
			fishLines = append(fishLines, line.Index)
			if len(fishLines) == fishSize {
				// We found enough lines, so we have a fish.
				// fishLines contains the base lines that make up the fish,
				// but we need to save the set containing the cover lines.
				coverLocs = locs
				return true
			}
			// Check recursively until we have enough lines.
			if checkLines(lines[i+1:], locs) {
				return true
			}

			// No fish found, so backtrack the last line and keep trying.
			fishLines = fishLines[:len(fishLines)-1]
		}
		// No more candidate lines, so we don't have a fish.
		return false
	}

	step := NewStep(fishKind).WithValues(val)
	if checkLines(candidates, base1.Unsolved[val]) {
		// We found a fish.
		bases := transformSlice(fishLines, func(x int) *House {
			return baseLines[x]
		})
		covers := transformSlice(coverLocs.Values(), func(y int) *House {
			return coverLines[y]
		})
		locSet := bitset.FromValues16(fishLines...)
		if s.eliminateFromOtherLocsMulti(covers, valueSet, locSet, step) {
			s.applyStep(step.
				WithBases(bases...).
				WithCovers(covers...))
			return true
		}
	}

	return false
}

func (s *Solver) findHiddenRectangle() (found bool) {
	// Find all 4-cell rectangles in 2 rows, 2 columns, and 2 boxes.
	for r1 := 0; r1 < 9; r1++ {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := 0; c1 < 9; c1++ {
				for c2 := c1 + 1; c2 < 9; c2++ {
					cells := [4]*puzzle.Cell{
						s.puzzle.Get(r1, c1), s.puzzle.Get(r1, c2),
						s.puzzle.Get(r2, c1), s.puzzle.Get(r2, c2),
					}

					// All must be unsolved
					allUnsolved := true
					for _, c := range cells {
						if c.IsSolved() {
							allUnsolved = false
							break
						}
					}
					if !allUnsolved {
						continue
					}

					// Must be in exactly 2 boxes
					boxes := make(map[int]bool)
					for _, c := range cells {
						boxes[c.Box()] = true
					}
					if len(boxes) != 2 {
						continue
					}

					// Find candidate pairs {x, y} common to all 4 cells.
					common := cells[0].Candidates.
						Intersection(cells[1].Candidates).
						Intersection(cells[2].Candidates).
						Intersection(cells[3].Candidates)

					if common.Size() < 2 {
						continue
					}

					vals := common.Values()
					for i := 0; i < len(vals); i++ {
						for j := i + 1; j < len(vals); j++ {
							x, y := vals[i], vals[j]

							// Check each corner as the potential target corner C1.
							for idx1 := 0; idx1 < 4; idx1++ {
								idx4 := idx1 ^ 3 // diagonal (C4)
								idx2 := idx1 ^ 1 // same row (C2)
								idx3 := idx1 ^ 2 // same column (C3)

								c1, c2, c3, c4 := cells[idx1], cells[idx2], cells[idx3], cells[idx4]

								// C4 must be bivalue {x, y}
								xy := bitset.FromValues16(x, y)
								if !c4.Candidates.Equal(xy) {
									continue
								}

								// Check for Hidden Rectangle on candidate x (eliminate y from C1)
								if s.checkHiddenRectangleForValue(c1, c2, c3, c4, x, y) {
									return true
								}
								// Check for Hidden Rectangle on candidate y (eliminate x from C1)
								if s.checkHiddenRectangleForValue(c1, c2, c3, c4, y, x) {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (s *Solver) checkHiddenRectangleForValue(c1, c2, c3, c4 *puzzle.Cell, val, otherVal int) bool {
	// val must be a conjugate pair in Row(c1) restricted to c1, c2.
	row := s.rows[c1.Row]
	if row.NumLocations(val) != 2 {
		return false
	}
	if !row.Locations(val).Contains(c1.Col) || !row.Locations(val).Contains(c2.Col) {
		return false
	}

	// val must be a conjugate pair in Col(c1) restricted to c1, c3.
	col := s.columns[c1.Col]
	if col.NumLocations(val) != 2 {
		return false
	}
	if !col.Locations(val).Contains(c1.Row) || !col.Locations(val).Contains(c3.Row) {
		return false
	}

	// Elimination: otherVal from c1.
	if c1.HasCandidate(otherVal) {
		step := NewStep(kindHiddenRectangle).
			WithValues(val, otherVal).
			WithIndices(c1.Index(), c2.Index(), c3.Index(), c4.Index())
		step.DeleteCandidate(c1.Index(), otherVal)
		s.applyStep(step)
		return true
	}
	return false
}

func (s *Solver) findFinnedXWings() (found bool) {
	check := func(baseLines, coverLines []*House) bool {
		return s.checkFinnedXWings(baseLines, coverLines)
	}
	// Rows as base, columns as cover.
	if check(s.rows, s.columns) {
		return true
	}
	// Columns as base, rows as cover.
	return check(s.columns, s.rows)
}

func (s *Solver) checkFinnedXWings(baseLines, coverLines []*House) bool {
	for val := 1; val <= 9; val++ {
		// Find candidate base lines with at least 2 candidates.
		var candidates []*House
		for _, h := range baseLines {
			if h.NumLocations(val) >= 2 {
				candidates = append(candidates, h)
			}
		}

		if len(candidates) < 2 {
			continue
		}

		// Try each pair of lines as base lines (b1, b2).
		// Note: we try (i, j) and (j, i) because the fin can be on either one.
		for i := 0; i < len(candidates); i++ {
			for j := 0; j < len(candidates); j++ {
				if i == j {
					continue
				}
				b1, b2 := candidates[i], candidates[j]

				// b1 must have exactly 2 candidates at L1, L2.
				if b1.NumLocations(val) != 2 {
					continue
				}
				locs1 := b1.Locations(val).Values()
				l1, l2 := locs1[0], locs1[1]

				// b2 has candidates. At least one must be at L2.
				if !b2.Locations(val).Contains(l2) {
					continue
				}

				// Fins are all candidates in b2 that are NOT at l2.
				fins := b2.Locations(val).Difference(bitset.FromValues16(l2))
				if fins.Empty() {
					continue
				}

				// If it's a perfect X-Wing, b2.Locations(val) would be {l1, l2}.
				// A Finned X-Wing has additional candidates in b2, but they
				// must all be in the same box as (b2, l1).
				corner21 := b2.Cells[l1]
				boxIdx := corner21.Box()

				allFinsInBox := true
				for _, finLoc := range fins.Values() {
					if b2.Cells[finLoc].Box() != boxIdx {
						allFinsInBox = false
						break
					}
				}

				if allFinsInBox {
					// Found a Finned/Sashimi X-Wing!
					// Elimination: Cells in Box(boxIdx) that are on cover line l1
					// (excluding corners themselves).
					step := NewStep(kindFinnedXWing).
						WithValues(val).
						WithBases(b1, b2).
						WithCovers(coverLines[l1], coverLines[l2])

					foundElim := false
					box := s.boxes[boxIdx]
					for _, cell := range box.Cells {
						// Identify cover line index.
						var cellLineIdx int
						if b1.Kind == kindRow {
							cellLineIdx = cell.Col
						} else {
							cellLineIdx = cell.Row
						}

						if cellLineIdx != l1 {
							continue
						}

						// Cannot be the corner from base b1.
						if cell.Index() == b1.Cells[l1].Index() {
							continue
						}
						// Cannot be any cell in b2 (the fish line itself).
						if (b1.Kind == kindRow && cell.Row == b2.Index) ||
							(b1.Kind == kindColumn && cell.Col == b2.Index) {
							continue
						}

						if cell.HasCandidate(val) {
							step.DeleteCandidate(cell.Index(), val)
							foundElim = true
						}
					}

					if foundElim {
						// Collect indices for highlighting.
						indices := []int{b1.Cells[l1].Index(), b1.Cells[l2].Index(), b2.Cells[l2].Index()}
						for _, fLoc := range fins.Values() {
							indices = append(indices, b2.Cells[fLoc].Index())
						}
						s.applyStep(step.WithIndices(indices...))
						return true
					}
				}
			}
		}
	}
	return false
}

// SINGLE-DIGIT TECHNIQUES

func (s *Solver) findSkyscraper() (found bool) {
	check := func(baseLines []*House) bool {
		return s.checkSkyscraper(baseLines)
	}
	// We check for Skyscrapers where the base lines are rows (meaning the strong
	// links are horizontal) and then where the base lines are columns (vertical
	// strong links).
	return check(s.rows) || check(s.columns)
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
				if s.eliminatePeerTargets(val, top1, top2, step) {
					s.applyStep(step.WithIndices(
						top1.Index(), top2.Index(), floor1.Index(), floor2.Index(),
					))
					return true
				}
			}
		}
	}

	return false
}

func (s *Solver) eliminatePeerTargets(
	val int, c1, c2 *puzzle.Cell, step *SolutionStep,
) bool {
	found := false

	// Check all cells that see c1 to see if they also see c2.
	// We check the 3 houses that contain c1.
	houses := []*House{
		s.rows[c1.Row],
		s.columns[c1.Col],
		s.boxes[c1.Box()],
	}

	checked := make(map[int]bool)

	for _, h := range houses {
		for _, cell := range h.Cells {
			// Skip the cells themselves.
			if cell.Index() == c1.Index() || cell.Index() == c2.Index() {
				continue
			}

			// Avoid checking the same cell twice (e.g. if it's in multiple shared houses)
			idx := cell.Row*9 + cell.Col
			if checked[idx] {
				continue
			}
			checked[idx] = true

			// If the cell contains the candidate value and sees c2, then it
			// sees both c1 and c2, so we can eliminate the candidate.
			if cell.HasCandidate(val) && seesCell(cell, c2) {
				step.DeleteCandidate(cell.Index(), val)
				found = true
			}
		}
	}

	return found
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

			if s.eliminatePeerTargets(val, tail1, tail2, step) {
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

func (s *Solver) findSimpleColoring() bool {
	globalStep := NewStep(kindSimpleColoring)
	found := false

	for val := 1; val <= 9; val++ {
		adj := make(map[int][]int)
		addLink := func(u, v int) {
			for _, existing := range adj[u] {
				if existing == v {
					return
				}
			}
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		activeLits := make(map[int]bool)
		for _, h := range s.houses {
			// Rule 1: exactly 2 cells
			locs := h.Unsolved[val].Values()
			if len(locs) == 2 {
				u := h.Cells[locs[0]].Index()
				v := h.Cells[locs[1]].Index()
				addLink(u, v)
				activeLits[u] = true
				activeLits[v] = true
			}

			// Grouped strong links
			switch h.Kind {
			case kindRow:
				activeBoxes := []int{}
				for bIdxInRow := 0; bIdxInRow < 3; bIdxInRow++ {
					boxIdx := (h.Index/3)*3 + bIdxInRow
					hasCandidate := false
					for _, cell := range s.boxes[boxIdx].Cells {
						if cell.Row == h.Index && cell.HasCandidate(val) {
							hasCandidate = true
							break
						}
					}
					if hasCandidate {
						activeBoxes = append(activeBoxes, bIdxInRow)
					}
				}
				if len(activeBoxes) == 2 {
					u := 81 + h.Index*3 + activeBoxes[0]
					v := 81 + h.Index*3 + activeBoxes[1]
					addLink(u, v)
					activeLits[u] = true
					activeLits[v] = true
				}
			case kindColumn:
				activeBoxes := []int{}
				for bIdxInCol := 0; bIdxInCol < 3; bIdxInCol++ {
					boxIdx := (h.Index/3) + bIdxInCol*3
					hasCandidate := false
					for _, cell := range s.boxes[boxIdx].Cells {
						if cell.Col == h.Index && cell.HasCandidate(val) {
							hasCandidate = true
							break
						}
					}
					if hasCandidate {
						activeBoxes = append(activeBoxes, bIdxInCol)
					}
				}
				if len(activeBoxes) == 2 {
					u := 108 + h.Index*3 + activeBoxes[0]
					v := 108 + h.Index*3 + activeBoxes[1]
					addLink(u, v)
					activeLits[u] = true
					activeLits[v] = true
				}
			case kindBox:
				// Rows in box
				activeRows := []int{}
				for rInBox := 0; rInBox < 3; rInBox++ {
					r := (h.Index/3)*3 + rInBox
					hasCandidate := false
					for _, cell := range h.Cells {
						if cell.Row == r && cell.HasCandidate(val) {
							hasCandidate = true
							break
						}
					}
					if hasCandidate {
						activeRows = append(activeRows, r)
					}
				}
				if len(activeRows) == 2 {
					u := 81 + activeRows[0]*3 + (h.Index % 3)
					v := 81 + activeRows[1]*3 + (h.Index % 3)
					addLink(u, v)
					activeLits[u] = true
					activeLits[v] = true
				}
				// Cols in box
				activeCols := []int{}
				for cInBox := 0; cInBox < 3; cInBox++ {
					c := (h.Index%3)*3 + cInBox
					hasCandidate := false
					for _, cell := range h.Cells {
						if cell.Col == c && cell.HasCandidate(val) {
							hasCandidate = true
							break
						}
					}
					if hasCandidate {
						activeCols = append(activeCols, c)
					}
				}
				if len(activeCols) == 2 {
					u := 108 + activeCols[0]*3 + (h.Index / 3)
					v := 108 + activeCols[1]*3 + (h.Index / 3)
					addLink(u, v)
					activeLits[u] = true
					activeLits[v] = true
				}
			}
		}

		visited := make(map[int]bool)
		for startLit := range activeLits {
			if visited[startLit] {
				continue
			}
			component := make(map[int]int)
			queue := []int{startLit}
			component[startLit] = 0
			visited[startLit] = true
			head := 0
			for head < len(queue) {
				u := queue[head]
				head++
				for _, v := range adj[u] {
					if _, ok := component[v]; ok {
						continue
					}
					component[v] = 1 - component[u]
					visited[v] = true
					queue = append(queue, v)
				}
			}
			colorUnits := [2][]int{}
			for lit, color := range component {
				colorUnits[color] = append(colorUnits[color], lit)
			}

			for color := 0; color < 2; color++ {
				units := colorUnits[color]
				contradiction := false
				for i := 0; i < len(units); i++ {
					for j := i + 1; j < len(units); j++ {
						if s.unitsShareHouse(val, units[i], units[j]) {
							contradiction = true
							break
						}
					}
					if contradiction {
						break
					}
				}
				if contradiction {
					found = true
					for _, lit := range colorUnits[color] {
						for _, cIdx := range s.getLitCells(val, lit) {
							globalStep.DeleteCandidate(cIdx, val)
						}
					}
					for _, lit := range colorUnits[1-color] {
						cells := s.getLitCells(val, lit)
						if len(cells) == 1 {
							globalStep.PlaceCandidate(cells[0], val)
						} else {
							s.eliminateFromGroupHouses(val, cells, globalStep)
						}
					}
				}
			}

			for cIdx := 0; cIdx < 81; cIdx++ {
				if !s.puzzle.Cell(cIdx).HasCandidate(val) {
					continue
				}
				inComponent := false
				for lit := range component {
					for _, uCell := range s.getLitCells(val, lit) {
						if uCell == cIdx {
							inComponent = true
							break
						}
					}
					if inComponent {
						break
					}
				}
				if !inComponent {
					sees0 := false
					for _, lit := range colorUnits[0] {
						if s.cellSeesAll(cIdx, s.getLitCells(val, lit)) {
							sees0 = true
							break
						}
					}
					if sees0 {
						sees1 := false
						for _, lit := range colorUnits[1] {
							if s.cellSeesAll(cIdx, s.getLitCells(val, lit)) {
								sees1 = true
								break
							}
						}
						if sees1 {
							globalStep.DeleteCandidate(cIdx, val)
							found = true
						}
					}
				}
			}
		}
	}
	if found {
		s.applyStep(globalStep)
	}
	return found
}

func (s *Solver) getLitCells(v, lit int) []int {
	if lit < 81 {
		return []int{lit}
	}
	cells := []int{}
	if lit < 108 { // rbLit
		litIdx := lit - 81
		r := litIdx / 3
		bIdxInRow := litIdx % 3
		boxIdx := (r/3)*3 + bIdxInRow
		for _, cell := range s.boxes[boxIdx].Cells {
			if cell.Row == r && cell.HasCandidate(v) {
				cells = append(cells, cell.Index())
			}
		}
	} else { // cbLit
		litIdx := lit - 108
		c := litIdx / 3
		bIdxInCol := litIdx % 3
		boxIdx := (c/3) + bIdxInCol*3
		for _, cell := range s.boxes[boxIdx].Cells {
			if cell.Col == c && cell.HasCandidate(v) {
				cells = append(cells, cell.Index())
			}
		}
	}
	return cells
}

func (s *Solver) unitsShareHouse(v, lit1, lit2 int) bool {
	cells1 := s.getLitCells(v, lit1)
	cells2 := s.getLitCells(v, lit2)
	houses1 := s.getHousesContainingAll(cells1)
	houses2 := s.getHousesContainingAll(cells2)
	for _, h1 := range houses1 {
		for _, h2 := range houses2 {
			if h1 == h2 {
				return true
			}
		}
	}
	return false
}

func (s *Solver) getHousesContainingAll(cells []int) []int {
	if len(cells) == 0 {
		return nil
	}
	res := []int{}
	// Row
	r := cells[0] / 9
	allInRow := true
	for _, c := range cells {
		if c/9 != r {
			allInRow = false
			break
		}
	}
	if allInRow {
		res = append(res, r)
	}
	// Col
	col := cells[0] % 9
	allInCol := true
	for _, c := range cells {
		if c%9 != col {
			allInCol = false
			break
		}
	}
	if allInCol {
		res = append(res, 9+col)
	}
	// Box
	b := (cells[0] / 9 / 3) * 3 + (cells[0] % 9 / 3)
	allInBox := true
	for _, c := range cells {
		if (c/9/3)*3+(c%9/3) != b {
			allInBox = false
			break
		}
	}
	if allInBox {
		res = append(res, 18+b)
	}
	return res
}

func (s *Solver) cellSeesAll(c int, cells []int) bool {
	if len(cells) == 0 {
		return false
	}
	for _, other := range cells {
		if c == other || !s.sees(c, other) {
			return false
		}
	}
	return true
}

func (s *Solver) eliminateFromGroupHouses(v int, cells []int, step *SolutionStep) {
	houses := s.getHousesContainingAll(cells)
	for _, hIdx := range houses {
		var h *House
		switch {
		case hIdx < 9:
			h = s.rows[hIdx]
		case hIdx < 18:
			h = s.columns[hIdx-9]
		default:
			h = s.boxes[hIdx-18]
		}

		cellSet := make(map[int]bool)
		for _, c := range cells {
			cellSet[c] = true
		}

		for _, cell := range h.Cells {
			if !cellSet[cell.Index()] && cell.HasCandidate(v) {
				step.DeleteCandidate(cell.Index(), v)
			}
		}
	}
}

// UNIQUENESS TECHNIQUES

func (s *Solver) findUniqueRectangleType1() (found bool) {
	b := s.puzzle
	// Check each cell with exactly 2 candidate values to see if it is the base
	// corner of a unique rectangle.
	for r := range 9 {
		for c := range 9 {
			cell := b.Get(r, c)
			if cell.NumCandidates() != 2 {
				continue
			}
			if s.checkUniqueRectangleForCell(cell) {
				return true
			}
		}
	}

	return false
}

func (s *Solver) checkUniqueRectangleForCell(base *puzzle.Cell) (found bool) {
	b := s.puzzle

	// Look for a cell in the same row as base with the same pair of candidates.
	var rowWing *puzzle.Cell
	for c := range 9 {
		if c != base.Col {
			cell := b.Get(base.Row, c)
			if sameCandidates(base, cell) {
				rowWing = cell
				break
			}
		}
	}
	if rowWing == nil {
		return false
	}

	// Look for a cell in the same column as base with the same pair of candidates.
	var colWing *puzzle.Cell
	for r := range 9 {
		if r != base.Row {
			cell := b.Get(r, base.Col)
			if sameCandidates(base, cell) {
				colWing = cell
				break
			}
		}
	}
	if colWing == nil {
		return false
	}

	// The 2 wing cells must be in different boxes, but one of them must be in
	// the same box as the base.
	if rowWing.Box() != colWing.Box() &&
		(rowWing.Box() == base.Box() || colWing.Box() == base.Box()) {
		// These cells form a unique rectangle, so we can eliminate their candidates
		// from the cell at the 4th corner of the rectangle, which will have the
		// same row as the column-wing and the same column as the row-wing.
		step := NewStep(kindUniqueRectangle1)
		if s.eliminateValuesFromCell(colWing.Row, rowWing.Col, base.Candidates, step) {
			s.applyStep(step.
				WithValues(base.Candidates.Values()...).
				WithIndices(base.Index(), rowWing.Index(), colWing.Index()))
			return true
		}
	}

	return false
}

func (s *Solver) findUniqueRectangleType2() (found bool) {
	// Find all 4-cell rectangles in 2 rows, 2 columns, and 2 boxes.
	for r1 := 0; r1 < 9; r1++ {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := 0; c1 < 9; c1++ {
				for c2 := c1 + 1; c2 < 9; c2++ {
					cells := [4]*puzzle.Cell{
						s.puzzle.Get(r1, c1), s.puzzle.Get(r1, c2),
						s.puzzle.Get(r2, c1), s.puzzle.Get(r2, c2),
					}

					// All must be unsolved
					allUnsolved := true
					for _, c := range cells {
						if c.IsSolved() {
							allUnsolved = false
							break
						}
					}
					if !allUnsolved {
						continue
					}

					// Must be in exactly 2 boxes
					boxes := make(map[int]bool)
					for _, c := range cells {
						boxes[c.Box()] = true
					}
					if len(boxes) != 2 {
						continue
					}

					// Find candidate pairs {x, y} common to all 4 cells.
					common := cells[0].Candidates.
						Intersection(cells[1].Candidates).
						Intersection(cells[2].Candidates).
						Intersection(cells[3].Candidates)

					if common.Size() < 2 {
						continue
					}

					vals := common.Values()
					for i := 0; i < len(vals); i++ {
						for j := i + 1; j < len(vals); j++ {
							x, y := vals[i], vals[j]
							xy := bitset.FromValues16(x, y)

							var bivalueCells, extraCells []*puzzle.Cell
							for _, c := range cells {
								if c.Candidates.Equal(xy) {
									bivalueCells = append(bivalueCells, c)
								} else {
									extraCells = append(extraCells, c)
								}
							}

							if len(bivalueCells) == 2 && len(extraCells) == 2 {
								e1 := extraCells[0].Candidates.Difference(xy)
								e2 := extraCells[1].Candidates.Difference(xy)

								if e1.Size() == 1 && e1.Equal(e2) {
									z := e1.Values()[0]

									// Z can be eliminated from any cell that sees BOTH extraCells.
									commonPeers := s.commonPeers(extraCells[0].Index(), extraCells[1].Index())
									step := NewStep(kindUniqueRectangle2).
										WithValues(x, y, z).
										WithIndices(cells[0].Index(), cells[1].Index(), cells[2].Index(), cells[3].Index())

									for _, pIdx := range commonPeers {
										if s.puzzle.Cell(pIdx).HasCandidate(z) {
											step.DeleteCandidate(pIdx, z)
										}
									}

									if len(step.deletedCandidates) > 0 {
										s.applyStep(step)
										return true
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (s *Solver) findUniqueRectangleType3() (found bool) {
	// Find all 4-cell rectangles in 2 rows, 2 columns, and 2 boxes.
	for r1 := 0; r1 < 9; r1++ {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := 0; c1 < 9; c1++ {
				for c2 := c1 + 1; c2 < 9; c2++ {
					cells := [4]*puzzle.Cell{
						s.puzzle.Get(r1, c1), s.puzzle.Get(r1, c2),
						s.puzzle.Get(r2, c1), s.puzzle.Get(r2, c2),
					}

					// All must be unsolved
					allUnsolved := true
					for _, c := range cells {
						if c.IsSolved() {
							allUnsolved = false
							break
						}
					}
					if !allUnsolved {
						continue
					}

					// Must be in exactly 2 boxes
					boxes := make(map[int]bool)
					for _, c := range cells {
						boxes[c.Box()] = true
					}
					if len(boxes) != 2 {
						continue
					}

					// Find candidate pairs {x, y} common to all 4 cells.
					common := cells[0].Candidates.
						Intersection(cells[1].Candidates).
						Intersection(cells[2].Candidates).
						Intersection(cells[3].Candidates)

					if common.Size() < 2 {
						continue
					}

					vals := common.Values()
					for i := 0; i < len(vals); i++ {
						for j := i + 1; j < len(vals); j++ {
							x, y := vals[i], vals[j]
							xy := bitset.FromValues16(x, y)

							var bivalueCells, extraCells []*puzzle.Cell
							for _, c := range cells {
								if c.Candidates.Equal(xy) {
									bivalueCells = append(bivalueCells, c)
								} else {
									extraCells = append(extraCells, c)
								}
							}

							if len(bivalueCells) == 2 && len(extraCells) == 2 {
								var h *House
								if extraCells[0].Row == extraCells[1].Row {
									h = s.rows[extraCells[0].Row]
								} else if extraCells[0].Col == extraCells[1].Col {
									h = s.columns[extraCells[0].Col]
								}

								if h == nil {
									continue
								}

								extraCandidates := bitset.Union(extraCells[0].Candidates, extraCells[1].Candidates).
									Difference(xy)

								if extraCandidates.Empty() {
									continue
								}

								if s.checkURType3NakedSubset(h, extraCells, cells[:], extraCandidates, xy) {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (s *Solver) checkURType3NakedSubset(
	h *House,
	urExtraCells []*puzzle.Cell,
	urCells []*puzzle.Cell,
	extraCandidates ValSet,
	xy ValSet,
) bool {
	K := extraCandidates.Size()
	if K > 4 {
		return false
	}

	// Try subset sizes N from K up to 4.
	for N := K; N <= 4; N++ {
		var potentialCells []*puzzle.Cell
		for _, cell := range h.Cells {
			if cell.IsSolved() {
				continue
			}
			if cell.Index() == urExtraCells[0].Index() || cell.Index() == urExtraCells[1].Index() {
				continue
			}
			if cell.NumCandidates() <= N && cell.Candidates.Intersects(extraCandidates) {
				potentialCells = append(potentialCells, cell)
			}
		}

		if len(potentialCells) < N-1 {
			continue
		}

		if s.findURType3SubsetCombination(h, urExtraCells, urCells, extraCandidates, potentialCells, N, xy) {
			return true
		}
	}
	return false
}

func (s *Solver) findURType3SubsetCombination(
	h *House,
	urExtraCells []*puzzle.Cell,
	urCells []*puzzle.Cell,
	extraCandidates ValSet,
	potentialCells []*puzzle.Cell,
	N int,
	xy ValSet,
) bool {
	var check func(start int, current []*puzzle.Cell) bool
	check = func(start int, current []*puzzle.Cell) bool {
		if len(current) == N-1 {
			union := extraCandidates
			for _, c := range current {
				union = bitset.Union(union, c.Candidates)
			}

			if union.Size() == N {
				subsetLocs := bitset.BitSet16(0)
				subsetLocs.Add(h.locFromIndex(urExtraCells[0].Index()))
				subsetLocs.Add(h.locFromIndex(urExtraCells[1].Index()))
				for _, c := range current {
					subsetLocs.Add(h.locFromIndex(c.Index()))
				}

				step := NewStep(kindUniqueRectangle3).
					WithValues(xy.Values()...).
					WithHouse(h)

				if s.eliminateFromOtherLocs(h, union, subsetLocs, step) {
					indices := make([]int, 0, 4+len(current))
					for _, c := range urCells {
						indices = append(indices, c.Index())
					}
					for _, c := range current {
						indices = append(indices, c.Index())
					}
					s.applyStep(step.WithIndices(indices...))
					return true
				}
			}
			return false
		}

		for i := start; i < len(potentialCells); i++ {
			if check(i+1, append(current, potentialCells[i])) {
				return true
			}
		}
		return false
	}

	return check(0, nil)
}

func (s *Solver) findUniqueRectangleType4() (found bool) {
	// Find all 4-cell rectangles in 2 rows, 2 columns, and 2 boxes.
	for r1 := 0; r1 < 9; r1++ {
		for r2 := r1 + 1; r2 < 9; r2++ {
			for c1 := 0; c1 < 9; c1++ {
				for c2 := c1 + 1; c2 < 9; c2++ {
					cells := [4]*puzzle.Cell{
						s.puzzle.Get(r1, c1), s.puzzle.Get(r1, c2),
						s.puzzle.Get(r2, c1), s.puzzle.Get(r2, c2),
					}

					// All must be unsolved
					allUnsolved := true
					for _, c := range cells {
						if c.IsSolved() {
							allUnsolved = false
							break
						}
					}
					if !allUnsolved {
						continue
					}

					// Must be in exactly 2 boxes
					boxes := make(map[int]bool)
					for _, c := range cells {
						boxes[c.Box()] = true
					}
					if len(boxes) != 2 {
						continue
					}

					// Find candidate pairs {x, y} common to all 4 cells.
					common := cells[0].Candidates.
						Intersection(cells[1].Candidates).
						Intersection(cells[2].Candidates).
						Intersection(cells[3].Candidates)

					if common.Size() < 2 {
						continue
					}

					vals := common.Values()
					for i := 0; i < len(vals); i++ {
						for j := i + 1; j < len(vals); j++ {
							x, y := vals[i], vals[j]
							xy := bitset.FromValues16(x, y)

							var bivalueCells, extraCells []*puzzle.Cell
							for _, c := range cells {
								if c.Candidates.Equal(xy) {
									bivalueCells = append(bivalueCells, c)
								} else {
									extraCells = append(extraCells, c)
								}
							}

							if len(bivalueCells) == 2 && len(extraCells) == 2 {
								// Type 4: One of the candidates (x or y) must be a conjugate pair
								// in a house shared by the extraCells.
								for _, val := range []int{x, y} {
									otherVal := x
									if val == x {
										otherVal = y
									}

									// Find shared house of extraCells.
									var houses []*House
									if extraCells[0].Row == extraCells[1].Row {
										houses = append(houses, s.rows[extraCells[0].Row])
									}
									if extraCells[0].Col == extraCells[1].Col {
										houses = append(houses, s.columns[extraCells[0].Col])
									}

									for _, h := range houses {
										// Is val a conjugate pair in this house, restricted to the two extraCells?
										if h.NumLocations(val) == 2 {
											locs := h.Locations(val).Values()
											c1Idx := h.Cells[locs[0]].Index()
											c2Idx := h.Cells[locs[1]].Index()

											if (c1Idx == extraCells[0].Index() && c2Idx == extraCells[1].Index()) ||
												(c1Idx == extraCells[1].Index() && c2Idx == extraCells[0].Index()) {
												// Elimination: otherVal can be removed from extraCells.
												step := NewStep(kindUniqueRectangle4).
													WithValues(x, y).
													WithHouse(h)

												foundElim := false
												for _, c := range extraCells {
													if c.HasCandidate(otherVal) {
														step.DeleteCandidate(c.Index(), otherVal)
														foundElim = true
													}
												}

												if foundElim {
													s.applyStep(step.WithIndices(
														cells[0].Index(), cells[1].Index(),
														cells[2].Index(), cells[3].Index(),
													))
													return true
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func (s *Solver) eliminateValuesFromCell(
	r, c int, values ValSet, step *SolutionStep,
) bool {
	cell := s.puzzle.Get(r, c)
	found := false
	for _, v := range values.Values() {
		if cell.HasCandidate(v) {
			step.DeleteCandidate(r*9+c, v)
			found = true
		}
	}
	return found
}

// BRUTE FORCE (LAST RESORT) TECHNIQUE

// findBruteForce uses a brute force search to find a solution for any remaining
// unsolved cells.  The search uses Donald Knuth's Algorithm X with the "Dancing
// Links" technique.
func (s *Solver) findBruteForce() bool {
	dl := NewDancingLinks(s.puzzle)
	dlOptions := &DancingLinksOptions{
		EnableDebug: s.EnableDebug,
		TimeLimit:   5 * time.Second,
	}

	solved, _ := dl.SolveWithStats(dlOptions)
	if solved {
		s.applyBruteForceSteps(dl)
	}

	return solved
}

func (s *Solver) applyBruteForceSteps(dl *DancingLinks) {
	steps := dl.GetSolution()
	for _, step := range steps {
		s.appendNextStep(NewStep(kindBruteForce).WithPlacedValue(step.Index, step.Value))
		// Place the value in the puzzle.
		s.puzzle.PlaceValue(step.Index, step.Value)
	}
}

func (s *Solver) findSueDeCoq() bool {
	deleted := make([]bitset.BitSet16, 81)

	// 54 intersections: 27 box-row, 27 box-column
	for bIdx := 0; bIdx < 9; bIdx++ {
		box := s.boxes[bIdx]
		// Box bIdx intersects 3 rows:
		for rIdx := (bIdx / 3) * 3; rIdx < (bIdx/3)*3+3; rIdx++ {
			step := NewStep(kindSueDeCoq)
			if s.checkSueDeCoq(box, s.rows[rIdx], step, deleted) {
				s.applyStep(step)
				return true
			}
		}
		// Box bIdx intersects 3 columns:
		for cIdx := (bIdx % 3) * 3; cIdx < (bIdx%3)*3+3; cIdx++ {
			step := NewStep(kindSueDeCoq)
			if s.checkSueDeCoq(box, s.columns[cIdx], step, deleted) {
				s.applyStep(step)
				return true
			}
		}
	}

	return false
}

func (s *Solver) checkSueDeCoq(box, line *House, step *SolutionStep, deleted []bitset.BitSet16) bool {
	// 1. Identify intersection cells S (unsolved)
	var sCells []*puzzle.Cell
	for _, cell := range box.Cells {
		if (line.Kind == kindRow && cell.Row == line.Index) ||
			(line.Kind == kindColumn && cell.Col == line.Index) {
			if !cell.IsSolved() {
				sCells = append(sCells, cell)
			}
		}
	}

	if len(sCells) < 2 {
		return false
	}

	// 2. Identify potential cells in Box and Line (outside intersection)
	var bPossible []*puzzle.Cell
	for _, cell := range box.Cells {
		if cell.IsSolved() {
			continue
		}
		inS := false
		for _, sc := range sCells {
			if cell.Index() == sc.Index() {
				inS = true
				break
			}
		}
		if !inS {
			bPossible = append(bPossible, cell)
		}
	}

	var lPossible []*puzzle.Cell
	for _, cell := range line.Cells {
		if cell.IsSolved() {
			continue
		}
		inS := false
		for _, sc := range sCells {
			if cell.Index() == sc.Index() {
				inS = true
				break
			}
		}
		if !inS {
			lPossible = append(lPossible, cell)
		}
	}

	found := false
	// 3. Try all combinations of subsetB (0 to 2 cells) and subsetL (0 to 2 cells)
	for bSize := 0; bSize <= 2 && bSize <= len(bPossible); bSize++ {
		var checkB func(startB int, currentB []*puzzle.Cell)
		checkB = func(startB int, currentB []*puzzle.Cell) {
			if found {
				return
			}
			if len(currentB) == bSize {
				var candB bitset.BitSet16
				for _, bc := range currentB {
					candB.Union(bc.Candidates.Difference(deleted[bc.Index()]))
				}

				for lSize := 0; lSize <= 2 && lSize <= len(lPossible); lSize++ {
					if found {
						break
					}
					if bSize == 0 && lSize == 0 {
						continue
					}

					var checkL func(startL int, currentL []*puzzle.Cell)
					checkL = func(startL int, currentL []*puzzle.Cell) {
						if found {
							return
						}
						if len(currentL) == lSize {
							var candL bitset.BitSet16
							for _, lc := range currentL {
								candL.Union(lc.Candidates.Difference(deleted[lc.Index()]))
							}

							// S = sCells ∪ currentB ∪ currentL
							numCells := len(sCells) + len(currentB) + len(currentL)

							var currentCandI bitset.BitSet16
							for _, sc := range sCells {
								currentCandI.Union(sc.Candidates.Difference(deleted[sc.Index()]))
							}

							candS := bitset.Union(currentCandI, candB, candL)

							if candS.Size() == numCells {
								if s.applySueDeCoq(sCells, currentB, currentL, currentCandI, candB, candL, box, line, step, deleted) {
									found = true
								}
							}
							return
						}

						for i := startL; i < len(lPossible); i++ {
							checkL(i+1, append(currentL, lPossible[i]))
						}
					}
					checkL(0, nil)
				}
				return
			}

			for i := startB; i < len(bPossible); i++ {
				checkB(i+1, append(currentB, bPossible[i]))
			}
		}
		checkB(0, nil)
		if found {
			break
		}
	}

	return found
}

func (s *Solver) applySueDeCoq(
	sCells, bCells, lCells []*puzzle.Cell,
	candI, candB, candL bitset.BitSet16,
	box, line *House,
	step *SolutionStep,
	deleted []bitset.BitSet16,
) bool {
	// S = sCells ∪ bCells ∪ lCells
	candS := bitset.Union(candI, candB, candL)

	foundElim := false

	// Elimination 1: Digits in (candI ∪ candB) \ candL from Box \ (sCells ∪ bCells)
	elimB := bitset.Union(candI, candB).Difference(candL)
	if !elimB.Empty() {
		for _, cell := range box.Cells {
			if cell.IsSolved() {
				continue
			}
			// Is it in sCells or bCells?
			inS := false
			for _, sc := range sCells {
				if cell.Index() == sc.Index() {
					inS = true
					break
				}
			}
			if !inS {
				for _, bc := range bCells {
					if cell.Index() == bc.Index() {
						inS = true
						break
					}
				}
			}
			if inS {
				continue
			}

			common := cell.Candidates.Intersection(elimB)
			for v := range common.All() {
				if !deleted[cell.Index()].Contains(v) {
					step.DeleteCandidate(cell.Index(), v)
					deleted[cell.Index()].Add(v)
					foundElim = true
				}
			}
		}
	}

	// Elimination 2: Digits in (candI ∪ candL) \ candB from Line \ (sCells ∪ lCells)
	elimL := bitset.Union(candI, candL).Difference(candB)
	if !elimL.Empty() {
		for _, cell := range line.Cells {
			if cell.IsSolved() {
				continue
			}
			// Is it in sCells or lCells?
			inS := false
			for _, sc := range sCells {
				if cell.Index() == sc.Index() {
					inS = true
					break
				}
			}
			if !inS {
				for _, lc := range lCells {
					if cell.Index() == lc.Index() {
						inS = true
						break
					}
				}
			}
			if inS {
				continue
			}

			common := cell.Candidates.Intersection(elimL)
			for v := range common.All() {
				if !deleted[cell.Index()].Contains(v) {
					step.DeleteCandidate(cell.Index(), v)
					deleted[cell.Index()].Add(v)
					foundElim = true
				}
			}
		}
	}

	if foundElim {
		step.WithBases(box, line)
		for _, v := range candS.Values() {
			found := false
			for _, ev := range step.values {
				if ev == v {
					found = true
					break
				}
			}
			if !found {
				step.values = append(step.values, v)
			}
		}

		cells := append(append(append([]*puzzle.Cell{}, sCells...), bCells...), lCells...)
		for _, c := range cells {
			idx := c.Index()
			found := false
			for _, eidx := range step.indices {
				if eidx == idx {
					found = true
					break
				}
			}
			if !found {
				step.indices = append(step.indices, idx)
			}
		}
		return true
	}
	return false
}
