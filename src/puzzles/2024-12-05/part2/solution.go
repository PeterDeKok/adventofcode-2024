package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"slices"
	"strconv"
	"strings"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int
	var ordersI int

	lookup := make(map[int][]int)
	stage := "lookup"

	for i, line := range input.LineReader(rd) {
		if line == "" {
			if build.DEBUG {
				fmt.Printf("\n\n%v\n\n", lookup)
			}
			stage = "orders"
			continue
		}

		if build.DEBUG {
			fmt.Printf("[parse: %d | %s]: %s\n", i, stage, line)
		}

		switch stage {
		case "lookup":
			v := strings.Split(line, "|")
			if len(v) != 2 {
				panic("Invalid constriant")
			}

			key, err := strconv.Atoi(v[0])
			if err != nil {
				panic("Failed to parse key")
			}

			value, err := strconv.Atoi(v[1])
			if err != nil {
				panic("Failed to parse key")
			}

			constraints, ok := lookup[key]
			if !ok {
				constraints = make([]int, 0, 10)
				lookup[key] = constraints
			}

			lookup[key] = append(lookup[key], value)
			if build.DEBUG {
				fmt.Printf("  [rule: %d]: %s\n", i, line)
			}
		default:
			_, ok := handlePrintLine(lookup, line)

			if !ok {
				sum += fixPrintLine(lookup, line)
			}

			if build.DEBUG {
				fmt.Printf("  [order: %d]: %v %s\n", ordersI, ok, line)
			}

			ordersI++
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func handlePrintLine(lookup map[int][]int, line string) (int, bool) {
	strs := strings.Split(line, ",")
	nrs := make([]int, len(strs))
	nrsLookup := make(map[int]int, len(strs))

	for i, v := range strs {
		nr, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Errorf("failed to parse value %d", i))
		}
		nrs[i] = nr

		constraints, ok := lookup[nr]
		if !ok {
			if build.DEBUG {
				fmt.Printf("    [check: %d] no constraint\n", nr)
			}
			nrsLookup[nr] = nr
			continue
		}

		if build.DEBUG {
			fmt.Printf("    [check: %d]: against %v\n", nr, constraints)
		}

		for _, before := range constraints {
			if _, wrong := nrsLookup[before]; wrong {
				if build.DEBUG {
					fmt.Printf("    [check: %d]: not before %d\n", nr, before)
				}
				return 0, false
			} else {
				if build.DEBUG {
					fmt.Printf("    [check: %d]: checked against %d\n", nr, before)
				}
			}
		}

		nrsLookup[nr] = nr
	}

	mid := nrs[len(nrs)/2]

	if build.DEBUG {
		fmt.Printf("    [check]: OK: mid: %d\n", mid)
	}
	return mid, true
}

func fixPrintLine(lookup map[int][]int, line string) int {
	strs := strings.Split(line, ",")
	nrs := make([]int, len(strs))

	for i, v := range strs {
		nr, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Errorf("failed to parse value %d", i))
		}
		nrs[i] = nr
	}

	slices.SortFunc(nrs, func(a, b int) int {
		la, oka := lookup[a]
		lb, okb := lookup[b]

		if oka && slices.Contains(la, b) {
			return -1
		}

		if okb && slices.Contains(lb, a) {
			return 1
		}

		return 0
	})

	mid := nrs[len(nrs)/2]

	fmt.Printf("    [sorted]: mid: %d\n", mid)

	return mid
}
