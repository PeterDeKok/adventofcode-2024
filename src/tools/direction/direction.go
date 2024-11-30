package direction

type Dir byte

const (
	Up Dir = 1 << iota
	Right
	Down
	Left
)

const (
	All = H | V
	H   = Right | Left
	V   = Up | Down
)

func (d Dir) TurnLeft() Dir {
	switch d {
	case /*  8 */ Left:
		return Down
	case /*  4 */ Down:
		return Right
	case /*  2 */ Right:
		return Up
	case /*  1 */ Up:
		return Left
	default:
		panic("Turning combined direction")
	}
}

func (d Dir) TurnRight() Dir {
	switch d {
	case /*  8 */ Left:
		return Up
	case /*  4 */ Down:
		return Left
	case /*  2 */ Right:
		return Down
	case /*  1 */ Up:
		return Right
	default:
		panic("Turning combined direction")
	}
}

func (d Dir) Reverse() Dir {
	switch d {
	case /*  8 */ Left:
		return Right
	case /*  4 */ Down:
		return Up
	case /*  2 */ Right:
		return Left
	case /*  1 */ Up:
		return Down
	default:
		panic("Turning combined direction")
	}
}

func (d Dir) Rune() rune {
	switch d {
	case /* 15 */ Left | Down | Right | Up:
		return '┼'
	case /* 14 */ Left | Down | Right:
		return '┬'
	case /* 13 */ Left | Down | Up:
		return '┤'
	case /* 12 */ Left | Down:
		return '┐'
	case /* 11 */ Left | Right | Up:
		return '┴'
	case /* 10 */ Left | Right:
		return '─'
	case /*  9 */ Left | Up:
		return '┘'
	case /*  8 */ Left:
		return '╴'
	case /*  7 */ Down | Right | Up:
		return '├'
	case /*  6 */ Down | Right:
		return '┌'
	case /*  5 */ Down | Up:
		return '│'
	case /*  4 */ Down:
		return '╷'
	case /*  3 */ Right | Up:
		return '└'
	case /*  2 */ Right:
		return '╶'
	case /*  1 */ Up:
		return '╵'
	}

	panic("Direction not recognized")
}

func (d Dir) Y() int {
	dy := 0

	if d&Up > 0 {
		dy -= 1
	}

	if d&Down > 0 {
		dy += 1
	}

	return dy
}

func (d Dir) X() int {
	dx := 0

	if d&Left > 0 {
		dx -= 1
	}

	if d&Right > 0 {
		dx += 1
	}

	return dx
}
