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
	var sum int
	fs := parseInput(rd)

	dedefrag(fs)

	for i, f := range fs {
		if f != nil {
			sum += i * f.id

			if build.DEBUG {
				fmt.Printf("    %d (%d * %d = %d)\n", sum, i, f.id, i*f.id)
			}
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func dedefrag(fs []*Block) {
	if build.DEBUG {
		printfs(fs)
	}

	for i, j := len(fs)-1, 0; i >= j; i-- {
		f := fs[i]

		if build.DEBUG {
			printfs(fs)
		}

		if f == nil {
			continue
		}

		for ; j < len(fs)-1 && fs[j] != nil; j++ {
			// no-op
		}

		if i < j {
			return
		}

		fs[j] = f
		fs[i] = nil
	}

	if build.DEBUG {
		printfs(fs)
	}
}

func printfs(fs []*Block) {
	str := make([]rune, len(fs))

	for i, f := range fs {
		if f == nil {
			str[i] = '.'
		} else {
			str[i] = []rune(fmt.Sprintf("%d", f.id%10))[0]
		}
	}

	fmt.Println(string(str))
}

func parseInput(rd io.Reader) (fs []*Block) {
	file := true
	id := 0
	fs = make([]*Block, 0, 10000)

	for _, r := range input.CharReader(rd) {
		if build.DEBUG {
			fmt.Printf("%s", string(r))
		}

		l := int(r - '0')

		if file {
			for i := 0; i < l; i++ {
				fs = append(fs, &Block{
					id: id,
				})
			}
			id++
		} else {
			for i := 0; i < l; i++ {
				fs = append(fs, nil)
			}
		}

		file = !file
	}

	if build.DEBUG {
		fmt.Println("")

		for i, f := range fs {
			if f != nil {
				fmt.Printf("%5d > %d\n", i, f.id)
			} else {
				fmt.Printf("%5d\n", i)
			}
		}
	}

	return fs
}

type Block struct {
	id int
}
