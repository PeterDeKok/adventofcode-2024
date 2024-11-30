package direction

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
)

func TestDir(t *testing.T) {
	assert.Equal(t, 0b00000001, Up)
	assert.Equal(t, 0b00000010, Right)
	assert.Equal(t, 0b00000100, Down)
	assert.Equal(t, 0b00001000, Left)
}

func TestDir_TurnRight(t *testing.T) {
	assert.Equal(t, Right, Up.TurnRight())
	assert.Equal(t, Down, Right.TurnRight())
	assert.Equal(t, Left, Down.TurnRight())
	assert.Equal(t, Up, Left.TurnRight())

	assert.Equal(t, 0b00000010, Up.TurnRight())
	assert.Equal(t, 0b00000100, Right.TurnRight())
	assert.Equal(t, 0b00001000, Down.TurnRight())
	assert.Equal(t, 0b00000001, Left.TurnRight())
}

func TestDir_Y(t *testing.T) {
	assert.Equal(t, 0, (Left | Down | Right | Up).Y())
	assert.Equal(t, 1, (Left | Down | Right).Y())
	assert.Equal(t, 0, (Left | Down | Up).Y())
	assert.Equal(t, 1, (Left | Down).Y())
	assert.Equal(t, -1, (Left | Right | Up).Y())
	assert.Equal(t, 0, (Left | Right).Y())
	assert.Equal(t, -1, (Left | Up).Y())
	assert.Equal(t, 0, (Left).Y())
	assert.Equal(t, 0, (Down | Right | Up).Y())
	assert.Equal(t, 1, (Down | Right).Y())
	assert.Equal(t, 0, (Down | Up).Y())
	assert.Equal(t, 1, (Down).Y())
	assert.Equal(t, -1, (Right | Up).Y())
	assert.Equal(t, 0, (Right).Y())
	assert.Equal(t, -1, (Up).Y())
}

func TestDir_X(t *testing.T) {
	assert.Equal(t, 0, (Left | Down | Right | Up).X())
	assert.Equal(t, 0, (Left | Down | Right).X())
	assert.Equal(t, -1, (Left | Down | Up).X())
	assert.Equal(t, -1, (Left | Down).X())
	assert.Equal(t, 0, (Left | Right | Up).X())
	assert.Equal(t, 0, (Left | Right).X())
	assert.Equal(t, -1, (Left | Up).X())
	assert.Equal(t, -1, (Left).X())
	assert.Equal(t, 1, (Down | Right | Up).X())
	assert.Equal(t, 1, (Down | Right).X())
	assert.Equal(t, 0, (Down | Up).X())
	assert.Equal(t, 0, (Down).X())
	assert.Equal(t, 1, (Right | Up).X())
	assert.Equal(t, 1, (Right).X())
	assert.Equal(t, 0, (Up).X())
}
