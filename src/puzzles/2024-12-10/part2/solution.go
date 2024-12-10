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
				return walker(g, chart, i, th, &sum)
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

func walker(g *errgroup.Group, chart *grid.Grid[Cell, *Cell], i int, c *Cell, sum *int32) error {
	if c.Weight == 9 {
		atomic.AddInt32(sum, 1)

		return nil
	}

	right := chart.Get(c.Y(), c.X()+1)
	if right != nil && right.Weight == c.Weight+1 {
		if err := walker(g, chart, i, right, sum); err != nil {
			return err
		}
	}

	down := chart.Get(c.Y()+1, c.X())
	if down != nil && down.Weight == c.Weight+1 {
		if err := walker(g, chart, i, down, sum); err != nil {
			return err
		}
	}

	left := chart.Get(c.Y(), c.X()-1)
	if left != nil && left.Weight == c.Weight+1 {
		if err := walker(g, chart, i, left, sum); err != nil {
			return err
		}
	}

	up := chart.Get(c.Y()-1, c.X())
	if up != nil && up.Weight == c.Weight+1 {
		if err := walker(g, chart, i, up, sum); err != nil {
			return err
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	return nil
}

func parseInput(rd io.Reader) (g *grid.Grid[Cell, *Cell], th []*Cell, err error) {
	trailheads := make([]*Cell, 0, 100)

	g = grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err = g.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			if c.Weight == 0 {
				trailheads = append(trailheads, c)
			}

			return c, nil
		})

		if build.DEBUG {
			fmt.Println(line)
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
}

func CreateCell(y, x int, r rune) *Cell {
	c := &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),
		Weight:   int(r - '0'),
	}

	return c
}

func (c *Cell) Bytes() []byte {
	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if c.Weight == 0 {
		return '#'
	}

	return rune('0' + c.Weight)
}
