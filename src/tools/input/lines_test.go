package input

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func TestLineReader_NormalInput(t *testing.T) {
	rd := strings.NewReader("line1\nline2\nline3")
	var results []string

	LineReader(rd)(func(_ int, line string) bool {
		results = append(results, line)
		return true
	})

	assert.EqualSlice(t, []string{"line1", "line2", "line3"}, results)
}

func TestLineReader_EmptyInput(t *testing.T) {
	rd := strings.NewReader("")
	var results []string

	LineReader(rd)(func(_ int, line string) bool {
		results = append(results, line)
		return true
	})

	assert.HasLen[string](t, 0, results)
}

func TestLineReader_PartialLines(t *testing.T) {
	lines := []string{
		// 4096 is the default buffersize (which is used to split prefixes) of the bufio Reader
		// int truncating of the remainder is negated by adding 1 extra.
		// The 1 extra is to ensure the input lines are always above the buffersize,
		// even if the repeating string has no remainder.
		strings.Repeat("line1", (4096/len("line1"))+2),
		strings.Repeat("line2", (4096/len("line2"))+2),
	}
	rd := strings.NewReader(strings.Join(lines, "\n"))
	var results []string

	LineReader(rd)(func(_ int, line string) bool {
		results = append(results, line)
		return true
	})

	assert.EqualSlice(t, lines, results)
}

func TestLineReader_PanicOnUnexpectedError(t *testing.T) {
	assert.ShouldPanic(t, func() {
		LineReader(utils.ErrorReader{})(func(_ int, line string) bool {
			return true
		})
	})
}

func TestLineIsStrSliceReader_NormalInput(t *testing.T) {
	rd := strings.NewReader("field1 field2\nfield3 field4")
	var results [][]string

	LineIsStrSliceReader(rd)(func(_ int, fields []string) bool {
		results = append(results, fields)
		return true
	})

	assert.EqualSlice2D(t, [][]string{{"field1", "field2"}, {"field3", "field4"}}, results)
}

func TestLineIsStrSliceReader_EmptyLines(t *testing.T) {
	rd := strings.NewReader("field1 field2\n\nfield3")
	var results [][]string

	LineIsStrSliceReader(rd)(func(_ int, fields []string) bool {
		results = append(results, fields)
		return true
	})

	assert.EqualSlice2D(t, [][]string{{"field1", "field2"}, {}, {"field3"}}, results)
}

func TestLineIsIntSliceReader_NormalInput(t *testing.T) {
	rd := strings.NewReader("1 2\n3 4")
	var results [][]int

	LineIsIntSliceReader(rd)(func(_ int, numbers []int) bool {
		results = append(results, numbers)
		return true
	})

	assert.EqualSlice2D(t, [][]int{{1, 2}, {3, 4}}, results)
}

func TestLineIsIntSliceReader_InvalidInput(t *testing.T) {
	rd := strings.NewReader("1 a\n3 4")

	assert.ShouldPanic(t, func() {
		LineIsIntSliceReader(rd)(func(_ int, _ []int) bool {
			return true
		})
	})
}

func TestLineIsIntSliceReader_EmptyLines(t *testing.T) {
	rd := strings.NewReader("1 2\n\n3 4")
	var results [][]int

	LineIsIntSliceReader(rd)(func(_ int, numbers []int) bool {
		results = append(results, numbers)
		return true
	})

	assert.EqualSlice2D(t, [][]int{{1, 2}, {}, {3, 4}}, results)
}
