package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	chart, err := parseInput(rd)
	if err != nil {
		return err
	}

	for c := range chart.Iter() {
		if edges, area, err := walker(chart, c); err != nil {
			return err
		} else {
			sum += edges * area
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func walker(chart *grid.Grid[Cell, *Cell], c *Cell) (edges, area int, err error) {
	if c.Seen {
		return 0, 0, nil
	}

	c.Seen = true
	edges = 0
	area = 1

	if right := chart.Get(c.Y(), c.X()+1); right == nil || right.R != c.R {
		c.Edges++
	} else if ed, ar, err := walker(chart, right); err != nil {
		return 0, 0, err
	} else {
		edges += ed
		area += ar
	}

	if down := chart.Get(c.Y()+1, c.X()); down == nil || down.R != c.R {
		c.Edges++
	} else if ed, ar, err := walker(chart, down); err != nil {
		return 0, 0, err
	} else {
		edges += ed
		area += ar
	}

	if left := chart.Get(c.Y(), c.X()-1); left == nil || left.R != c.R {
		c.Edges++
	} else if ed, ar, err := walker(chart, left); err != nil {
		return 0, 0, err
	} else {
		edges += ed
		area += ar
	}

	if up := chart.Get(c.Y()-1, c.X()); up == nil || up.R != c.R {
		c.Edges++
	} else if ed, ar, err := walker(chart, up); err != nil {
		return 0, 0, err
	} else {
		edges += ed
		area += ar
	}

	edges += c.Edges

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return edges, area, err
		}
	}

	return edges, area, nil
}

func parseInput(rd io.Reader) (g *grid.Grid[Cell, *Cell], err error) {
	g = grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err = g.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			return c, nil
		})

		if build.DEBUG {
			fmt.Println("")
		}

		if err != nil {
			return g, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return g, err
		}
	}

	return g, nil
}

type Cell struct {
	grid.BaseCell
	R     rune
	Edges int
	Seen  bool
}

func CreateCell(y, x int, r rune) *Cell {
	return &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),
		R:        r,
	}
}

func (c *Cell) Bytes() []byte {
	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if !c.Seen {
		return '░'
	}

	if c.Edges > 0 {
		return rune('0' + c.Edges)
	}

	return '#'
}
