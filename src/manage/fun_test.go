package manage

import (
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func TestHexToAnsi(t *testing.T) {
	t.SkipNow()
	fmt.Printf("\\033%s\n", HexToAnsi("#886655")[1:])
}

func TestParseFunLines(t *testing.T) {
	t.Skip("The html and expected strings do not belong together. The advent calendar on the website is not static.")
	str, err := ParseFunLines(strings.NewReader(funLoggedIn))
	assert.NoErr(t, err)
	assert.Equal(t, funLoggedInExpected, str)

	fmt.Println(funLoggedInExpected)
}

func TestParseColor(t *testing.T) {
	values := []utils.Expected[string, struct {
		A, B string
		Ok   bool
	}]{
		{
			Value: ".calendar .calendar-color-6y { color:#ffff66; text-shadow:0 0 5px #ffff66; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "6y", B: "\x1b[38;2;255;255;102m", Ok: true},
		},
		{
			Value: ".calendar .calendar-color-1w1 { color: #ffff66; text-shadow:0 0 5px #ffff66; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "1w1", B: "\x1b[38;2;255;255;102m", Ok: true},
		},
		{
			Value: ".calendar .calendar-color-9n { text-shadow:0 0 3px #456efe,0 0 5px #456efe,0 0 10px #456efe,0 0 15px #456efe; color: #66ff66; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "9n", B: "\x1b[38;2;102;255;102m", Ok: true},
		},
		{
			Value: ".calendar .calendar-color-1w2 { text-shadow:0 0 3px #456efe,0 0 5px #456efe,0 0 10px #456efe,0 0 15px #456efe; color: #66ff66; } .calendar .calendar-color-0r::before { content:\"*\"; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "1w2", B: "\x1b[38;2;102;255;102m", Ok: true},
		},
		{
			Value: ".calendar .calendar-color-w { color: #ccc; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "w", B: "\x1b[38;2;204;204;204m", Ok: true},
		},
		{
			Value: ".calendar i { font-style:normal; display:inline-block; width:.6em; line-height:.6em; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "", B: "", Ok: false},
		},
		{
			// This value is an edge case, but it should never resolve anyway, so its fine to leave it!
			Value: ".calendar .calendar-color-0r::before { content:\"*\"; position:absolute; color:#ffff66; transform:translate(-.5px,-.6em) scale(.5); text-shadow: 0 0 25px #ffff66, 0 0 20px #ffff66, 0 0 15px #ffff66, 0 0 10px #ffff66, 0 0 5px #ffff66; }",
			Expected: struct {
				A, B string
				Ok   bool
			}{A: "0r::before", B: "\x1b[38;2;255;255;102m", Ok: true},
		},
	}

	for i, exp := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			a, b, ok := parseColor(exp.Value)
			assert.Equal(t, exp.Expected.A, a)
			assert.Equal(t, exp.Expected.B, b)
			assert.Equal(t, exp.Expected.Ok, ok)
		})
	}
}

func TestParseCalendarFirst(t *testing.T) {
	values := []utils.Expected[string, string]{
		{
			Value:    `<pre class="calendar">          .-----.          .------------------.       `,
			Expected: "\033[38;2;204;204;204m          .-----.          .------------------.\033[m",
		},
	}

	for i, exp := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, exp.Expected, parseCalendarFirst(exp.Value))
		})
	}
}

func TestParseCalendar(t *testing.T) {
	values := []utils.Expected[string, string]{
		{
			Value:    `<a aria-label="Day 1, two stars" href="/2024/day/1" class="calendar-day1 calendar-verycomplete">       <span class="calendar-color-w">.--'</span><span class="calendar-color-3s">~</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3s">~</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">.-'</span> <span class="calendar-color-6y">*</span>       <span class="calendar-color-8n">\</span>  <span class="calendar-color-8n">/</span>     <span class="calendar-color-w">'-.</span>  <span class="calendar-day"> 1</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>`,
			Expected: "\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m.--'\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m.-'\033[m\033[38;2;102;102;102m \033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m       \033[m\033[38;2;136;102;85m\\\033[m\033[38;2;102;102;102m  \033[m\033[38;2;136;102;85m/\033[m\033[38;2;102;102;102m     \033[m\033[38;2;204;204;204m'-.\033[m\033[38;2;102;102;102m\033[m",
		},
		{
			Value:    `<a aria-label="Day 13" href="/2024/day/13" class="calendar-day13">       |        |        |        |    |        |  <span class="calendar-day">13</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>`,
			Expected: "\033[38;2;102;102;102m       |        |        |        |    |        |\033[m",
		},
		{
			Value:    `<a aria-label="Day 16, one star" href="/2024/day/16" class="calendar-day16 calendar-complete">       | '.~  '.|        | : :::::|    |<i>─</i><i>─</i><i>┤</i>AoC<i>├</i>o|  <span class="calendar-day">16</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>`,
			Expected: "\033[38;2;102;102;102m       | '.~  '.|        | : :::::|    |──┤AoC├o|\033[m",
		},
	}

	colors := map[string]string{
		"w":  "\033[38;2;204;204;204m",
		"3s": "\033[38;2;227;181;133m",
		"6y": "\033[38;2;255;255;102m",
		"8n": "\033[38;2;136;102;85m",
	}

	for i, exp := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result, ok := parseCalendar(exp.Value, colors)
			assert.Equal(t, true, ok)
			assert.Equal(t, exp.Expected, result)
		})
	}
}

func TestParseCalendarRaw(t *testing.T) {
	values := []utils.Expected[string, string]{
		{
			Value:    `'----------------------'   '------------------'       `,
			Expected: "\033[38;2;204;204;204m'----------------------'   '------------------'\033[m",
		},
	}

	for i, exp := range values {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, exp.Expected, parseCalendarRaw(exp.Value))
		})
	}
}

const funLoggedInExpected =
// Mismatch between HTML and parsed
// Parsed is from different HTML (they differ in details)
"\x1b[38;2;204;204;204m          .-----.          .------------------.\033[m\n" +
	"\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m.--'\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m.-'\033[m\033[38;2;102;102;102m \033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m       \033[m\033[38;2;136;102;85m\\\033[m\033[38;2;102;102;102m  \033[m\033[38;2;136;102;85m/\033[m\033[38;2;102;102;102m     \033[m\033[38;2;204;204;204m'-.\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m.--'\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m  \033[m\033[38;2;0;204;0m,\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\033[38;2;0;153;0m>\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;153;0mo\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m<\033[m\033[38;2;102;102;102m   \033[m\033[38;2;136;102;85m\\_\\_\\|_/__/\033[m\033[38;2;102;102;102m   \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m.---'\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m:\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m'\033[m\033[38;2;102;102;102m\033[m\033[38;2;85;85;187m(~)\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;204;0m,\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;153;0m>\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;0;0m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m>\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;102;255mO\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m<\033[m\033[38;2;102;102;102m \033[m\033[38;2;255;0;0mo\033[m\033[38;2;102;102;102m\033[m\033[38;2;136;102;85m-_/\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m.\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m()\033[m\033[38;2;102;102;102m\033[m\033[38;2;136;102;85m__------\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;72;136;19m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;94;171;180m..\033[m\033[38;2;102;102;102m\033[m\033[38;2;77;139;3m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m'.\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m\"\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m'\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m>\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;102;255mO\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m>\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;153;0mo\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m<\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;0;0m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;153;0m<\033[m\033[38;2;102;102;102m \033[m\033[38;2;136;102;85m\\____\033[m\033[38;2;102;102;102m       \033[m\033[38;2;0;204;0m.'\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;66;115;34m_\033[m\033[38;2;102;102;102m\033[m\033[38;2;94;171;180m.~.\033[m\033[38;2;102;102;102m\033[m\033[38;2;72;136;19m_@\033[m\033[38;2;102;102;102m\033[m\033[38;2;227;181;133m'..\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;227;181;133m~\033[m\033[38;2;102;102;102m \033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m \033[m\033[38;2;170;170;170m_|\033[m\033[38;2;102;102;102m \033[m\033[38;2;170;170;170m|_\033[m\033[38;2;102;102;102m   \033[m\033[38;2;204;204;204m..\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m\\_\033[m\033[38;2;102;102;102m\033[m\033[38;2;136;102;85m\\_\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m..\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m| ||| #@##'''...|        |...     .'  '.'''../..|  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;66;115;34m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;255m~~~\033[m\033[38;2;102;102;102m\033[m\033[38;2;77;139;3m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;127;189;57m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;77;139;3m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;127;189;57m@\033[m\033[38;2;102;102;102m\033[m\033[38;2;66;115;34m#@\033[m\033[38;2;102;102;102m\033[m\033[38;2;72;136;19m#\033[m\033[38;2;102;102;102m\033[m\033[38;2;77;139;3m@\033[m\033[38;2;102;102;102m   \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;165;168;175m/\\\033[m\033[38;2;102;102;102m \033[m\033[38;2;162;81;81m''.\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m   \033[m\033[38;2;204;204;255m-\033[m\033[38;2;102;102;102m\033[m\033[38;2;212;221;228m/\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m:\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;94;171;180m~~.\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m.--.\033[m\033[38;2;102;102;102m _____  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m \033[m\033[38;2;165;168;175m/\033[m\033[38;2;102;102;102m\033[m\033[38;2;223;35;8m~\033[m\033[38;2;102;102;102m\033[m\033[38;2;165;168;175m\\\033[m\033[38;2;102;102;102m \033[m\033[38;2;162;81;81m'.\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m \033[m\033[38;2;204;204;255m-\033[m\033[38;2;102;102;102m \033[m\033[38;2;212;221;228m/\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m.'\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m'---'\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m|[][]_\\-\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;223;35;8m~\033[m\033[38;2;102;102;102m\033[m\033[38;2;165;168;175m/\033[m\033[38;2;102;102;102m \033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m \033[m\033[38;2;165;168;175m\\\033[m\033[38;2;102;102;102m \033[m\033[38;2;162;81;81m:\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;255m..'\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m------- \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m   \033[m\033[38;2;165;168;175m/\\\033[m\033[38;2;102;102;102m \033[m\033[38;2;162;81;81m.'\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;255m'''\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;200;255m~~~~~\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;255m.......\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;102m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;165;168;175m/\\\033[m\033[38;2;102;102;102m \033[m\033[38;2;162;81;81m..'\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;0;181;237m.\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m.\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m-\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m-\033[m\033[38;2;102;102;102m  \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m\033[m\033[38;2;162;81;81m'''\033[m\033[38;2;102;102;102m\033[m\033[38;2;51;51;51m::\033[m\033[38;2;102;102;102m\033[m\033[38;2;255;255;102m:\033[m\033[38;2;102;102;102m\033[m\033[38;2;51;51;51m::\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\033[38;2;255;255;255m.\033[m\033[38;2;102;102;102m     \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        |    |        |  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        |    |        |  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        |    |        |  \033[m\n" +
	"\033[38;2;102;102;102m       | '.~  '.|        | :.:::::|    |──┤AoC├o|  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        |    |        |  \033[m\n" +
	"\033[38;2;102;102;102m       \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m'..'\033[m\033[38;2;102;102;102m \033[m\033[38;2;0;204;0m.'\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m        \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m \033[m\033[38;2;51;51;51m.\033[m\033[38;2;102;102;102m '\033[m\033[38;2;69;110;254mo\033[m\033[38;2;102;102;102m \033[m\033[38;2;51;51;51m..\033[m\033[38;2;102;102;102m\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m    \033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m┘\033[m\033[38;2;255;255;102m*\033[m\033[38;2;102;102;102m┤yrs├─\033[m\033[38;2;204;204;204m|\033[m\033[38;2;102;102;102m  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        |    |        |  \033[m\n" +
	"\033[38;2;102;102;102m       |        |        |        '.  .'        |  \033[m\n" +
	"\033[38;2;102;102;102m.------'        '------. |          ''          |  \033[m\n" +
	"\033[38;2;102;102;102m|                      | |                      |  \033[m\n" +
	"\033[38;2;102;102;102m|                      | |                      |  \033[m\n" +
	"\033[38;2;102;102;102m|                      | |                      |  \033[m\n" +
	"\033[38;2;102;102;102m|                      | '-.                  .-'  \033[m\n" +
	"\033[38;2;204;204;204m'----------------------'   '------------------'  \033[m"

const funLoggedIn = `<!DOCTYPE html>
<html lang="en-us">
<head>
<meta charset="utf-8"/>
<title>Advent of Code 2024</title>
<link rel="stylesheet" type="text/css" href="/static/style.css?31"/>
<link rel="stylesheet alternate" type="text/css" href="/static/highcontrast.css?1" title="High Contrast"/>
<link rel="shortcut icon" href="/favicon.png"/>
<script>window.addEventListener('click', function(e,s,r){if(e.target.nodeName==='CODE'&&e.detail===3){s=window.getSelection();s.removeAllRanges();r=document.createRange();r.selectNodeContents(e.target);s.addRange(r);}});</script>
</head><!--




Oh, hello!  Funny seeing you here.

I appreciate your enthusiasm, but you aren't going to find much down here.
There certainly aren't clues to any of the puzzles.  The best surprises don't
even appear in the source until you unlock them for real.

Please be careful with automated requests; I'm not a massive company, and I can
only take so much traffic.  Please be considerate so that everyone gets to play.

If you're curious about how Advent of Code works, it's running on some custom
Perl code. Other than a few integrations (auth, analytics, social media), I
built the whole thing myself, including the design, animations, prose, and all
of the puzzles.

The puzzles are most of the work; preparing a new calendar and a new set of
puzzles each year takes all of my free time for 4-5 months. A lot of effort
went into building this thing - I hope you're enjoying playing it as much as I
enjoyed making it for you!

If you'd like to hang out, I'm @was.tl on Bluesky, @ericwastl@hachyderm.io on
Mastodon, and @ericwastl on Twitter.

- Eric Wastl


















































-->
<body>
<header><div><h1 class="title-global"><a href="/">Advent of Code</a></h1><nav><ul><li><a href="/2024/about">[About]</a></li><li><a href="/2024/events">[Events]</a></li><li><a href="https://cottonbureau.com/people/advent-of-code" target="_blank">[Shop]</a></li><li><a href="/2024/settings">[Settings]</a></li><li><a href="/2024/auth/logout">[Log Out]</a></li></ul></nav><div class="user">Peter de Kok <a href="/2024/support" class="supporter-badge" title="Advent of Code Supporter">(AoC++)</a> <span class="star-count">26*</span></div></div><div><h1 class="title-event">&nbsp;&nbsp;<span class="title-event-wrap">{year=&gt;</span><a href="/2024">2024</a><span class="title-event-wrap">}</span></h1><nav><ul><li><a href="/2024">[Calendar]</a></li><li><a href="/2024/support">[AoC++]</a></li><li><a href="/2024/sponsors">[Sponsors]</a></li><li><a href="/2024/leaderboard">[Leaderboard]</a></li><li><a href="/2024/stats">[Stats]</a></li></ul></nav></div></header>

<div id="sidebar">
<div id="sponsor"><div class="quiet">Our <a href="/2024/sponsors">sponsors</a> help make Advent of Code possible:</div><div class="sponsor"><a href="/2024/sponsors/redirect?url=https%3A%2F%2Faoc%2Einfi%2Enl%2F%3Fmtm%5Fcampaign%3Daoc2024%26mtm%5Fsource%3Daoc" target="_blank" onclick="if(ga)ga('send','event','sponsor','sidebar',this.href);" rel="noopener">Infi</a> - Er is slecht weer op komst en het is een lange tocht vanaf de Noordpool... Help jij de kerstman veilig door het luchtruim te navigeren?</div></div>
</div><!--/sidebar-->

<main>
<style>
.calendar .calendar-color-6y { color:#ffff66; text-shadow:0 0 5px #ffff66; }
.calendar .calendar-color-6b { color:#009900; }
.calendar .calendar-color-3y { color:#ffff66; text-shadow:0 0 5px #ffff66, 0 0 10px #ffff66; }
.calendar .calendar-color-8n { color:#886655; }
.calendar i { font-style:normal; display:inline-block; width:.6em; line-height:.6em; }
.calendar .calendar-color-3m { color:#d4dde4; }
.calendar .calendar-color-2g2 { color:#7fbd39; }
.calendar .calendar-color-7y { color:#ffff66; text-shadow:0 0 5px #ffff66; }
.calendar .calendar-color-0l { color:#ccccff; }
.calendar .calendar-color-2w { color:#ffffff; }
.calendar .calendar-color-w { color: #ccc; }
.calendar .calendar-color-6t { color:#aaaaaa; }
.calendar .calendar-color-1w3 { color:#00a2db; }
.calendar .calendar-color-6u { color: #0066ff; text-shadow: 0 0 5px #0066ff; }
.calendar .calendar-color-3g { color:#00cc00; }
.calendar .calendar-color-8i { color:#ff0000; text-shadow:0 0  5px #ff0000, 0 0 10px #ff0000, 0 0 15px #ff0000; }
.calendar .calendar-color-0r { color:#ff0000; position:relative; } .calendar .calendar-color-0r::before { content:"*"; position:absolute; color:#ffff66; transform:translate(-.5px,-.6em) scale(.5); text-shadow: 0 0 25px #ffff66, 0 0 20px #ffff66, 0 0 15px #ffff66, 0 0 10px #ffff66, 0 0 5px #ffff66; }
.calendar .calendar-color-3w { color:#ffffff; }
.calendar .calendar-color-3v { color:#df2308; text-shadow:0 0 5px #df2308, 0 0 10px #df2308; }
.calendar .calendar-color-9n { text-shadow:0 0 3px #456efe,0 0 5px #456efe,0 0 10px #456efe,0 0 15px #456efe; color:#456efe; }
.calendar .calendar-color-2g0 { color:#488813; }
.calendar .calendar-color-3a { color:#a5a8af; }
.calendar .calendar-color-3s { color:#e3b585; }
.calendar .calendar-color-1w1 { color:#00c8ff; }
.calendar .calendar-color-6r { color: #ff0000; text-shadow: 0 0 5px #ff0000; }
.calendar .calendar-color-8w { color:#cccccc; }
.calendar .calendar-color-8e { color:#cccccc; }
.calendar .calendar-color-2u { color:#5eabb4; }
.calendar .calendar-color-3i { color:#a25151; }
.calendar .calendar-color-1w2 { color:#00b5ed; }
.calendar .calendar-color-3l { color:#ccccff; }
.calendar .calendar-color-1s { color:#ffffff; }
.calendar .calendar-color-6o { color: #ff9900; text-shadow: 0 0 5px #ff9900; }
.calendar .calendar-color-2g3 { color:#427322; }
.calendar .calendar-color-0w { color:#ffffff; }
.calendar .calendar-color-6d { color:#333333; }
.calendar .calendar-color-3b { color:#5555bb; }
.calendar .calendar-color-2g1 { color:#4d8b03; }
</style>
<pre class="calendar">          .-----.          .------------------.
<a aria-label="Day 1, two stars" href="/2024/day/1" class="calendar-day1 calendar-verycomplete">       <span class="calendar-color-w">.--'</span><span class="calendar-color-3s">~</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3s">~</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">.-'</span> <span class="calendar-color-6y">*</span>       <span class="calendar-color-8n">\</span>  <span class="calendar-color-8n">/</span>     <span class="calendar-color-w">'-.</span>  <span class="calendar-day"> 1</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 2, two stars" href="/2024/day/2" class="calendar-day2 calendar-verycomplete">    <span class="calendar-color-w">.--'</span><span class="calendar-color-3s">~</span>  <span class="calendar-color-3g">,</span><span class="calendar-color-3y">*</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span>  <span class="calendar-color-6b">&gt;</span><span class="calendar-color-6o">o</span><span class="calendar-color-6b">&lt;</span>   <span class="calendar-color-8n">\_\_\|_/__/</span>   <span class="calendar-color-w">|</span>  <span class="calendar-day"> 2</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 3, two stars" href="/2024/day/3" class="calendar-day3 calendar-verycomplete"><span class="calendar-color-w">.---'</span><span class="calendar-color-3s">:</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3g">'</span><span class="calendar-color-3b">(~)</span><span class="calendar-color-3g">,</span> <span class="calendar-color-3s">~</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span> <span class="calendar-color-6b">&gt;</span><span class="calendar-color-6r">@</span><span class="calendar-color-6b">&gt;</span><span class="calendar-color-6u">O</span><span class="calendar-color-6b">&lt;</span> <span class="calendar-color-8i">o</span><span class="calendar-color-8n">-_/</span><span class="calendar-color-8e">.</span><span class="calendar-color-8w">()</span><span class="calendar-color-8n">__------</span><span class="calendar-color-w">|</span>  <span class="calendar-day"> 3</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 4, two stars" href="/2024/day/4" class="calendar-day4 calendar-verycomplete"><span class="calendar-color-w">|</span><span class="calendar-color-2g0">@</span><span class="calendar-color-2u">..</span><span class="calendar-color-2g1">@</span><span class="calendar-color-3s">'.</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3g">&quot;</span> <span class="calendar-color-3g">'</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-6b">&gt;</span><span class="calendar-color-6u">O</span><span class="calendar-color-6b">&gt;</span><span class="calendar-color-6o">o</span><span class="calendar-color-6b">&lt;</span><span class="calendar-color-6r">@</span><span class="calendar-color-6b">&lt;</span> <span class="calendar-color-8n">\____</span>       <span class="calendar-color-3g">.'</span><span class="calendar-color-w">|</span>  <span class="calendar-day"> 4</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 5, two stars" href="/2024/day/5" class="calendar-day5 calendar-verycomplete"><span class="calendar-color-w">|</span><span class="calendar-color-2g3">_</span><span class="calendar-color-2u">.~.</span><span class="calendar-color-2g0">_@</span><span class="calendar-color-3s">'..</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3s">~</span> <span class="calendar-color-3y">*</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span> <span class="calendar-color-6t">_|</span> <span class="calendar-color-6t">|_</span>    <span class="calendar-color-w">..</span><span class="calendar-color-8w">\_</span><span class="calendar-color-8n">\_</span> <span class="calendar-color-3g">..'</span><span class="calendar-color-3y">*</span> <span class="calendar-color-w">|</span>  <span class="calendar-day"> 5</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 6, one star" href="/2024/day/6" class="calendar-day6 calendar-complete">| ||| @@ #'''...|        |...     .'  '.'''../..|  <span class="calendar-day"> 6</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 7, two stars" href="/2024/day/7" class="calendar-day7 calendar-verycomplete"><span class="calendar-color-w">|</span><span class="calendar-color-2g3">@</span><span class="calendar-color-2w">~~~</span><span class="calendar-color-2g2">@#</span><span class="calendar-color-2g3">@</span> <span class="calendar-color-2g0">#</span><span class="calendar-color-2g1">@</span>  <span class="calendar-color-2g2">#</span>  <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-3a">/\</span> <span class="calendar-color-3i">''.</span>  <span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span>   <span class="calendar-color-3l">-</span><span class="calendar-color-3m">/</span>  <span class="calendar-color-3w">:</span><span class="calendar-color-w">|</span>  <span class="calendar-day"> 7</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 8, two stars" href="/2024/day/8" class="calendar-day8 calendar-verycomplete"><span class="calendar-color-w">|</span><span class="calendar-color-2u">~~.</span><span class="calendar-color-w">.--.</span> _____  <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-3y">*</span> <span class="calendar-color-3a">/</span><span class="calendar-color-3v">~</span><span class="calendar-color-3a">\</span> <span class="calendar-color-3i">'.</span><span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span> <span class="calendar-color-3l">-</span> <span class="calendar-color-3m">/</span>  <span class="calendar-color-3w">.'</span><span class="calendar-color-w">|</span>  <span class="calendar-day"> 8</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 9, two stars" href="/2024/day/9" class="calendar-day9 calendar-verycomplete"><span class="calendar-color-w">'---'</span>  <span class="calendar-color-w">|</span>|[][]_\-<span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-3v">~</span><span class="calendar-color-3a">/</span> <span class="calendar-color-3y">*</span> <span class="calendar-color-3a">\</span> <span class="calendar-color-3i">:</span><span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span>  <span class="calendar-color-3y">*</span><span class="calendar-color-3w">..'</span>  <span class="calendar-color-w">|</span>  <span class="calendar-day"> 9</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 10, two stars" href="/2024/day/10" class="calendar-day10 calendar-verycomplete">       <span class="calendar-color-w">|</span>------- <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span>   <span class="calendar-color-3a">/\</span> <span class="calendar-color-3i">.'</span><span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span><span class="calendar-color-3w">'''</span><span class="calendar-color-1w1">~~~~~</span><span class="calendar-color-w">|</span>  <span class="calendar-day">10</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 11, two stars" href="/2024/day/11" class="calendar-day11 calendar-verycomplete">       <span class="calendar-color-w">|</span><span class="calendar-color-0l">.......</span><span class="calendar-color-0r">|</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-3a">/\</span> <span class="calendar-color-3i">..'</span>  <span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span>   <span class="calendar-color-1s">.</span> <span class="calendar-color-1w2">.~</span> <span class="calendar-color-w">|</span>  <span class="calendar-day">11</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 12, two stars" href="/2024/day/12" class="calendar-day12 calendar-verycomplete">       <span class="calendar-color-w">|</span>  <span class="calendar-color-0w">-</span>  <span class="calendar-color-0w">-</span>  <span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span><span class="calendar-color-3i">'''</span><span class="calendar-color-6d">::</span><span class="calendar-color-6y">:</span><span class="calendar-color-6d">::</span><span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span>  <span class="calendar-color-1s">.</span>    <span class="calendar-color-1w3">.</span><span class="calendar-color-w">|</span>  <span class="calendar-day">12</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 13" href="/2024/day/13" class="calendar-day13">       |        |        |        |    |        |  <span class="calendar-day">13</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 14" href="/2024/day/14" class="calendar-day14">       |        |        |        |    |        |  <span class="calendar-day">14</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 15" href="/2024/day/15" class="calendar-day15">       |        |        |        |    |        |  <span class="calendar-day">15</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 16, one star" href="/2024/day/16" class="calendar-day16 calendar-complete">       | '.~  '.|        | : :::::|    |<i>─</i><i>─</i><i>┤</i>AoC<i>├</i>o|  <span class="calendar-day">16</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 17" href="/2024/day/17" class="calendar-day17">       |        |        |        |    |        |  <span class="calendar-day">17</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 18, two stars" href="/2024/day/18" class="calendar-day18 calendar-verycomplete">       <span class="calendar-color-w">|</span> <span class="calendar-color-3g">'..'</span> <span class="calendar-color-3g">.'</span><span class="calendar-color-w">|</span>        <span class="calendar-color-w">|</span>   '<span class="calendar-color-9n">o</span>   <span class="calendar-color-w">|</span>    <span class="calendar-color-w">|</span><i>┘</i><span class="calendar-color-7y">*</span><i>┤</i>yrs<i>├</i><i>─</i><span class="calendar-color-w">|</span>  <span class="calendar-day">18</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 19" href="/2024/day/19" class="calendar-day19">       |        |        |        |    |        |  <span class="calendar-day">19</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 20" href="/2024/day/20" class="calendar-day20">       |        |        |        '.  .'        |  <span class="calendar-day">20</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 21" href="/2024/day/21" class="calendar-day21">.------'        '------. |          ''          |  <span class="calendar-day">21</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 22" href="/2024/day/22" class="calendar-day22">|                      | |                      |  <span class="calendar-day">22</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 23" href="/2024/day/23" class="calendar-day23">|                      | |                      |  <span class="calendar-day">23</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 24" href="/2024/day/24" class="calendar-day24">|                      | |                      |  <span class="calendar-day">24</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
<a aria-label="Day 25" href="/2024/day/25" class="calendar-day25">|                      | '-.                  .-'  <span class="calendar-day">25</span> <span class="calendar-mark-complete">*</span><span class="calendar-mark-verycomplete">*</span></a>
'----------------------'   '------------------'
</pre>
</main>

<!-- ga -->
<script>
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','//www.google-analytics.com/analytics.js','ga');
ga('create', 'UA-69522494-1', 'auto');
ga('set', 'anonymizeIp', true);
ga('send', 'pageview');
</script>
<!-- /ga -->
</body>
</html>`
