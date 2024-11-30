package puzzle

const readmeContent = `# Advent of Code 2024 | Day {{day}} | Part {{part}}

Eric Wastl deserves the credit for creating these challenges year after year. Therefore, the problem description, input and Solutions won't be available in this repo.
However, a link to the problem description will be provided and the solution can be generated locally.

| Year        | 2024                    |
|-------------|-------------------------|
| Date        | 2024-12-{{daypadded}}              |
| Part        | {{part}}                       |
| Finished on |                         |

## Problem

[adventofcode.com/2024/day/{{day}}](https://adventofcode.com/2024/day/{{day}})

## Runtime

` + "```" + `

` + "```" + `
`

const solutionContent = `package main

import (
	"context"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	for _, line := range input.LineReader(rd) {
		panic(line)
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}
`

const solutionTestContent = `package main

import (
	"context"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strings"
	"testing"
	"time"
)

func TestSolution_Sample(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ll := logger.CreateIterationLogger(ctx)
	go ll.AllFmt()

	out := &strings.Builder{}

	fpI, err := os.Open("sample-input-1.txt")
	if err != nil {
		panic(err)
	}
	fpE, err := os.ReadFile("sample-expected-1.txt")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	if err := Solution(ctx, ll, fpI, out); err != nil {
		panic(err)
	} else if str := out.String(); str != string(fpE) {
		t.Fatalf("[FAIL] got %s, want %s\n", str, string(fpE))
	} else {
		t.Logf("[-OK-]: %s -> (%s)\n", str, time.Since(start))
	}
}

func TestSolution_Input(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ll := logger.CreateIterationLogger(ctx)
	go ll.AllFmt()

	out := &strings.Builder{}

	fpI, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	if err := Solution(ctx, ll, fpI, out); err != nil {
		panic(err)
	} else {
		t.Logf("[RSLT]: %s -> (%s)\n", out, time.Since(start))
	}
}

func BenchmarkSolution_Sample(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ll := logger.CreateIterationLogger(ctx)
	go ll.AllFmt()

	var input string

	if bInput, err := os.ReadFile("sample-input-1.txt"); err != nil {
		panic(err)
	} else {
		input = string(bInput)
	}

    for i := 0; i < b.N; i++ {
        b.StopTimer()

		out := &strings.Builder{}
		ir := strings.NewReader(input)

        b.StartTimer()

		if err := Solution(ctx, ll, ir, out); err != nil {
			panic(err)
		}
    }
}

func BenchmarkSolution_Input(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ll := logger.CreateIterationLogger(ctx)
	go ll.AllFmt()

	var input string

	if bInput, err := os.ReadFile("input.txt"); err != nil {
		panic(err)
	} else {
		input = string(bInput)
	}

    for i := 0; i < b.N; i++ {
        b.StopTimer()

		out := &strings.Builder{}
		ir := strings.NewReader(input)

        b.StartTimer()

		if err := Solution(ctx, ll, ir, out); err != nil {
			panic(err)
		}
    }
}
`

const statsContent = `{
  "nr": {{part}},
  "solutions": [],
  "star": null
}
`
