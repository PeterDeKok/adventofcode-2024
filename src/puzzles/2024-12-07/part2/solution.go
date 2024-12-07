package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/distribute"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/math"
	"strconv"
	"strings"
)

func Solution(ctx context.Context, il *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	err := distribute.Lines(ctx, il, rd, func(ctx context.Context, i int, l string) error {
		if v, ok, err := handleLineRecursion(il, i, l); err != nil {
			return err
		} else if ok {
			sum += v
		}

		return nil
	})

	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func handleLineRecursion(ll *logger.IterationLogger, lineNr int, line string) (int, bool, error) {
	expected, operands, err := parseLine(line)
	if err != nil || len(operands) == 0 {
		return 0, false, err
	} else if len(operands) == 1 {
		return expected, operands[0] == expected, nil
	}

	prev := make([]interface{}, 0, 100)
	prev = append(prev, operands[0])

	var reducer func(carry int, operands []int, prev []interface{}) bool
	reducer = func(carry int, operands []int, prev []interface{}) bool {
		if len(operands) == 0 {
			if carry == expected {
				if build.DEBUG {
					ll.LogDebugf(lineNr, "found a solution: %d = %s", expected, fmt.Sprint(prev...))
				}
				return true
			} else {
				if build.DEBUG {
					ll.LogDebugf(lineNr, "NOT a solution: %d != %s", expected, fmt.Sprint(prev...))
				}
				return false
			}
		}

		if carry > expected {
			if build.DEBUG {
				ll.LogDebugf(lineNr, "NOT a solution: %d < %s ...", expected, fmt.Sprint(prev...))
			}
			return false
		}

		v := operands[0]
		next := operands[1:]

		return reducer(carry+v, next, append(prev[:], "+", v)) ||
			reducer(carry*v, next, append(prev[:], "*", v)) ||
			reducer(math.Concat(carry, v), next, append(prev[:], "||", v))
	}

	return expected, reducer(operands[0], operands[1:], prev), nil
}

func parseLine(line string) (expected int, operands []int, err error) {
	nrs := strings.Fields(line)

	expected, err = strconv.Atoi(strings.Trim(nrs[0], ":"))
	if err != nil {
		return expected, nil, fmt.Errorf("invalid expected")
	}

	if len(nrs) <= 1 {
		return expected, nil, fmt.Errorf("no operands (want %d)", expected)
	}

	if len(nrs)-2 > 32 {
		return expected, nil, fmt.Errorf("%d opSlots does not fit in uint32", len(nrs)-2)
	}

	operands = make([]int, len(nrs)-1)

	for i, operand := range nrs[1:] {
		if operands[i], err = strconv.Atoi(operand); err != nil {
			return expected, nil, fmt.Errorf("invalid operand [%d] %s", i, operand)
		}
	}

	return expected, operands, nil
}
