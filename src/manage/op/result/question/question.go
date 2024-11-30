package question

import "context"

type Question struct {
	Q string

	Answers   []*Answer
	Answer    *Answer

	cancelled bool
	answeredCtx    context.Context
	closeAnswerCtx func()
}

func New(ctx context.Context, q string, answers []*Answer) *Question {
	cctx, cancel := context.WithCancel(ctx)

	return &Question{
		Q:       q,
		Answers: answers,

		answeredCtx:    cctx,
		closeAnswerCtx: cancel,
	}
}

func (q *Question) Wait() {
	select {
	case <-q.answeredCtx.Done():
		// Answered or cancelled
	}
}

func (q *Question) GiveAnswer(answer *Answer) {
	select {
	case <-q.answeredCtx.Done():
		// Already answered (or cancelled)
	default:
		q.Answer = answer
		if answer == nil {
			q.cancelled = true
		}
	}

	q.closeAnswerCtx()
}

func (q *Question) Cancel() {
	select {
	case <-q.answeredCtx.Done():
		// Already answered (or cancelled)
	default:
		q.Answer = nil
		q.cancelled = true
	}

	q.closeAnswerCtx()
}

func (q *Question) Open() bool {
	select {
	case <-q.answeredCtx.Done():
		// Already answered (or cancelled)
		return false
	default:
		return !q.cancelled && q.Answer == nil
	}
}

type Answer struct {
	Key   string
	Title string
}
