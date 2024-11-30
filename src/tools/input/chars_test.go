package input

import (
	"bytes"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func TestCharReader_NormalInput(t *testing.T) {
	rd := strings.NewReader("hello")
	var results []rune

	CharReader(rd)(func(_ int, v rune) bool {
		results = append(results, v)
		return true
	})

	assert.Equal(t, "hello", string(results))
}

func TestCharReader_EmptyInput(t *testing.T) {
	rd := strings.NewReader("")
	var results []rune

	CharReader(rd)(func(_ int, v rune) bool {
		results = append(results, v)
		return true
	})

	assert.Equal(t, "", string(results))
}

func TestCharReader_NonUTF8Input(t *testing.T) {
	rd := bytes.NewReader([]byte{0xff, 0xfe, 0xfd})
	var results []rune

	CharReader(rd)(func(_ int, v rune) bool {
		results = append(results, v)
		return true
	})

	assert.EqualSlice(t, []rune{0xff, 0xfe, 0xfd}, results)
}

func TestCharReader_LongInput(t *testing.T) {
	// 1GB of 'a'
	size := 1024 * 1024 * 1024
	rd := strings.NewReader(strings.Repeat("a", size))
	count := 0

	CharReader(rd)(func(_ int, v rune) bool {
		if v != 'a' {
			utils.Fatalf(t, v, 'a')
		}
		count++
		return true
	})

	assert.Equal(t, size, count)
}

func TestCharReader_ErrorHandling(t *testing.T) {
	assert.ShouldPanic(t, func() {
		CharReader(utils.ErrorReader{})(func(_ int, _ rune) bool {
			return true
		})
	})
}

func TestCharReader_MultiByteInput(t *testing.T) {
	rd := strings.NewReader("aあ") // 'あ' is a multi-byte rune
	var results []rune

	CharReader(rd)(func(_ int, v rune) bool {
		results = append(results, v)
		return true
	})

	// Expect only single-byte characters, 'a' and the individual bytes of 'あ'
	assert.EqualSlice(t, []rune{'a', 0xe3, 0x81, 0x82}, results)
}
