package solver

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kpitt/sudoku/internal/puzzle"
)

// RegressionStats holds the results of a regression test run.
type RegressionStats struct {
	Total  int
	Passed int
	Failed int
	Start  time.Time
	End    time.Time
}

// RunRegressionFile runs all test cases in the specified file.
func RunRegressionFile(filename string) (*RegressionStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := &RegressionStats{
		Start: time.Now(),
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		stats.Total++
		if pass, err := runTestCase(line, lineNum); err != nil {
			fmt.Printf("Line %d: ERROR: %v\n", lineNum, err)
			stats.Failed++
		} else if pass {
			stats.Passed++
		} else {
			fmt.Printf("Line %d: FAILED: %s\n", lineNum, line)
			stats.Failed++
		}
	}

	stats.End = time.Now()
	return stats, scanner.Err()
}

func runTestCase(line string, lineNum int) (bool, error) {
	p, err := puzzle.FromHodokuString(line)
	if err != nil {
		return false, fmt.Errorf("parse error: %v", err)
	}

	parts := strings.Split(line, ":")
	// :technique:candidates:givens:deleted:eliminations:placements:extra
	targetTechID := parts[1]
	// Strip suffixes like -1, -2, -x
	if idx := strings.Index(targetTechID, "-"); idx != -1 {
		targetTechID = targetTechID[:idx]
	}

	targetCandidatesStr := parts[2]
	expectedElimsStr := parts[5]
	expectedPlacementsStr := parts[6]

	targetCandidates := parseHodokuDigits(targetCandidatesStr)

	targetKind, ok := hodokuIDToKind[targetTechID]
	if !ok {
		return false, fmt.Errorf("unknown Hodoku technique ID: %s", targetTechID)
	}

	opts := &Options{
		// Don't enable brute force or debug during regression
		DisableAutomaticSingles: true,
	}
	s := NewSolver(p, opts)
	s.processInitialValues()

	// Find the technique in the solver's list
	var tech *Technique
	if int(targetKind) < len(s.techniques) {
		tech = &s.techniques[targetKind]
	}

	if tech == nil || tech.Check == nil {
		// If it's Naked Single, it's handled automatically during initialization or by applyStep.
		// We need to check if any step was already recorded.
		if targetKind != kindNakedSingle {
			return false, fmt.Errorf("technique %s (ID %s) not implemented or has no Check function", techName(targetKind), targetTechID)
		}
	} else {
		// Run the specific technique
		tech.Check()
	}

	// Now verify the results
	return verifyResults(s, targetCandidates, expectedElimsStr, expectedPlacementsStr)
}

func verifyResults(s *Solver, targetCandidates []int, expectedElimsStr, expectedPlacementsStr string) (bool, error) {
	expectedElims := parseHodokuCandidates(expectedElimsStr)
	expectedPlacements := parseHodokuCandidates(expectedPlacementsStr)

	// Collect all eliminations and placements from the solver's recorded steps
	actualElims := make([]Candidate, 0)
	actualPlacements := make([]Candidate, 0)

	for _, step := range s.solution {
		// Filter by target candidates
		for _, dc := range step.deletedCandidates {
			if containsDigit(targetCandidates, dc.Value) {
				actualElims = append(actualElims, dc)
			}
		}
		if step.IsSingle() {
			for i, idx := range step.indices {
				val := step.values[i]
				if containsDigit(targetCandidates, val) {
					actualPlacements = append(actualPlacements, Candidate{Index: idx, Value: val})
				}
			}
		}
	}

	// Compare
	if !compareCandidates(actualElims, expectedElims) {
		fmt.Printf("  Eliminations mismatch (target candidates %v):\n    Expected: %v\n    Actual:   %v\n", targetCandidates, formatCandidates(expectedElims), formatCandidates(actualElims))
		return false, nil
	}

	if !compareCandidates(actualPlacements, expectedPlacements) {
		fmt.Printf("  Placements mismatch (target candidates %v):\n    Expected: %v\n    Actual:   %v\n", targetCandidates, formatCandidates(expectedPlacements), formatCandidates(actualPlacements))
		return false, nil
	}

	return true, nil
}

func parseHodokuDigits(s string) []int {
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	res := make([]int, 0, len(parts))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		// Some might be digits, some might be cell refs?
		// For Singles, it's just the digit.
		// For others, it might be more.
		// Let's just try to parse digits from the string.
		for _, c := range p {
			if c >= '1' && c <= '9' {
				res = append(res, int(c-'0'))
			}
		}
	}
	return res
}

func containsDigit(digits []int, d int) bool {
	if len(digits) == 0 {
		return true // If no target candidates specified, include all
	}
	for _, x := range digits {
		if x == d {
			return true
		}
	}
	return false
}

func parseHodokuCandidates(s string) []Candidate {
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	res := make([]Candidate, 0, len(parts))
	for _, p := range parts {
		if len(p) != 3 {
			continue
		}
		val := int(p[0] - '0')
		row := int(p[1] - '1')
		col := int(p[2] - '1')
		res = append(res, Candidate{Index: row*9 + col, Value: val})
	}
	return res
}

func compareCandidates(a, b []Candidate) bool {
	if len(a) != len(b) {
		return false
	}
	// Sort both to compare
	sortCandidates(a)
	sortCandidates(b)
	for i := range a {
		if a[i].Index != b[i].Index || a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}

func sortCandidates(c []Candidate) {
	sort.Slice(c, func(i, j int) bool {
		if c[i].Index != c[j].Index {
			return c[i].Index < c[j].Index
		}
		return c[i].Value < c[j].Value
	})
}

func formatCandidates(cs []Candidate) string {
	if len(cs) == 0 {
		return "none"
	}
	sortCandidates(cs)
	var parts []string
	for _, c := range cs {
		parts = append(parts, fmt.Sprintf("%d@%s", c.Value, puzzle.FormatCell(c.Index)))
	}
	return strings.Join(parts, " ")
}

func techName(k techniqueKind) string {
	names := []string{
		"Full House", "Naked Single", "Hidden Single", "Locked Candidates (Pointing)", "Locked Candidates (Claiming)",
		"Naked Pair", "Naked Triple", "Hidden Pair", "Hidden Triple", "Naked Quadruple", "Hidden Quadruple",
		"X-Wing", "Swordfish", "Jellyfish", "Remote Pair", "BUG+1", "Skyscraper", "2-String Kite",
		"Empty Rectangle", "W-Wing", "XY-Wing", "XYZ-Wing", "Avoidable Rectangle",
		"Unique Rectangle Type 1", "Unique Rectangle Type 2", "Unique Rectangle Type 3", "Unique Rectangle Type 4",
		"Hidden Rectangle", "Finned X-Wing", "Finned Swordfish", "Finned Jellyfish", "Sue de Coq",
		"Simple Coloring", "Multi-Coloring", "X-Chain", "XY-Chain", "Nice Loop", "AIC", "ALS-XZ", "ALS-XY-Wing", "Brute Force",
	}
	if int(k) < len(names) {
		return names[k]
	}
	return "Unknown"
}

var hodokuIDToKind = map[string]techniqueKind{
	"0000": kindFullHouse,
	"0003": kindNakedSingle,
	"0002": kindHiddenSingle,
	"0110": kindNakedPair,
	"0111": kindNakedTriple,
	"0100": kindLockedCandidatesPointing,
	"0101": kindLockedCandidatesClaiming,
	"0200": kindNakedPair,
	"0201": kindNakedTriple,
	"0202": kindNakedQuadruple,
	"0210": kindHiddenPair,
	"0211": kindHiddenTriple,
	"0212": kindHiddenQuadruple,
	"0300": kindXWing,
	"0301": kindSwordfish,
	"0302": kindJellyfish,
	"0310": kindFinnedXWing,
	"0320": kindFinnedXWing, // Sashimi handled by Finned
	"0311": kindFinnedSwordfish,
	"0321": kindFinnedSwordfish,
	"0312": kindFinnedJellyfish,
	"0322": kindFinnedJellyfish,
	"0400": kindSkyscraper,
	"0401": kindTwoStringKite,
	"0402": kindEmptyRectangle,
	"0701": kindXChain,
	"0702": kindAIC,
	"0708": kindAIC,
	"0707": kindNiceLoop,
	"0703": kindRemotePair,
	"0610": kindBUG,
	"0803": kindWWing,
	"0800": kindXYWing,
	"0801": kindXYZWing,
	"0901": kindALSXZ,
	"0902": kindALSXYWing,
	"0600": kindUniqueRectangle1,
	"0601": kindUniqueRectangle2,
	"0602": kindUniqueRectangle3,
	"0603": kindUniqueRectangle4,
	"0606": kindHiddenRectangle,
	"0607": kindAvoidableRectangle,
	"0608": kindAvoidableRectangle, // AR Type 2 handled by same func
	"1101": kindSueDeCoq,
}
