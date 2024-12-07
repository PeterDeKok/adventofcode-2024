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
	sum := 0

	err := distribute.Lines(ctx, il, rd, func(ctx context.Context, i int, l string) error {
		if v, ok, err := handleLine(i, l); err != nil {
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

func handleLine(lineNr int, line string) (int, bool, error) {
	if build.DEBUG {
		fmt.Printf("[%d] %s\n", lineNr, line)
	}

	nrs := strings.Fields(line)

	expected, err := strconv.Atoi(strings.Trim(nrs[0], ":"))

	if err != nil {
		return 0, false, fmt.Errorf("[%d] invalid expected", lineNr)
	}

	if len(nrs) <= 1 {
		return 0, false, fmt.Errorf("[%d] no operands (want %d)", lineNr, expected)
	}

	if len(nrs)-2 > 32 {
		return 0, false, fmt.Errorf("[%d] %d opSlots does not fit in uint32", lineNr, len(nrs)-2)
	}

	operands := make([]int, len(nrs)-1)

	for i, operand := range nrs[1:] {
		if operands[i], err = strconv.Atoi(operand); err != nil {
			return 0, false, fmt.Errorf("[%d] invalid operand %s", i, operand)
		}
	}

	operators := []func(a, b int) int{
		func(a, b int) int { return a + b },
		//		func(a, b int) int { return a - b },
		func(a, b int) int { return a * b },
		//		func(a, b int) int { return a / b },
	}
	opstr := []string{
		"+",
		//		"-",
		"*",
		//		"/",
	}

	if !math.IsPowerOfTwoUint32(uint32(len(opstr))) || len(opstr) != len(operators) {
		return 0, false, fmt.Errorf("[%d] nr of operands should be a power of 2 %d", lineNr, len(opstr))
	}

	opSlots := uint32(len(operands) - 1)
	lim := math.PowUint32(uint32(len(operators)), opSlots)
	lenOp := uint32(len(operators))
	opBit := lenOp - 1

	if build.DEBUG {
		fmt.Printf("[%d]   opSlots=%d lim=%d lenOp=%d opBit=%d\n", lineNr, opSlots, lim, lenOp, opBit)
	}

	for j := uint32(0); j < lim; j++ {
		carry := operands[0]
		str := fmt.Sprintf("%d", carry)
		if build.DEBUG {
			fmt.Printf("[%d]   j: %d\n", lineNr, j)
		}

		for k := uint32(0); k < opSlots; k++ {
			opi := (j >> (k)) & opBit
			if build.DEBUG {
				fmt.Printf("[%d]     j: %d k: %d, op: %d (%s)\n", lineNr, j, k, opi, opstr[opi])
			}
			carry = operators[opi](carry, operands[k+1])
			str += fmt.Sprintf(" %s %d (=%d)", opstr[opi], operands[k], carry)
		}

		if build.DEBUG {
			fmt.Printf("[%d]       %s\n", lineNr, str)
		}
		if carry == expected {
			if build.DEBUG {
				fmt.Printf("[%d]       found good combination %d\n", lineNr, carry)
			}
			return carry, true, nil
		}
	}

	return 0, false, nil
}
