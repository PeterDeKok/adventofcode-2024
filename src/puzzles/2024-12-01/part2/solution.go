package main

import (
	"context"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"strings"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	listA := make([]int, 0, 1024)
	listB := make([]int, 0, 1024)
	countListB := make(map[int]int, len(listB))

	for _, line := range input.LineReader(rd) {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)

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

	for _, v := range listB {
		if c, ok := countListB[v]; ok {
			countListB[v] = c + v
		} else {
			countListB[v] = v
		}
	}

	for _, v := range listA {
		if c, ok := countListB[v]; ok {
			sum += c
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}
