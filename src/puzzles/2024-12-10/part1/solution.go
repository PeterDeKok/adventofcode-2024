package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/distribute"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"sync/atomic"
)

func Solution(ctx context.Context, il *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int32

	chart, ths, err := parseInput(rd)
	if err != nil {
		return err
	}

	if err := distribute.Group(ctx, il, func(g *errgroup.Group, gctx context.Context) error {
		for i, th := range ths {
			i, th := i, th

			g.Go(func() error {
				if n, err := walker(g, chart, i, th, &sum); err != nil {
					return err
				} else {
					atomic.AddInt32(&sum, n)
				}

				return nil
			})
		}

		return nil
	}); err != nil {
		return err
	}

	if _, err := w.Write([]byte(strconv.Itoa(int(sum)))); err != nil {
		return err
	}

	return nil
}

func walker(g *errgroup.Group, chart *grid.Grid[Cell, *Cell], i int, c *Cell, sum *int32) (n int32, err error) {
	if len(c.Seen) <= i {
		ns := make([]bool, i+1, (i+1)*3)
		copy(ns, c.Seen[:])
		c.Seen = ns
	}

	if c.Weight == 9 {
		if !c.Seen[i] {
			n = 1
		}

		c.Seen[i] = true
		return n, nil
	}

	c.Seen[i] = true
	w := c.Weight + 1

	if right := chart.Get(c.Y(), c.X()+1); right != nil && right.Weight == w {
		if nn, err := walker(g, chart, i, right, sum); err != nil {
			return n, err
		} else {
			n += nn
		}
	}

	if down := chart.Get(c.Y()+1, c.X()); down != nil && down.Weight == w {
		if nn, err := walker(g, chart, i, down, sum); err != nil {
			return n, err
		} else {
			n += nn
		}
	}

	if left := chart.Get(c.Y(), c.X()-1); left != nil && left.Weight == w {
		if nn, err := walker(g, chart, i, left, sum); err != nil {
			return n, err
		} else {
			n += nn
		}
	}

	if up := chart.Get(c.Y()-1, c.X()); up != nil && up.Weight == w {
		if nn, err := walker(g, chart, i, up, sum); err != nil {
			return n, err
		} else {
			n += nn
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return n, err
		}
	}

	return n, nil
}

func parseInput(rd io.Reader) (g *grid.Grid[Cell, *Cell], th []*Cell, err error) {
	trailheads := make([]*Cell, 0, 100)

	g = grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err = g.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			if r == '0' {
				trailheads = append(trailheads, c)
			}

			return c, nil
		})

		if build.DEBUG {
			fmt.Println("")
		}

		if err != nil {
			return g, trailheads, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return g, trailheads, err
		}
	}

	return g, trailheads, nil
}

type Cell struct {
	grid.BaseCell
	Weight int
	Seen   []bool
}

func CreateCell(y, x int, r rune) *Cell {
	return &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),
		Weight:   int(r - '0'),
		Seen:     make([]bool, 277),
	}
}

func (c *Cell) Bytes() []byte {
	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if c.Weight == 0 {
		return '#'
	}

	if len(c.Seen) > 0 {
		return rune('0' + c.Weight)
	}

	return '░'
}
