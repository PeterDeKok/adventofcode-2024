package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/math"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var safe int

	for i, line := range input.LineIsIntSliceReader(rd) {
		lineResult, ok, err := handleLine(line)
		if !ok && err != nil {
			panic(fmt.Errorf("failed to handle line %d: %v", i, err))
		} else if !ok {
			panic(fmt.Errorf("failed to handle line %d: unknown error", i))
		}

		if build.DEBUG {
			if err != nil {
				fmt.Printf("[line %d]: %v | %v\n", i, lineResult, err)
			} else {
				fmt.Printf("[line %d]: %v\n", i, lineResult)
			}
		}

		if lineResult {
			safe++
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(safe))); err != nil {
		return err
	}

	return nil
}

func handleLine(line []int) (bool, bool, error) {
	if len(line) < 2 {
		return false, false, fmt.Errorf("invalid report, len = %d", len(line))
	}

	var dirUp, dirDown, dirZero int

	for i, v := range line[1:] {
		if v > line[i] {
			dirUp++
		} else if v < line[i] {
			dirDown++
		} else {
			dirZero++
		}
	}

	if dirZero > 0 {
		return false, true, fmt.Errorf("zero diff")
	} else if dirUp > 0 && dirDown > 0 {
		return false, true, fmt.Errorf("not strictly increasing or decreasing")
	}

	for i, v := range line[1:] {
		if diff := math.AbsDiff(v, line[i]); diff > 3 {
			return false, true, fmt.Errorf("diff too large d[%d, %d] = %d", v, line[i], diff)
		}
	}

	return true, true, nil
}
