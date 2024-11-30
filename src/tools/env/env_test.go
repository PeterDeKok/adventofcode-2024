package env

import (
	"errors"
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func TestEnvResult_String(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")

	len76 := strings.Repeat("a", 76)
	len152 := len76 + len76
	len152WithNL := "aaaaa\n" + len76 + "aaaaa\n" + len76
	values := []utils.Expected[envResult[string], string]{
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: "AOCTESTENV description"},
			Expected: fmt.Sprintf(
				"%s\n    %s\n%s\n",
				"AOC_TEST_ENV",
				"AOCTESTENV description",
				"",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: len152},
			Expected: fmt.Sprintf(
				"%s\n    %s\n    %s\n%s\n",
				"AOC_TEST_ENV",
				len76, len76,
				"",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: len152WithNL},
			Expected: fmt.Sprintf(
				"%s\n    aaaaa\n    %s\n    aaaaa\n    %s\n%s\n",
				"AOC_TEST_ENV",
				len76, len76,
				"",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: "AOCTESTENV description", err: errors.New("test err")},
			Expected: fmt.Sprintf(
				"%s\n    %s\n%s\n",
				"AOC_TEST_ENV",
				"AOCTESTENV description",
				"\x1b[31m    err: test err\x1b[0m\n",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: "AOCTESTENV description", err: errors.New(len152)},
			Expected: fmt.Sprintf(
				"%s\n    %s\n%s\n%s\n%s\n\n",
				"AOC_TEST_ENV",
				"AOCTESTENV description",
				"\x1b[31m    err: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
				"\x1b[31m    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
				"\x1b[31m    aaaaa\x1b[0m                                                                       ",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: "AOCTESTENV description", err: errors.New(len152WithNL)},
			Expected: fmt.Sprintf(
				"%s\n    %s\n%s\n%s\n%s\n%s\n\n",
				"AOC_TEST_ENV",
				"AOCTESTENV description",
				"\x1b[31m    err: aaaaa\x1b[0m                                                                  ",
				"\x1b[31m    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
				"\x1b[31m    aaaaa\x1b[0m                                                                       ",
				"\x1b[31m    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
			),
		},
		{
			Value: envResult[string]{key: "AOC_TEST_ENV", description: "AOCTESTENV description", err: errors.New(len152WithNL), errExtra: errors.New("extra")},
			Expected: fmt.Sprintf(
				"%s\n    %s\n%s\n%s\n%s\n%s\n%s\n\n",
				"AOC_TEST_ENV",
				"AOCTESTENV description",
				"\x1b[31m    err: aaaaa\x1b[0m                                                                  ",
				"\x1b[31m    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
				"\x1b[31m    aaaaa\x1b[0m                                                                       ",
				"\x1b[31m    aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\x1b[0m",
				"\x1b[31m    : extra\x1b[0m                                                                     ",
			),
		},
	}

	for i, exp := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.TypeOf[fmt.Stringer](t, exp.Value)
			assert.Equal(t, exp.Expected, exp.Value.String())
		})
	}
}
