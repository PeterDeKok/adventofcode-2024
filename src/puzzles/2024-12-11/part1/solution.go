package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"strings"
)

const blinks = 25

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int
	defer func() {
		if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
			panic(err)
		}
	}()

	stones, err := parseInput(rd)
	if err != nil {
		return err
	}

	if build.DEBUG {
		for _, rr := range stones {
			fmt.Printf("%s ", string(rr))
		}
		fmt.Printf("\n\n")
	}

	// Note;
	// The last change to any stone is potentially skipped if it doesn't result
	// in a new stone.
	queue := make([]*TimeLoop, len(stones), 10*len(stones))

	for i, s := range stones {
		queue[i] = &TimeLoop{
			stone: s,
			i:     blinks,
		}
	}

	for i := 0; i < len(queue); i++ {
		if add, err := walker(queue[i], &queue); err != nil {
			return err
		} else {
			sum += add
			if build.DEBUG {
				fmt.Println(add)
			}
		}
	}

	sum += len(queue)

	return nil
}

func walker(tl *TimeLoop, queue *[]*TimeLoop) (add int, err error) {
	for ; tl.i > 0; tl.i-- {
		before := string(tl.stone)

		if len(tl.stone)%2 == 0 {
			ns := tl.stone[len(tl.stone)/2:]
			tl.stone = tl.stone[:len(tl.stone)/2]

			trim := 0
			for ; trim < len(ns)-1; trim++ {
				if ns[trim] != '0' {
					break
				}
			}
			ns = ns[trim:]

			if len(ns) == 0 {
				panic("Mistakes were made")
			}

			*queue = append(*queue, &TimeLoop{
				stone: ns,
				i:     tl.i - 1,
			})

			if build.DEBUG {
				fmt.Printf("stone split %*s -> %s %s\n", 25-tl.i-1, before, string(tl.stone), string(ns))
			}

			continue
		}

		if tl.i == 1 {
			if build.DEBUG {
				fmt.Printf("skip last i %*s\n", 25-tl.i-1, string(tl.stone))
			}
			// Skip the last step. It won't result in more stones.
			// We do skip the last changes to the slice, it won't resemble the true state!!!
			return 0, nil
		}

		if len(tl.stone) == 1 && tl.stone[0] == '0' {
			// Nope, not cheating, it's part of the problem description!
			// 25 = 19778 - 1 = 19777
			// 24 = 12343 - 1 = 12342
			// 23 = 8268 - 1 = 8267
			// 22 = 5602 - 1 = 5601
			// 21 = 3572 - 1 = 3571
			// 20 = 2377 - 1 = 2376
			// 19 = 1546 - 1 = 1545
			// 18 = 1059 - 1 = 1058
			// 17 = 667 - 1 = 666
			// 16 = 418 - 1 = 417
			// 15 = 328 - 1 = 327
			// 14 = 200 - 1 = 199
			// 13 = 110 - 1 = 109
			// 12 = 81 - 1 = 80
			// 11 = 62 - 1 = 61
			// 10 = 39 - 1 = 38
			// 9 = 20 - 1 = 19
			// 8 = 16 - 1 = 15
			// 7 = 14 - 1 = 13
			// 6 = 7 - 1 = 6
			// 5 = 4 - 1 = 3
			// 4 = 4 - 1 = 3
			// 3 = 2 - 1 = 1
			// 2 = 1 - 1 = 0
			// 1 = 1 - 1 = 0

			lookup := []int{
				25: 19777,
				24: 12342,
				23: 8267,
				22: 5601,
				21: 3571,
				20: 2376,
				19: 1545,
				18: 1058,
				17: 666,
				16: 417,
				15: 327,
				14: 199,
				13: 109,
				12: 80,
				11: 61,
				10: 38,
				9:  19,
				8:  15,
				7:  13,
				6:  6,
				5:  3,
				4:  3,
				3:  1,
				2:  0,
				1:  0,
				0:  0,
			}

			add = lookup[tl.i]
			tl.stone = []rune{'-'}
			tl.i = 0

			return add, nil
		}

		nr := 0
		for _, r := range tl.stone {
			nr = nr*10 + int(r-'0')
		}
		nr *= 2024
		tl.stone = []rune(strconv.Itoa(nr))

		if build.DEBUG {
			fmt.Printf("stone*2024  %*s * 2024 -> %s\n", 25-tl.i-1, before, string(tl.stone))
		}
	}

	return 0, nil
}

func parseInput(rd io.Reader) ([][]rune, error) {
	stones := make([][]rune, 0, 1000)

	brd := bufio.NewReader(rd)
	line, isPrefix, err := brd.ReadLine()

	if err != nil {
		return stones, err
	} else if isPrefix {
		return stones, fmt.Errorf("not entire line read")
	}

	for _, w := range strings.Fields(string(line)) {
		stones = append(stones, []rune(w))
	}

	return stones, nil
}

type TimeLoop struct {
	stone []rune
	i     int
}
