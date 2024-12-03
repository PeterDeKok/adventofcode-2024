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

var mulEnabled bool = true

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var h Opt

	var firstDigit *digit
	mods := make([]Result, 0)

	for i, r := range input.CharReader(rd) {
		if h == nil {
			if build.DEBUG {
				fmt.Printf("  [rune %d]: new mul\n\n", i)
			}
			h = MulOrSwitch()
		}

		if build.DEBUG {
			fmt.Printf("[rune %d] %v\n", i, r)
		}
		next := h.Next(r)

		if hh, ok := h.(*sw); ok && hh.done {
			mulEnabled = hh.enable
		} else if hh, ok := h.(*digit); ok && hh.nth == 0 && hh.done {
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
			h = MulOrSwitch().Next(r)
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

type sw struct {
	enable bool
	l      int
	done   bool
}

func Switch() Opt {
	return &sw{}
}

func (s *sw) Done() bool { return s.done }
func (s *sw) Next(r rune) Opt {
	switch {
	case s.l == 0 && r == 'd':
		fallthrough
	case s.l == 1 && r == 'o':
		s.l++
		return s
	case s.l == 2:
		switch r {
		case '(':
			s.enable = true
			s.l++
			return s
		case 'n':
			s.enable = false
			s.l++
			return s
		}
	case s.l == 3:
		switch {
		case s.enable && r == ')':
			s.done = true
			s.l++
			return nil
		case s.enable == false && r == '\'':
			s.l++
			return s
		}
	case s.enable == false && s.l == 4 && r == 't':
		fallthrough
	case s.enable == false && s.l == 5 && r == '(':
		s.l++
		return s
	case s.enable == false && s.l == 6 && r == ')':
		s.done = true
		s.l++
		return nil
	default:
		return nil
	}

	return nil
}

type or struct{}

func MulOrSwitch() Opt {
	return &or{}
}

func (or *or) Done() bool { return false }
func (or *or) Next(r rune) Opt {
	switch {
	case r == 'm' && mulEnabled:
		return Mul().Next(r)
	case r == 'd':
		return Switch().Next(r)
	default:
		return nil
	}
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
