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
	var h Opt

	var firstDigit *digit
	mods := make([]Result, 0)

	for i, r := range input.CharReader(rd) {
		if h == nil {
			if build.DEBUG {
				fmt.Printf("  [rune %d]: new mul\n\n", i)
			}
			h = Mul()
		}

		if build.DEBUG {
			fmt.Printf("[rune %d] %v\n", i, r)
		}
		next := h.Next(r)

		if hh, ok := h.(*digit); ok && hh.nth == 0 && hh.done {
			// done 1st
			if build.DEBUG {
				fmt.Printf("  [rune %d | 1st]: %d\n\n", i, hh.Int())
			}
			firstDigit = hh
		} else if hh, ok = h.(*digit); ok && hh.nth == 1 && hh.done {
			if firstDigit == nil {
				panic("First digit nil, should not happen")
			}
			// done 2nd
			result := MulResult(firstDigit, hh)
			if build.DEBUG {
				fmt.Printf("  [rune %d | result]: %s\n\n", i, result)
			}
			mods = append(mods, result)
			firstDigit = nil
		} else if next == nil {
			// Reset
			firstDigit = nil
			h = Mul().Next(r)
			continue
		}

		h = next
	}

	sum := 0

	for _, result := range mods {
		sum += result.Result()
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

type Opt interface {
	Next(r rune) Opt
	Done() bool
}

type mul struct {
	l int
}

func Mul() Opt {
	return &mul{}
}

func (m *mul) Done() bool { return false }
func (m *mul) Next(r rune) Opt {
	switch {
	case m.l == 0 && r == 'm':
		fallthrough
	case m.l == 1 && r == 'u':
		fallthrough
	case m.l == 2 && r == 'l':
		m.l++
		return m
	case m.l == 3 && r == '(':
		m.l++
		return Digit(0)
	default:
		return nil
	}
}

type digit struct {
	nth  int
	d    []rune
	done bool
}

func Digit(nth int) Opt {
	return &digit{
		nth: nth,
		d:   make([]rune, 0, 3),
	}
}

func (n *digit) Done() bool { return n.done }
func (n *digit) Next(r rune) Opt {
	l := len(n.d)

	switch {
	case l < 3 && r >= '0' && r <= '9':
		n.d = append(n.d, r)
		return n
	case l > 0 && n.nth == 0 && r == ',':
		n.done = true
		return Digit(1)
	case l > 0 && n.nth == 1 && r == ')':
		n.done = true
		return nil
	default:
		return nil
	}
}
func (n *digit) Int() int {
	sum := 0

	for _, r := range n.d {
		sum = (sum * 10) + int(r-'0')
	}

	return sum
}

type Result interface {
	Result() int
}
type mulResult struct {
	left  *digit
	right *digit
}

func MulResult(left, right *digit) Result {
	return &mulResult{
		left:  left,
		right: right,
	}
}

func (r *mulResult) Result() int {
	return r.left.Int() * r.right.Int()
}
func (r *mulResult) String() string {
	return fmt.Sprintf("%d * %d = %d", r.left.Int(), r.right.Int(), r.Result())
}