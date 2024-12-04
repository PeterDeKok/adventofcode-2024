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
	word := "MAS"
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
			found := 0
			for i, got := range getAlldirections(lines, y, x, word, 1) {
				if got == word {
					if build.DEBUG {
						fmt.Printf("  [word %d]: OK %s\n", i, got)
					}
					found++
				} else {
					if build.DEBUG {
						fmt.Printf("  [word %d]: %s\n", i, got)
					}
				}
			}

			if found >= 2 {
				sum++
			}
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func getAlldirections(lines [][]rune, y, x int, word string, center int) []string {
	if y < 0 || y >= len(lines) || x < 0 || x >= len(lines[y]) {
		return []string{}
	}

	if len(word) <= 1 {
		// Prevent getting the same char 8 times
		return []string{string(lines[y][x])}
	}

	words := make([]string, 0, 8)

	dirs := [][2]int{
		{-1, -1},
		{-1, 1},
		{1, 1},
		{1, -1},
	}

	for _, d := range dirs {
		if build.DEBUG {
			fmt.Println("")
		}
		if w := getWord(lines, y, x, d[0], d[1], word, center); len(w) > 0 {
			words = append(words, w)
		}
	}

	return words
}

func getWord(lines [][]rune, y, x, dy, dx int, word string, center int) string {
	result := ""

	for i := range word {
		y2, x2 := (y-center*dy)+i*dy, (x-center*dx)+i*dx

		if y2 < 0 || y2 >= len(lines) || x2 < 0 || x2 >= len(lines[y2]) {
			return ""
		}

		if build.DEBUG {
			fmt.Printf("      TRY: %d,%d [%d,%d]: %s\n", y2, x2, dy, dx, string(lines[y2][x2]))
		}

		result += string(lines[y2][x2])
	}

	return result
}
