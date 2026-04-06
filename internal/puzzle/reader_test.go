package puzzle

import (
	"os"
	"strings"
	"testing"
)

func TestFromFile_LargeInput(t *testing.T) {
	// Create a large temporary file (2MB)
	tmpFile, err := os.CreateTemp("", "large_puzzle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	data := make([]byte, 2*1024*1024)
	for i := range data {
		data[i] = 'A'
	}
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	_, err = FromFile(tmpFile)
	if err == nil {
		t.Error("expected error for large input, got nil")
	} else if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Errorf("expected maximum size error, got %v", err)
	}
}

func TestFromFile_ValidInput(t *testing.T) {
	puzzleStr := "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
	tmpFile, err := os.CreateTemp("", "valid_puzzle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(puzzleStr)); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	p, err := FromFile(tmpFile)
	if err != nil {
		t.Fatalf("expected no error for valid input, got %v", err)
	}
	if p == nil {
		t.Fatal("expected puzzle, got nil")
	}
}

func TestFromFile_Truncation(t *testing.T) {
	// Create a file slightly larger than maxPuzzleSize
	size := maxPuzzleSize + 100
	data := strings.Repeat("A", size)

	tmpFile, err := os.CreateTemp("", "truncated_puzzle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(data)); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	// We check that truncation error is handled properly.
	_, err = FromFile(tmpFile)
	if err == nil {
		t.Error("expected error for truncated input, got nil")
	} else if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Errorf("expected maximum size error, got %v", err)
	}
}
