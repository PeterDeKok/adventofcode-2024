package result

import (
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"strings"
)

type OpResult interface {
	Done() bool
}

type Result interface {
	OpResult() OpResult

	AddRows(value ...[]string) Result
	AddRow(value ...string) Result
	Rows() [][]string

	SetTotal(i int) Result
	AddTotal(i int) Result
	SetSteps(i int) Result
	Increment(diff int) Result
	Total() int
	Steps() int
	Progress() float64
	Ask(q *question.Question)
	Question() *question.Question
	OK() bool
	SetDone() Result
	Done() bool

	AddError(err error, cols ...string) Result
	Error() error

	Listen(fn func(r Result))

	fmt.Stringer
}

type result struct {
	opResult OpResult

	done     bool
	ok       bool
	rows     [][]string
	total    int
	steps    int
	err      error
	question *question.Question
	listener func(rr Result)
}

func New(opResult ...OpResult) Result {
	r := &result{
		rows: make([][]string, 0),
	}

	if len(opResult) > 0 {
		// TODO Clean this up after all other results are gone here...
		r.opResult = opResult[0]
	}

	return r
}

func OrNew(r Result, addTotal int) Result {
	if r == nil {
		r = New()
	}

	r.AddTotal(addTotal)

	return r
}

func (r *result) OpResult() OpResult {
	return r.opResult
}

func (r *result) Rows() [][]string {
	return r.rows
}

func (r *result) addRows(value ...[]string) Result {
	r.rows = append(r.rows, value...)

	return r
}

func (r *result) AddRows(value ...[]string) Result {
	r.rows = append(r.rows, value...)

	if len(value) > 0 {
		r.callListener()
	}

	return r
}

func (r *result) AddRow(value ...string) Result {
	return r.AddRows(value)
}

func (r *result) AddError(err error, cols ...string) Result {
	col2 := "error"
	l := 1
	if len(cols) > 0 {
		col2 = cols[0]
		l = len(cols)
	}

	cells := make([]string, 0, 2+l)
	cells = append(cells,
		EmojiCross+" ",
		col2,
		err.Error(),
	)
	if len(cols) > 0 {
		cells = append(cells, cols[1:]...)
	}

	r.err = err
	r.addRows(cells)
	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) SetTotal(total int) Result {
	r.total = total
	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) AddTotal(i int) Result {
	r.total += i
	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) SetSteps(steps int) Result {
	r.steps = steps
	if r.total < r.steps {
		r.total = r.steps
	}

	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) Increment(diff int) Result {
	r.steps += diff
	if r.total < r.steps {
		r.total = r.steps
	}
	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) Total() int {
	return r.total
}

func (r *result) Steps() int {
	return r.steps
}

func (r *result) Progress() float64 {
	if r.OK() {
		return 1
	}

	return float64(r.steps) / float64(r.total)
}

func (r *result) Ask(q *question.Question) {
	r.question = q

	r.callListener()

	r.question.Wait()

	r.checkOkAndDone()
	r.callListener()
}

func (r *result) Question() *question.Question {
	return r.question
}

func (r *result) OK() bool {
	return r.err == nil && r.steps >= r.total && (r.question == nil || r.question.Answer != nil)
}

func (r *result) SetDone() Result {
	r.done = true
	r.checkOkAndDone()
	r.callListener()

	return r
}

func (r *result) Done() bool {
	return r.done
}

func (r *result) Error() error {
	return r.err
}

func (r *result) Listen(fn func(rr Result)) {
	if r.listener != nil {
		// Already listening
		panic("already listening")
	}

	r.listener = fn
}

func (r *result) callListener() {
	if r.listener != nil {
		r.listener(r)
	}
}

func (r *result) checkOkAndDone() {
	r.done = r.OK() || r.err != nil || (r.steps >= r.total && r.question != nil && !r.question.Open())
}

func (r *result) String() string {
	rows := make([]string, len(r.rows))

	for i, row := range r.rows {
		rows[i] = fmt.Sprint(row)
	}

	return fmt.Sprintf("result: [done: %v] [ok: %v] [listener: %v] [%d/%d → %.0f%%]\nerr: %v\nrows:\n%s\n",
		r.done,
		r.ok,
		r.listener != nil,
		r.steps,
		r.total,
		r.Progress()*100,
		r.err,
		strings.Join(rows, "\n"),
	)
}
