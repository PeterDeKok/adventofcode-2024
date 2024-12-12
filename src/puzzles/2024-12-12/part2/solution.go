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

	areas := make([]*Area, 0, 200)

	for c := range chart.Iter() {
		area := &Area{R: c.R}

		if ok, err := walker(chart, c, area); err != nil {
			return err
		} else if ok {
			areas = append(areas, area)
		}
	}

	for _, area := range areas {
		sum += area.Corners * area.Area
		if build.DEBUG {
			fmt.Printf("%s: a:%2d * c:%2d = %4d -> %6d\n", string(area.R), area.Area, area.Corners, area.Area*area.Corners, sum)
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func walker(chart *grid.Grid[Cell, *Cell], c *Cell, area *Area) (bool, error) {
	if c.Seen {
		return false, nil
	}

	c.Seen = true

	area.Corners += c.Area.Corners
	area.Area += c.Area.Area
	c.Area = area
	area.Area++

	var u, r, d, l bool
	var uc, rc, dc, lc *Cell

	if rc = chart.Get(c.Y(), c.X()+1); rc == nil || rc.R != c.R {
		r = true
	} else if _, err := walker(chart, rc, area); err != nil {
		return false, err
	}

	if dc = chart.Get(c.Y()+1, c.X()); dc == nil || dc.R != c.R {
		d = true
	} else if _, err := walker(chart, dc, area); err != nil {
		return false, err
	}

	if lc = chart.Get(c.Y(), c.X()-1); lc == nil || lc.R != c.R {
		l = true
	} else if _, err := walker(chart, lc, area); err != nil {
		return false, err
	}

	if uc = chart.Get(c.Y()-1, c.X()); uc == nil || uc.R != c.R {
		u = true
	} else if _, err := walker(chart, uc, area); err != nil {
		return false, err
	}

	ur, rd, dl, lu := u && r, r && d, d && l, l && u

	if ur {
		c.Area.Corners++
	}

	if rd {
		c.Area.Corners++
	}

	if dl {
		c.Area.Corners++
	}

	if lu {
		c.Area.Corners++
	}

	// Check if the opposite sides are part of the outside of a corner
	if ur {
		if urc := chart.Get(c.Y()-1, c.X()+1); urc != nil && uc != nil && rc != nil && urc.R == uc.R && urc.R == rc.R {
			urc.Area.Corners++
		}
	}

	if rd {
		if rdc := chart.Get(c.Y()+1, c.X()+1); rdc != nil && rc != nil && dc != nil && rdc.R == rc.R && rdc.R == dc.R {
			rdc.Area.Corners++
		}
	}

	if dl {
		if dlc := chart.Get(c.Y()+1, c.X()-1); dlc != nil && dc != nil && lc != nil && dlc.R == dc.R && dlc.R == lc.R {
			dlc.Area.Corners++
		}
	}

	if lu {
		if luc := chart.Get(c.Y()-1, c.X()-1); luc != nil && lc != nil && uc != nil && luc.R == lc.R && luc.R == uc.R {
			luc.Area.Corners++
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := chart.Fprint(os.Stdout); err != nil {
			return false, err
		}
	}

	return true, nil
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
	R    rune
	Seen bool
	Area *Area
}

func CreateCell(y, x int, r rune) *Cell {
	return &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),
		R:        r,
		Area:     &Area{},
	}
}

func (c *Cell) Bytes() []byte {
	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if !c.Seen {
		return '░'
	}

	return c.R
}

type Area struct {
	Area    int
	Corners int
	R       rune
}
