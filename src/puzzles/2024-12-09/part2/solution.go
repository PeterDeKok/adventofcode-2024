package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"strings"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	fs, free, size := parseInput(rd)

	sum := defrag(fs, free, size)

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func defrag(fs, free []*Block, size int) int {
	if build.DEBUG {
		printfs(fs, size)
	}

	sum := 0

	for i, skip := len(fs)-1, 0; i >= 0; i-- {
		f := fs[i]

		for j := skip; j < len(free); j++ {
			fr := free[j]

			if fr != nil && fr.l >= f.l && f.pos > fr.pos {
				f.pos = fr.pos
				fr.l -= f.l
				fr.pos += f.l

				if j == skip && fr.l <= 0 {
					skip++
				}

				if build.DEBUG {
					printfs(fs, size)
				}

				break
			}
		}

		for k := 0; k < f.l; k++ {
			sum += (f.pos + k) * f.id

			if build.DEBUG {
				fmt.Printf("(id: %d * %d = %d  -> %d\n", f.id, f.pos+k, f.id*(f.pos+k), sum)
			}
		}
	}

	return sum
}

func printfs(fs []*Block, size int) {
	str := []rune(strings.Repeat(".", size))

	for _, f := range fs {
		for i := 0; i < f.l; i++ {
			str[f.pos+i] = rune('0' + f.id)
		}
	}

	fmt.Println(string(str))
}

func parseInput(rd io.Reader) (fs []*Block, free []*Block, size int) {
	file := true
	id, pos := 0, 0
	fs = make([]*Block, 0, 10000)
	free = make([]*Block, 0, 10000)

	for _, r := range input.CharReader(rd) {
		l := int(r - '0')

		if file {
			fs = append(fs, &Block{
				id:  id,
				pos: pos,
				l:   l,
			})
			id++
		} else {
			free = append(free, &Block{
				pos: pos,
				l:   l,
			})
		}

		pos += l
		file = !file
	}

	return fs, free, pos
}

type Block struct {
	id  int
	pos int
	l   int
}
