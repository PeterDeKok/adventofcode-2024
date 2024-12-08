package main

import (
	"context"
	"fmt"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/distribute"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(ctx context.Context, il *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	chart, nodes, err := parseInput(rd)
	if err != nil {
		return err
	}

	err = distribute.Map(ctx, il, nodes, func(ctx context.Context, _ rune, v []*Cell) error {
		return handleFrequency(chart, v)
	})
	if err != nil {
		return err
	}

	for c := range chart.Iter() {
		if c.Antinode {
			sum++
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func handleFrequency(chart *grid.Grid[Cell, *Cell], v []*Cell) error {
	if len(v) < 2 {
		return nil
	}

	for i, c1 := range v[:len(v)-1] {
		for _, c2 := range v[i+1:] {
			dy, dx := c2.Y()-c1.Y(), c2.X()-c1.X()

			an1y, an1x := c1.Y()-dy, c1.X()-dx
			if c := chart.Get(an1y, an1x); c != nil {
				c.Antinode = true
			}

			an2y, an2x := c2.Y()+dy, c2.X()+dx
			if c := chart.Get(an2y, an2x); c != nil {
				c.Antinode = true
			}
		}
	}

	return nil
}

func parseInput(rd io.Reader) (chart *grid.Grid[Cell, *Cell], nodes map[rune][]*Cell, err error) {
	chart = grid.CreateGrid[Cell]()
	nodes = make(map[rune][]*Cell)

	for y, line := range input.LineReader(rd) {
		err = chart.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			if c.Node != nil {
				nn, ok := nodes[r]
				if !ok {
					nn = make([]*Cell, 0)
					nodes[r] = nn
				}

				nn = append(nn, c)
				nodes[r] = nn
			}

			return c, nil
		})

		if err != nil {
			return chart, nodes, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	return chart, nodes, nil
}

type Cell struct {
	grid.BaseCell
	Node     *rune
	Antinode bool
}

func CreateCell(y, x int, r rune) *Cell {
	var rr *rune

	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		rr = &r
	}

	return &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),

		Node: rr,
	}
}
