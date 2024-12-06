package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/direction"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"time"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	sum := 0

	chart, guard, err := parseInput(rd)
	if err != nil {
		return err
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	walk(chart, guard)

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	//	err = tools.DistributeMap(ctx, nodes, func(k rune, v []*Cell) error {
	//		return handleFrequency(chart, k, v)
	//	})
	//	if err != nil {
	//		return err
	//	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	for c := range chart.Iter() {
		if c.Seen > 0 {
			sum++
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func walk(chart *grid.Grid[Cell, *Cell], guard *Guard) {
	for steps := 0; ; steps++ {
		if guard.Walk() == nil {
			return
		}

		if build.DEBUG {
			// TODO Instead use the argument logger
			if err := chart.Fprint(os.Stdout); err != nil {
				panic(err)
			}

			time.Sleep(20 * time.Millisecond)
		}
	}
}

func parseInput(rd io.Reader) (chart *grid.Grid[Cell, *Cell], guard *Guard, err error) {
	chart = grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err = chart.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			if c.Seen > 0 && guard != nil {
				return nil, fmt.Errorf("failed to parse line %d: %v", y, err)
			} else if c.Seen > 0 {
				// Note; c.Seen could be multiple directions,
				// that would be a mistake. Not guarding for that here...
				guard = CreateGuard(chart, c, c.Seen)
			}

			return c, nil
		})

		if err != nil {
			return chart, guard, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	if guard == nil {
		return chart, guard, fmt.Errorf("failed to parse lines: missing guard")
	}

	return chart, guard, nil
}

type Cell struct {
	grid.BaseCell
	Obstacle bool
	Seen     direction.Dir
}

func CreateCell(y, x int, r rune) *Cell {
	c := &Cell{
		BaseCell: grid.CreateBaseCell(y, x, '.'),
	}

	if r == '#' {
		c.Obstacle = true
	} else if r == '^' {
		c.Seen |= direction.Up
	}

	return c
}

func (c *Cell) Bytes() []byte {
	if c == nil {
		return []byte("█")
	}

	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if c.Obstacle {
		return '█'
	}

	if c.Seen > 0 && c.Seen < 16 {
		return c.Seen.Rune()
	}

	return '░'
}

type Vec struct {
	*Cell
	Dir direction.Dir
}

func (v *Vec) String() string {
	return fmt.Sprintf("@(%d,%d → %s)", v.Y(), v.X(), string(v.Dir.Rune()))
}

func (v *Vec) Next() (y, x int) {
	return v.Y() + v.Dir.Y(), v.X() + v.Dir.X()
}

func (v *Vec) TurnRight() *Vec {
	d := v.Dir.TurnRight()

	return CreateVec(v.Cell, d)
}

func CreateVec(c *Cell, enteredInDirection direction.Dir) *Vec {
	return &Vec{
		Cell: c,
		Dir:  enteredInDirection,
	}
}

type Guard struct {
	chart   *grid.Grid[Cell, *Cell]
	History []*Vec
}

func (g *Guard) Last() *Vec {
	return g.History[len(g.History)-1]
}

func (g *Guard) Walk() *Vec {
	v := g.Last()
	nv := v

	nc := g.chart.Get(v.Next())
	if nc == nil {
		return nil
	}

	for i := 0; nc.Obstacle; i++ {
		if i > 3 {
			panic("full circle: nowhere to go")
		}

		nv = nv.TurnRight()

		nc = g.chart.Get(nv.Next())
		if nc == nil {
			return nil
		}
	}

	nv = CreateVec(nc, nv.Dir)
	nc.Seen |= nv.Dir
	v.Cell.Seen |= nv.Dir
	g.History = append(g.History, nv)

	return nv
}

func CreateGuard(chart *grid.Grid[Cell, *Cell], c *Cell, d direction.Dir) *Guard {
	history := make([]*Vec, 1, 10000)

	history[0] = CreateVec(c, d)
	return &Guard{
		chart:   chart,
		History: history,
	}
}
