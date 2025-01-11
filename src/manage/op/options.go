package op

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"time"
)

type Option struct {
	title       string
	description string
	disabled    bool
	run         func(ctx context.Context, manager *manage.Manager, part *puzzle.Part) result.Result
}

// Generic and debug options

var testr result.Result
var (
	OptionTest = &Option{
		title:       "test options",
		description: "Test options.",
		run: func(ctx context.Context, _ *manage.Manager, _ *puzzle.Part) result.Result {
			if testr == nil {
				testr = result.New()
				testr.SetTotal(100)
			}
			testr.AddRow("Testing")

			go func() {
				for i := 0; ; i++ {
					if i == 5 {
						testr.Ask(question.New(ctx,
							"Is the submitted answer correct?",
							[]*question.Answer{
								{Key: "no", Title: "No"},
								{Key: "yes", Title: "Yes"},
							},
						))
					}

					select {
					case <-ctx.Done():
						log.Info("Context done in test option")
						testr.SetTotal(testr.Steps())
						testr = nil
						return
					case <-time.NewTimer(time.Millisecond * 250).C:
						testr.AddRow(fmt.Sprintf("%d loop", i))
						testr.SetSteps(i)

						if testr.Done() {
							testr = nil
							return
						}
					}
				}
			}()

			return testr
		},
	}
)

// Main options

var (
	OptionLoadFunLines = &Option{
		title:       "update funlines",
		description: "Retrieve the advent calendar.",
		run: func(ctx context.Context, manager *manage.Manager, _ *puzzle.Part) result.Result {
			return manager.LoadRemoteFunLines(ctx)
		},
	}

	OptionCumulativeRuntime = &Option{
		title: "cumulative runtime",
		description: "Compute and print the cumulative runtime",
		run: func(ctx context.Context, manager *manage.Manager, _ *puzzle.Part) result.Result {
			var sum time.Duration
			r := result.New()
			defer r.SetDone()

			r.AddRow("", "", "", "part 1", "part 2")

			for _, p := range manager.Puzzles() {
				p1, p2 := "-", "-"

				if p.Part1 != nil && p.Part1.FastestSolution != nil && p.Part1.FastestSolution.RunResult != nil && p.Part1.FastestSolution.RunResult.Runtime != nil {
					rt := *p.Part1.FastestSolution.RunResult.Runtime
					p1 = rt.String()
					sum += rt
				}

				if p.Part2 != nil && p.Part2.FastestSolution != nil && p.Part2.FastestSolution.RunResult != nil && p.Part2.FastestSolution.RunResult.Runtime != nil {
					rt := *p.Part2.FastestSolution.RunResult.Runtime
					p2 = rt.String()
					sum += rt
				}

				r.AddRow(result.EmojiRunning, p.Title(), "", p1, p2)
			}

			r.AddRow("")
			r.AddRow(result.EmojiCheckMark, "Total runtime: ", "", sum.String())

			return r
		},
	}
)

// Part options

var (
	OptionPartInfo = &Option{
		title:       "info",
		description: "Overview & statistics of the puzzle part.",
		run: func(ctx context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.Info()
		},
	}
	OptionValidatePart = &Option{
		title:       "validate",
		description: "Validate the puzzle part directory.",
		run: func(_ context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.ValidateDir(nil)
		},
	}
	OptionPartBoilerplate = &Option{
		title:       "boilerplate",
		description: "Create the boilerplate puzzle part files.",
		run: func(_ context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.CreateBoilerplatePuzzlePartDir()
		},
	}
	OptionPartLoadRemote = &Option{
		title:       "load remote",
		description: "Load problem statement and puzzle input.",
		run: func(ctx context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.LoadRemote(ctx)
		},
	}
	OptionBuildPart = &Option{
		title:       "build",
		description: "Build the solution.",
		disabled:    true,
		run: func(ctx context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.Build(ctx)
		},
	}
	OptionRunPartSample = &Option{
		title:       "run samples",
		description: "Run the last build against the sample(s).",
		disabled:    true,
		run: func(ctx context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.RunSamples(ctx)
		},
	}
	OptionRunPart = &Option{
		title:       "run input",
		description: "Run the last build against the input",
		disabled:    true,
		run: func(ctx context.Context, _ *manage.Manager, part *puzzle.Part) result.Result {
			return part.RunInput(ctx)
		},
	}
	OptionRecordResult = &Option{
		title:       "record result",
		description: "Record the result of a submitted answer",
		disabled:    true,
		run: func(ctx context.Context, mng *manage.Manager, part *puzzle.Part) result.Result {
			return part.RecordResult(ctx)
		},
	}
)

func (o *Option) FilterValue() string {
	return o.title + " | " + o.description
}

func (o *Option) Title() string {
	return o.title
}

func (o *Option) Description() string {
	return o.description
}

func (o *Option) Disabled() bool {
	return o.disabled
}

func (o *Option) SetDisabled(v bool) *Option {
	o.disabled = v

	return o
}

func (o *Option) Run(ctx context.Context, mng *manage.Manager, p *puzzle.Part) result.Result {
	return o.run(ctx, mng, p)
}
