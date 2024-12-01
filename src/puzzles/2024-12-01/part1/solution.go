package main

import (
	"context"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"slices"
	"strconv"
	"strings"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	listA := make([]int, 0, 1024)
	listB := make([]int, 0, 1024)

	for _, line := range input.LineReader(rd) {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		if nrA, err := strconv.Atoi(parts[0]); err == nil {
			listA = append(listA, nrA)
		}

		if nrB, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			listB = append(listB, nrB)
		}
	}

	if len(listA) != len(listB) {
		panic("Incomatible lengts not yet implemented")
	}

	slices.Sort(listA)
	slices.Sort(listB)

	for i := 0; i < len(listA); i++ {
		if listA[i] > listB[i] {
			sum += listA[i] - listB[i]
		} else {
			sum += listB[i] - listA[i]
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}
