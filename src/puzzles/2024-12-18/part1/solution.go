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
	"strings"
	"time"
)

const MaxUint64 = ^uint64(0)
const DEBUGSPEED = 0

var HEIGHT = 71
var WIDTH = 71
var LIMIT = 1024

func Pre(run string) {
	switch run {
	case "sample-input-1.txt":
		HEIGHT = 7
		WIDTH = 7
		LIMIT = 12
	case "input.txt":
		HEIGHT = 71
		WIDTH = 71
		LIMIT = 1024
	}
}

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	g, err := grid.CreateFixedGrid[Cell, *Cell](HEIGHT, WIDTH, func(y, x int) (*Cell, error) {
		return CreateCell(y, x), nil
	})

	if err != nil {
		return err
	}

	s, e := g.Get(0, 0), g.Get(HEIGHT-1, WIDTH-1)
	s.Start = true
	e.End = true

	if build.DEBUG {
		for i := 0; i < HEIGHT; i++ {
			fmt.Println("")
		}
	}

	if err := parseInput(rd, g, LIMIT); err != nil {
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

func parseInput(rd io.Reader, g *grid.Grid[Cell, *Cell], limit int) (err error) {
	for lineNr, line := range input.LineReader(rd) {
		if lineNr >= limit {
			break
		}

		in := strings.FieldsFunc(line, func(r rune) bool {
			return r == ','
		})

		if len(in) != 2 {
			return fmt.Errorf("failed to parse line %d: expected 2 fields", lineNr)
		}

		y, errY := strconv.Atoi(in[1])
		if len(in) != 2 {
			return fmt.Errorf("failed to parse line %d: y field unexpected: %v", lineNr, errY)
		}

		x, errX := strconv.Atoi(in[0])
		if len(in) != 2 {
			return fmt.Errorf("failed to parse line %d: x field unexpected: %v", lineNr, errX)
		}

		c := g.Get(y, x)
		if c == nil {
			return fmt.Errorf("failed to parse line %d: (x: %d, y: %d) is not a cell", lineNr, x, y)
		}

		c.Safe = false
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return err
		}
	}

	return nil
}

type Cell struct {
	grid.BaseCell
	Safe            bool
	Start           bool
	End             bool
	LowestWeight    uint64
	LowestParent    *Cell
	IsShortestRoute bool
}

func CreateCell(y, x int) *Cell {
	return &Cell{
		BaseCell:     grid.CreateBaseCell(y, x, '.'),
		Safe:         true,
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
	if c == nil || !c.Safe {
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
				{d: spec.dd.TurnLeft(), m: 1},
				{d: spec.dd, m: 1},
				{d: spec.dd.TurnRight(), m: 1},
			} {
				if nc := g.Get(spec.cc.Dir(nd.d)); nc != nil && nc.Safe {
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
