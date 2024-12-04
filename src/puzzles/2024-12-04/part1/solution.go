package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	word := "XMAS"
	var sum int

	// Ensure all characters are properly indexed
	lines := make([][]rune, 0, 1024)
	for _, line := range input.LineReader(rd) {
		lines = append(lines, []rune(line))
	}

	for y, line := range lines {
		if build.DEBUG {
			fmt.Printf("\n[line %d]\n", y)
		}
		for x := 0; x < len(line); x++ {
			if build.DEBUG {
				fmt.Printf("  [char @ %d,%d]\n", y, x)
			}
			for i, got := range getAlldirections(lines, y, x, word) {
				if got == word {
					if build.DEBUG {
						fmt.Printf("  [word %d]: OK %s\n", i, got)
					}
					sum++
				} else {
					if build.DEBUG {
						fmt.Printf("  [word %d]: %s\n", i, got)
					}
				}
			}
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func getAlldirections(lines [][]rune, y, x int, word string) []string {
	if y < 0 || y >= len(lines) || x < 0 || x >= len(lines[y]) {
		return []string{}
	}

	if len(word) <= 1 {
		// Prevent getting the same char 8 times
		return []string{string(lines[y][x])}
	}

	words := make([]string, 0, 8)

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if x == 0 && y == 0 {
				continue
			}

			if w := getWord(lines, y, x, dy, dx, word); len(w) > 0 {
				words = append(words, w)
			}
		}
	}

	return words
}

func getWord(lines [][]rune, y, x, dy, dx int, word string) string {
	result := ""

	for i, _ := range word {
		y2, x2 := y+i*dy, x+i*dx

		if y2 < 0 || y2 >= len(lines) || x2 < 0 || x2 >= len(lines[y2]) {
			return ""
		}

		result += string(lines[y+i*dy][x+i*dx])
	}

	return result
}
