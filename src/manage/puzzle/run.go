package puzzle

import (
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"time"
)

type RunResult struct {
	InputFile  string         `json:"input_file"`
	Error      error          `json:"error"`
	Answer     string         `json:"-"`
	Runtime    *time.Duration `json:"runtime"`
	FinishedAt *Timestamp     `json:"finishedAt"`
	Expected   string         `json:"-"`
}

func (rr *RunResult) OK() bool {
	if rr == nil || rr.Error != nil || len(rr.Answer) == 0 {
		return false
	}

	return rr.Answer == rr.Expected
}

func (rr *RunResult) Deferred() bool {
	if rr == nil || rr.Error != nil || len(rr.Answer) == 0 {
		return false
	}

	return len(rr.Expected) == 0
}

func (rr *RunResult) mergeJson(rrJson *RunResult, r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if len(rrJson.InputFile) > 0 {
		rr.InputFile = rrJson.InputFile
	}
	if rrJson.Error != nil {
		rr.Error = rrJson.Error
	}
	if len(rrJson.Answer) > 0 {
		rr.Answer = rrJson.Answer
	}
	if rrJson.Runtime != nil && *rrJson.Runtime > 0 {
		rr.Runtime = rrJson.Runtime
	}
	if rrJson.FinishedAt != nil {
		rr.FinishedAt = rrJson.FinishedAt
	}
	if len(rrJson.Expected) > 0 {
		rr.Expected = rrJson.Expected
	}

	return r
}

func (rr *RunResult) ToSampleRunResult(name string, expected string, err error) *SampleRunResult {
	return &SampleRunResult{
		InputFile:  name,
		Error:      err,
		Answer:     rr.Answer,
		Timing:     rr.Runtime,
		FinishedAt: rr.FinishedAt,
		Expected:   expected,
	}
}

type SampleRunResult struct {
	InputFile  string         `json:"input_file"`
	Error      error          `json:"error"`
	Answer     string         `json:"answer"`
	Timing     *time.Duration `json:"timing"`
	FinishedAt *Timestamp     `json:"finishedAt"`
	Expected   string         `json:"expected"`
}

func (rr *SampleRunResult) OK() bool {
	if rr == nil || rr.Error != nil || len(rr.Answer) == 0 {
		return false
	}

	return rr.Answer == rr.Expected
}

func (rr *SampleRunResult) Deferred() bool {
	if rr == nil || rr.Error != nil || len(rr.Answer) == 0 {
		return false
	}

	return len(rr.Expected) == 0
}

func (rr *SampleRunResult) mergeJson(rrJson *SampleRunResult, r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if len(rrJson.InputFile) > 0 {
		rr.InputFile = rrJson.InputFile
	}
	if rrJson.Error != nil {
		rr.Error = rrJson.Error
	}
	if len(rrJson.Answer) > 0 {
		rr.Answer = rrJson.Answer
	}
	if rrJson.Timing != nil && *rrJson.Timing > 0 {
		rr.Timing = rrJson.Timing
	}
	if rrJson.FinishedAt != nil {
		rr.FinishedAt = rrJson.FinishedAt
	}
	if len(rrJson.Expected) > 0 {
		rr.Expected = rrJson.Expected
	}

	return r
}

type SolutionFilepair struct {
	Name     string
	Error    error
	In       io.Reader
	Expected string
}
