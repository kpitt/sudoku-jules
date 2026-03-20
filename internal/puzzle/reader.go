package puzzle

import (
	"io"
	"os"
	"strings"
	"unicode"
)

// FromFile reads a Sudoku puzzle from the specified file.
func FromFile(f *os.File) (*Puzzle, error) {
	var buf strings.Builder
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	_, _ = buf.Write(data)
	return FromString(buf.String())
}

// FromString creates a new puzzle from the given puzzle string.
// Empty lines and comments starting with a '#' are ignored.
// The following formats are supported:
//
// A single 81-digit string ('0' or '.' for an empty cell)
//
// A 9x9 grid of cells, with optional whitespace and border characters.
// Possible formats include:
//
//	```
//	16.54..7.   or   *-----------*   or   +-----+-----+-----+
//	..8..1.3.        |16.|54.|.7.|        |1 6 .|5 4 .|. 7 .|
//	.3.8.....        |..8|..1|.3.|        |. . 8|. . 1|. 3 .|
//	7...5..69        |.3.|8..|...|        |. 3 .|8 . .|. . .|
//	6..9.2.57        |---+---+---|        +-----+-----+-----+
//	.........        |7..|.5.|.69|        |7 . .|. 5 .|. 6 9|
//	....3..4.        |6..|9.2|.57|        |6 . .|9 . 2|. 5 7|
//	.......16        |...|...|...|        |. . .|. . .|. . .|
//	...1645..        |---+---+---|        +-----+-----+-----+
//	                 |...|.3.|.4.|        |. . .|. 3 .|. 4 .|
//	                 |...|...|.16|        |. . .|. . .|. 1 6|
//	                 |...|164|5..|        |. . .|1 6 4|5 . .|
//	                 *-----------*        +-----+-----+-----+
//	```
func FromString(s string) (*Puzzle, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, ":") {
		return FromHodokuString(s)
	}

	p := NewPuzzle()
	// Determine the newline sequence used in the input string, and split the
	// string into lines based on that sequence.
	sep := "\n"
	if strings.Contains(s, "\r\n") {
		sep = "\r\n"
	} else if strings.Contains(s, "\r") {
		sep = "\r"
	}
	lines := strings.Split(s, sep)
	i := 0
	for _, line := range lines {
		// Remove comments and extra whitespace from the line.
		line, _, _ = strings.Cut(line, "#")
		line = strings.TrimSpace(line)
		// Skip empty lines and border lines.
		// Any line with a sequence of 3 dashes ("---") is a border line.
		if line == "" || strings.Contains(line, "---") {
			continue
		}

		for pos, c := range line {
			// Discard whitespace and '|' border characters.
			if unicode.IsSpace(c) || c == '|' {
				continue
			}
			if i >= 81 {
				// We've already filled all 81 cells and skipped any whitespace
				// and border characters that might follow the last cell, so we
				// shouldn't have any more characters to process.
				return nil, errPuzzleFormat("extraneous characters: %q", line[pos:])
			}
			if c != '.' && !unicode.IsDigit(c) {
				return nil, errPuzzleFormat("invalid character: %c", c)
			}
			// Place a given for digits '1'-'9' and advance the index.
			// A '0' or '.' is an empty cell, so just advance the index without
			// placing a given.
			if c >= '1' && c <= '9' {
				val := int(c - '0')
				if !p.Cell(i).HasCandidate(val) {
					return nil, errPuzzleState("given value %d is not a candidate for %s",
						val, FormatCell(i))
				}
				p.GivenValue(i, val)
			}
			i++
		}

		if i >= 81 {
			// If we've filled all 81 cells, then we're done.
			// Just ignore any remaining lines.
			break
		}
	}

	if i < 81 {
		return nil, errPuzzleFormat("not enough cells")
	}

	return p, nil
}

// FromHodokuString creates a new puzzle from a Hodoku Library Format string:
// :technique:candidates:givens:deleted:eliminations:placements:extra
func FromHodokuString(s string) (*Puzzle, error) {
	if !strings.HasPrefix(s, ":") {
		return nil, errPuzzleFormat("not a Hodoku format string")
	}
	parts := strings.Split(s, ":")
	if len(parts) < 8 {
		return nil, errPuzzleFormat("invalid Hodoku format string (expected 8 parts, got %d)", len(parts))
	}

	givensStr := parts[3]
	deletedStr := parts[4]

	p := NewPuzzle()

	// Parse givens
	// Example: +6+3+42+8+9...
	// A '+' before a digit means it's a placed value, not a given.
	// '.' or '0' means empty.
	idx := 0
	for i := 0; i < len(givensStr) && idx < 81; i++ {
		c := givensStr[i]
		switch {
		case c == '+':
			i++
			if i >= len(givensStr) {
				return nil, errPuzzleFormat("invalid Hodoku givens: trailing +")
			}
			c = givensStr[i]
			if !unicode.IsDigit(rune(c)) || c == '0' {
				return nil, errPuzzleFormat("invalid Hodoku givens: + must be followed by 1-9")
			}
			val := int(c - '0')
			p.PlaceValue(idx, val)
			idx++
		case c == '.' || c == '0':
			idx++
		case unicode.IsDigit(rune(c)):
			val := int(c - '0')
			p.GivenValue(idx, val)
			idx++
		}
	}

	if idx < 81 {
		return nil, errPuzzleFormat("not enough cells in Hodoku givens (got %d, expected 81)", idx)
	}

	// Parse deleted candidates
	// Format: <candidate><line><col> separated by spaces
	if deletedStr != "" {
		deletedParts := strings.Fields(deletedStr)
		for _, dp := range deletedParts {
			if len(dp) != 3 {
				return nil, errPuzzleFormat("invalid deleted candidate format: %s", dp)
			}
			val := int(dp[0] - '0')
			row := int(dp[1] - '1')
			col := int(dp[2] - '1')
			if val < 1 || val > 9 || row < 0 || row > 8 || col < 0 || col > 8 {
				return nil, errPuzzleFormat("invalid deleted candidate values: %s", dp)
			}
			p.Get(row, col).RemoveCandidate(val)
		}
	}

	return p, nil
}
