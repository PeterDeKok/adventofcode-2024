package main

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/direction"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"time"
)

const MaxUint64 = ^uint64(0)
const DEBUGSPEED = 0

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	g, s, e, err := parseInput(rd)

	if err != nil {
		return err
	}

	for range s.WalkSomewhere(g, 0, direction.Right) {
		if build.DEBUG {
			// TODO Instead use the argument logger
			if err := g.Fprint(os.Stdout); err != nil {
				panic(err)
			}

			time.Sleep(DEBUGSPEED * time.Millisecond)
		}
	}

	for c := e; c.LowestParent != nil; c = c.LowestParent {
		c.IsShortestRoute = true
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(strconv.FormatUint(e.LowestWeight, 10))); err != nil {
		return err
	}

	return nil
}

func parseInput(rd io.Reader) (g *grid.Grid[Cell, *Cell], s, e *Cell, err error) {
	g = grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err = g.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			if c == nil {
				return c, nil
			} else if c.Start {
				s = c
			} else if c.End {
				e = c
			}

			return c, nil
		})

		if build.DEBUG {
			fmt.Println("")
		}

		if err != nil {
			return g, nil, nil, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return g, s, e, err
		}
	}

	if s == nil || e == nil {
		return g, s, e, fmt.Errorf("failed to parse input: missing start or end")
	}

	return g, s, e, nil
}

type Cell struct {
	grid.BaseCell
	Start           bool
	End             bool
	LowestWeight    uint64
	LowestParent    *Cell
	IsShortestRoute bool
}

func CreateCell(y, x int, r rune) *Cell {
	if r == '#' {
		return nil
	}

	return &Cell{
		BaseCell:     grid.CreateBaseCell(y, x, r),
		Start:        r == 'S',
		End:          r == 'E',
		LowestWeight: MaxUint64,
	}
}

func (c *Cell) Bytes() []byte {
	if c == nil {
		return []byte("█")
	}

	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if c == nil {
		return '█'
	}

	if c.Start {
		return '›'
	}

	if c.End {
		return '🮖'
	}

	if c.IsShortestRoute {
		return '¤'
	}

	if c.LowestParent != nil {
		return '·'
	}

	return '░'
}

type OpenSpec struct {
	cc    *Cell
	dd    direction.Dir
	carry uint64
}

type TryDir struct {
	d direction.Dir
	m uint64
}

func (c *Cell) WalkSomewhere(g *grid.Grid[Cell, *Cell], carry uint64, dir direction.Dir) iter.Seq2[*Cell, direction.Dir] {
	return func(yield func(*Cell, direction.Dir) bool) {
		open := []OpenSpec{{cc: c, dd: dir, carry: carry}}

		for i := 0; i < len(open); i++ {
			spec := open[i]
			for _, nd := range []TryDir{
				{d: spec.dd.TurnLeft(), m: 1001},
				{d: spec.dd, m: 1},
				{d: spec.dd.TurnRight(), m: 1001},
			} {
				if nc := g.Get(spec.cc.Dir(nd.d)); nc != nil {
					weight := spec.carry + nd.m

					if weight < nc.LowestWeight {
						nc.LowestWeight = weight
						nc.LowestParent = spec.cc

						if !nc.End {
							open = append(open, OpenSpec{cc: nc, dd: nd.d, carry: weight})
						}

						if !yield(nc, nd.d) {
							return
						}
					}
				}
			}
		}
	}
}
