package main

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
