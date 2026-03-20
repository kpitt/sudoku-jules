package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kpitt/sudoku/internal/puzzle"
	"github.com/kpitt/sudoku/internal/solver"
	"github.com/mattn/go-isatty"
)

func main() {
	var (
		inputFile   string
		bruteForce  bool
		logLevel    string
		testMode    bool
	)

	flag.StringVar(&inputFile, "file", "", "input file (defaults to stdin)")
	flag.BoolVar(&bruteForce, "brute-force", true, "enable brute-force fallback")
	flag.StringVar(&logLevel, "log-level", "info", "log level (info, debug, trace)")
	flag.BoolVar(&testMode, "test", false, "regression mode")
	flag.Parse()

	if testMode {
		runRegression(inputFile)
		return
	}

	var input io.Reader
	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			fatalError(err.Error())
		}
		defer f.Close()
		input = f
	} else {
		input = os.Stdin
		if isStdinTTY() {
			fmt.Println("Enter initial puzzle as 9 lines of 9 characters.")
			fmt.Println("Use any character other than the digits 1-9 for empty cells.")
			fmt.Println("(Ctrl+D to finish on Unix/Linux, Ctrl+Z then Enter on Windows):")
		}
	}

	// Read all from input to support FromString's format detection
	data, err := io.ReadAll(input)
	if err != nil {
		fatalError(err.Error())
	}

	p, err := puzzle.FromString(string(data))
	if err != nil {
		fatalError(err.Error())
	}

	color.HiBlue("Original Puzzle:")
	p.Print()
	fmt.Println()

	liveLog := false
	enableDebug := false
	switch strings.ToLower(logLevel) {
	case "debug":
		liveLog = true
	case "trace":
		liveLog = true
		enableDebug = true
	}

	opts := &solver.Options{
		LiveLog:          liveLog,
		EnableDebug:      enableDebug,
		EnableBruteForce: bruteForce,
	}
	s := solver.NewSolver(p, opts)
	s.Solve()

	if p.IsSolved() {
		fmt.Printf("%s (%v)\n\n",
			color.HiGreenString("✓ Solved successfully"), s.SolveTime)
		if s.IsNonUnique {
			fmt.Printf("%s\n\n", color.HiYellowString("Warning: Puzzle is non-unique (has multiple solutions)"))
		}
		color.HiBlue("Solution:")
		p.Print()
	} else {
		if s.IsUnsolvable {
			fmt.Printf("%s (%v)\n\n",
				color.HiRedString("✗ Puzzle is unsolvable"), s.SolveTime)
		} else {
			fmt.Printf("%s (%v)\n\n",
				color.HiRedString("✗ Failed to solve"), s.SolveTime)
		}
		color.HiBlue("Partial Solution:")
		p.PrintCandidateGrid()
		fmt.Println()
		p.PrintUnsolvedCounts()
	}

	// Only print solution if steps were not already live-logged.
	if !opts.LiveLog {
		fmt.Println()
		s.PrintSolution()
	}
}

func runRegression(filename string) {
	if filename == "" {
		fatalError("regression mode requires an input file via -file")
	}

	stats, err := solver.RunRegressionFile(filename)
	if err != nil {
		fatalError(err.Error())
	}

	fmt.Printf("\nRegression Results: %s\n", filename)
	fmt.Printf("  Total Tests:  %d\n", stats.Total)
	fmt.Printf("  Passed Tests: %d\n", stats.Passed)
	fmt.Printf("  Failed Tests: %d\n", stats.Failed)
	fmt.Printf("  Duration:     %v\n", stats.End.Sub(stats.Start))

	if stats.Failed > 0 {
		os.Exit(1)
	}
}

func isStdinTTY() bool {
	return isTerminal(os.Stdin)
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

func fatalError(msg string) {
	fmt.Fprintln(os.Stderr, color.HiRedString("error: %s", msg))
	os.Exit(1)
}
