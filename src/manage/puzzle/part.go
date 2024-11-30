package puzzle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"path"
	"path/filepath"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/info"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrBuildWithError = errors.New("failed to buid: part has an error")
	ErrSeeOutputFiles = errors.New("run failed in previous session, see output files")
)

// TODO Remote check if star (& thus OK status)
// TODO Load solution value
// TODO Load first OK solution value as next expected value

type Part struct {
	cnf *PartConfig
	l   *log.Logger
	wg  sync.WaitGroup

	Puzzle    *Puzzle     `json:"-"`
	Nr        int         `json:"nr"`
	Solutions []*Solution `json:"solutions"`

	FastestSolution    *Solution  `json:"-"`
	FirstValidSolution *Solution  `json:"-"`
	LastSolution       *Solution  `json:"-"`
	Star               *Timestamp `json:"star"`

	Validated bool  `json:"-"`
	Error     error `json:"-"`
}

type PartConfig struct {
	Ctx                    context.Context
	Logger                 *log.Logger
	Remote                 *remote.Client
	PuzzlesDir             string
	Puzzle                 *Puzzle
	Nr                     int
	SkipValidationOnCreate bool
}

func CreatePart(cnf *PartConfig) *Part {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "part", "part", cnf.Nr)
	} else {
		l = l.With("type", "part", "part", cnf.Nr)
	}

	p := &Part{
		cnf: cnf,
		l:   l,

		Puzzle:    cnf.Puzzle,
		Nr:        cnf.Nr,
		Solutions: make([]*Solution, 0, 100),
	}

	if ok, err := tools.DirExists(p.Path()); ok || err != nil || !cnf.SkipValidationOnCreate {
		p.ValidateDir(nil)
	}

	return p
}

func (p *Part) FilterValue() string {
	return fmt.Sprintf("%d | %s", p.Nr, p.Puzzle.Day.Format(time.DateOnly))
}

func (p *Part) Title() string {
	return fmt.Sprintf("Day %2d Part %d", p.Puzzle.Day.Day(), p.Nr)
}

func (p *Part) Description() string {
	if p.Error != nil {
		return lipgloss.NewStyle().Foreground(styles.ErrorColor).Render(p.Error.Error())
	}

	// TODO Most of this should move to a custom delegate maybe???
	// View/styling and data should not co-exist... But it is 'easy' (now)......

	str := p.RenderStar()
	str += "    "

	if p.OK() {
		str += lipgloss.NewStyle().Foreground(styles.DimmedColor).Render(p.FastestSolution.RunResult.Runtime.String())
	}

	return str
}

func (p *Part) RenderStar() string {
	if p.OK() {
		return lipgloss.NewStyle().Foreground(styles.StarColor).Render("*")
	}

	return lipgloss.NewStyle().Foreground(styles.VeryDimmedColor).Render("*")
}

func (p *Part) Path() string {
	return path.Join(
		p.Puzzle.Path(),
		fmt.Sprintf("part%d", p.Nr),
	)
}

func (p *Part) OK() bool {
	fvs := p.FirstValidSolution
	fs := p.FastestSolution

	return p.Validated &&
		fvs != nil && fvs.RunResult != nil && (fvs.Status == SolutionStatusValid || fvs.RunResult.OK()) &&
		(fs == fvs || (fs != nil && fs.RunResult != nil && (fvs.Status == SolutionStatusValid || fs.RunResult.OK())))
}

func (p *Part) CanRunSamples() bool {
	return p.Error == nil &&
		p.Validated &&
		p.LastSolution != nil &&
		!p.LastSolution.Status.IsBefore(SolutionStatusBuild) &&
		!p.LastSolution.Status.IsAfter(SolutionStatusSamplesInvalid)
}

func (p *Part) CanRunInput() bool {
	return p.Error == nil &&
		p.Validated &&
		p.LastSolution != nil &&
		!p.LastSolution.Status.IsBefore(SolutionStatusSamplesValid) &&
		!p.LastSolution.Status.IsAfter(SolutionStatusSamplesValid)
}

func (p *Part) CanRecordResult() bool {
	return p.Error == nil &&
		p.Validated &&
		p.LastSolution != nil &&
		p.LastSolution.Status == SolutionStatusReview
}

func (p *Part) Info() result.Result {
	// TODO Remove ones the exact structure of the Result response is known.
	return result.New(p.partInfoTmp()).SetDone()
}

// partInfoTmp is a temporary method until non of the Info() results are used anymore.
func (p *Part) partInfoTmp() *info.PartInfo {
	pi := info.NewPartInfo()

	{ // Summary
		pi.Summary.SolutionCount = strconv.Itoa(len(p.Solutions))

		var bestStatus SolutionStatus = SolutionStatusCreated
		for _, s := range p.Solutions {
			if s.Status.IsAfter(bestStatus) {
				bestStatus = s.Status
			}
		}
		pi.Summary.BestStatus = bestStatus.String()

		if p.OK() {
			rr := p.FastestSolution.RunResult
			pi.Summary.FastestRuntime = rr.Runtime.String()
			pi.Summary.CorrectAnswer = rr.Answer
			if fa := p.FirstValidSolution.RunResult.FinishedAt; !fa.IsUnix() {
				pi.Summary.FinishedAt = fa.Format(time.RFC822)
			} else {
				pi.Summary.FinishedAt = "unix"
			}
		}
	}

	{ // LastSolution
		ls := pi.LastSolution

		if s := p.LastSolution; s != nil {
			ls.Status = s.Status.String()

			if s.RunResult != nil && s.RunResult.Error != nil {
				ls.Error = s.RunResult.Error.Error()
			} else {
				for _, srr := range s.SampleRunResults {
					if srr.Error != nil {
						ls.Error = srr.Error.Error()
					}
				}
			}

			if rr := s.RunResult; rr != nil {
				if rr.Runtime != nil {
					ls.Runtime = rr.Runtime.String()
				}
				if !rr.FinishedAt.IsUnix() {
					ls.FinishedAt = rr.FinishedAt.Format(time.RFC822)
				}
				ls.Answer = rr.Answer
				ls.Expected = rr.Expected
			}
		}
	}

	return pi
}

func (p *Part) ValidateDir(r result.Result) result.Result {
	r = result.OrNew(r, 3)
	defer r.Increment(1)

	if _, err := ValidatePuzzlePartDir(r, p.Path()); err != nil {
		p.Validated = false
		p.Error = err

		return r
	} else {
		r.Increment(1)
	}

	if _, err := p.loadStats(r); err != nil {
		p.Validated = false
		p.Error = err

		return r
	} else {
		r.Increment(1)
	}

	p.Validated = true
	p.Error = nil

	return r
}

func (p *Part) CreateBoilerplatePuzzlePartDir() result.Result {
	emptyFiles := []string{
		"input.txt",
		"PROBLEM.md",
		"sample-expected-1.txt",
		"sample-input-1.txt",
	}

	r := result.New()
	r.AddTotal(7 + len(emptyFiles))
	defer r.SetDone()
	defer r.Increment(1)

	dir := p.Path()

	if ok, err := tools.DirExists(dir); err != nil {
		return r.AddError(err, "dir should not exist")
	} else if ok {
		return r.AddError(os.ErrExist, "dir should not exist")
	} else {
		r.AddRow(result.EmojiCheckMark, "dir should not exist")
		r.Increment(1)
	}

	if err := os.MkdirAll(path.Join(dir, "output"), 0775); err != nil {
		return r.AddError(err, "Puzzle part dir should be created")
	} else {
		r.AddRow(result.EmojiCheckMark, "Puzzle part dir should be created")
		r.Increment(1)
	}

	for _, emptyFile := range emptyFiles {
		if err := writeFile(dir, emptyFile, []byte{}); err != nil {
			return r.AddError(err, fmt.Sprintf("/output/%s should be created", emptyFile))
		} else {
			r.AddRow(result.EmojiCheckMark, fmt.Sprintf("/output/%s should be created", emptyFile))
			r.Increment(1)
		}
	}

	dayPadded := strings.Split(p.Puzzle.Day.Format(time.DateOnly), "-")[2]
	dayUnpadded := strings.TrimPrefix(dayPadded, "0")
	part := strconv.Itoa(p.Nr)

	readmeContentRendered := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(readmeContent,
		"{{day}}", dayUnpadded),
		"{{daypadded}}", dayPadded),
		"{{part}}", part)

	if err := writeFile(dir, "README.md", []byte(readmeContentRendered)); err != nil {
		return r.AddError(err, "/README.md should be generated")
	} else {
		r.AddRow(result.EmojiCheckMark, "/README.md should be generated")
		r.Increment(1)
	}

	if err := writeFile(dir, "solution.go", []byte(solutionContent)); err != nil {
		return r.AddError(err, "/solution.go should be generated")
	} else {
		r.AddRow(result.EmojiCheckMark, "/solution.go should be generated")
		r.Increment(1)
	}

	if err := writeFile(dir, "solution_test.go", []byte(solutionTestContent)); err != nil {
		return r.AddError(err, "/solution_test.go should be generated")
	} else {
		r.AddRow(result.EmojiCheckMark, "/solution_test.go should be generated")
		r.Increment(1)
	}

	statsContentRendered := strings.ReplaceAll(strings.ReplaceAll(statsContent, "{{day}}", dayPadded), "{{part}}", part)
	if err := writeFile(dir, "stats.json", []byte(statsContentRendered)); err != nil {
		return r.AddError(err, "/stats.json should be generated")
	} else {
		r.AddRow(result.EmojiCheckMark, "/stats.json should be generated")
		r.Increment(1)
	}

	return r
}

func (p *Part) saveStats() error {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(p.Path(), "stats.json"), b, 0660); err != nil {
		return fmt.Errorf("failed to save stats: %v", err)
	}

	return nil
}

func (p *Part) loadStats(r result.Result) (result.Result, error) {
	r = result.OrNew(r, 3)
	defer r.Increment(1)

	file, err := os.ReadFile(path.Join(p.Path(), "stats.json"))
	if err != nil {
		return r.AddError(err, "stats.json should be read"), err
	} else {
		r.AddRow(result.EmojiCheckMark, "stats.json should be read")
		r.Increment(1)
	}

	partJson := &Part{}

	if err := json.Unmarshal(file, partJson); err != nil {
		return r.AddError(err, "stats.json should parse"), err
	} else {
		r.AddRow(result.EmojiCheckMark, "stats.json should parse")
		r.Increment(1)
	}

	p.mergeJson(partJson, r)

	return r, r.Error()
}

func (p *Part) mergeJson(partJson *Part, r result.Result) result.Result {
	r = result.OrNew(r, 1)

	if p.Nr != partJson.Nr {
		r.AddError(fmt.Errorf("invalid part nr"), "stats.json part nr match")
		return r
	} else {
		r.AddRow(result.EmojiCheckMark, "stats.json part nr match")
	}

	if partJson.Star != nil {
		r.AddRow(result.EmojiCheckMark, "assign star", "", "replace")
		p.Star = partJson.Star
	} else {
		r.AddRow(result.EmojiEmpty, "assign star", "", "no star")
	}

	if partJson.Solutions == nil {
		p.Solutions = make([]*Solution, 0, 100)
	} else {
		r.AddTotal(len(partJson.Solutions))

		for _, sJson := range partJson.Solutions {
			s := p.getSolution(sJson.Nr)

			if s == nil {
				s = CreateSolution(&SolutionConfig{
					Ctx:    p.cnf.Ctx,
					Logger: p.l,
					Part:   p,
					Nr:     sJson.Nr,
				})

				p.Solutions = append(p.Solutions, s)

				r.AddRow(result.EmojiCheckMark, "retrieve solution", "", fmt.Sprintf("[%d] new", sJson.Nr))
			} else {
				r.AddRow(result.EmojiCheckMark, "retrieve solution", "", fmt.Sprintf("[%d]", sJson.Nr))
			}

			s.mergeJson(sJson, r)

			r.Increment(1)
		}

		slices.SortFunc(p.Solutions, func(a, b *Solution) int {
			return a.Nr - b.Nr
		})
	}

	p.updateLastSolution(r)
	p.updateValidSolutions(r)

	if len(p.Solutions) > 0 {
		r.AddRow(result.EmojiCheckMark, "solutions loaded", "", fmt.Sprintf("[%d x]", len(partJson.Solutions)))
	} else {
		r.AddRow(result.EmojiEmpty, "solutions loaded", "", "none")
	}

	return r.Increment(1)
}

func (p *Part) updateLastSolution(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if len(p.Solutions) == 0 {
		p.LastSolution = nil
		r.AddRow(result.EmojiEmpty, "last solution updated", "", "none")
		return r
	}

	p.LastSolution = p.Solutions[len(p.Solutions)-1]
	r.AddRow(result.EmojiCheckMark, "last solution updated", "", "id: "+p.LastSolution.NrStr())

	return r
}

func (p *Part) updateValidSolutions(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	var fastest *Solution
	var first *Solution

	for _, s := range p.Solutions {
		if s.Status != SolutionStatusValid || s.RunResult == nil {
			continue
		}

		if fastest == nil || *s.RunResult.Runtime < *fastest.RunResult.Runtime {
			fastest = s
		}

		if first == nil || s.RunResult.FinishedAt.Before(first.RunResult.FinishedAt) {
			first = s
		}
	}

	if fastest == nil {
		r.AddRow(result.EmojiEmpty, "fastest solution updated", "", "none")
	} else {
		r.AddRow(result.EmojiCheckMark, "fastest solution updated", "", "id: "+fastest.NrStr())
	}

	if first == nil {
		r.AddRow(result.EmojiEmpty, "first valid solution updated", "", "none")
	} else {
		r.AddRow(result.EmojiCheckMark, "first valid solution updated", "", "id: "+first.NrStr())
	}

	p.FastestSolution = fastest
	p.FirstValidSolution = first

	return r
}

func (p *Part) LoadRemote(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(1)

	go func() {
		defer r.SetDone()

		if p == p.Puzzle.Part2 && !p.Puzzle.Part1.OK() {
			r.AddError(fmt.Errorf("unfinished"), "previous part finished")
			return
		} else if p == p.Puzzle.Part2 {
			r.AddRow(result.EmojiCheckMark, "previous part finished")
		}

		p.loadRemoteInput(ctx, r)
		if r.Error() != nil {
			return
		}

		p.loadRemoteProblem(ctx, r)

		r.Increment(1)
	}()

	return r
}

func (p *Part) loadRemoteInput(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	inputUri := fmt.Sprintf("%d/day/%d/input", p.Puzzle.Day.Year(), p.Puzzle.Day.Day())
	file := "input.txt"

	f, err := os.OpenFile(path.Join(p.Path(), file), os.O_WRONLY, 0660)
	if err != nil {
		return r.AddError(err, "file should still exist and be writable", file)
	} else {
		r.AddRow(result.EmojiCheckMark, "file should still exist and be writable", "", file)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	return p.loadRemote(ctx, r, file, inputUri, f, remote.RateInput)
}

func (p *Part) loadRemoteProblem(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	problemUri := fmt.Sprintf("%d/day/%d", p.Puzzle.Day.Year(), p.Puzzle.Day.Day())
	file := "PROBLEM.md"

	body := bytes.NewBuffer(make([]byte, 0, 4096))

	r = p.loadRemote(ctx, r, file, problemUri, body, remote.RateProblem)

	// TODO Parse
	err := os.WriteFile(path.Join(p.Path(), "PROBLEM.html"), body.Bytes(), 0660)
	if err != nil {
		return r
	}

	return r
}

func (p *Part) loadRemote(ctx context.Context, r result.Result, file string, uri string, w io.Writer, rateLimitCategory string) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	fp := path.Join(p.Path(), file)

	if ok, err := tools.FileExists(fp); err != nil {
		return r.AddError(err, "file should exist", file)
	} else if !ok {
		return r.AddError(fmt.Errorf("not found"), "file should exist", file)
	} else if fc, err := os.ReadFile(fp); err != nil {
		return r.AddError(err, "file should exist", file)
	} else if len(fc) > 0 {
		return r.AddError(fmt.Errorf("file is not empty"), "file should exist", file)
	} else {
		r.AddRow(result.EmojiCheckMark, "file should exist", "", file)
	}

	if p.cnf.Remote == nil {
		return r.AddError(fmt.Errorf("remote not set"), "select remote", file)
	} else {
		r.AddRow(result.EmojiCheckMark, "select remote", "", file)
	}

	body, err := p.cnf.Remote.Get(ctx, uri, &remote.RequestOptions{
		RateLimitCategory: rateLimitCategory,
	})
	if err != nil {
		return r.AddError(err, "retrieve content from remote", file)
	} else {
		r.AddRow(result.EmojiCheckMark, "retrieve content from remote", "", file)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	if _, err := io.Copy(w, body); err != nil {
		return r.AddError(err, "should be written to file", file)
	} else {
		r.AddRow(result.EmojiCheckMark, "should be written to file", "", file)
	}

	return r
}

func (p *Part) Build(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(6)

	go func() {
		defer r.SetDone()
		defer r.Increment(1)

		p.ValidateDir(r)
		if r.Error() != nil {
			p.Error = ErrBuildWithError
			return
		} else {
			r.Increment(1)
		}

		p.l.Info("after part validate", "s", r.Steps(), "t", r.Total())
		r.AddRow(result.EmojiCheckMark, "loaded dir")

		nr := 0
		if p.LastSolution != nil {
			nr = p.LastSolution.Nr + 1
		}
		r.AddRow(result.EmojiCheckMark, "next solution nr", "", fmt.Sprintf("[id: %d]", nr))

		s := CreateSolution(&SolutionConfig{
			Ctx:    p.cnf.Ctx,
			Logger: p.l,
			Part:   p,
			Nr:     nr,
		})

		p.Solutions = append(p.Solutions, s)

		if err := p.saveStats(); err != nil {
			r.AddError(err, "save submission pre-build state")
			p.Error = err
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save submission pre-build state")
			r.Increment(1)
		}

		s.Fmt(ctx, r)
		if err := r.Error(); err != nil {
			r.AddError(err, "go fmt")
			p.Error = err
			p.l.Error("failed to go fmt part", "err", err)
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "go fmt")
			r.Increment(1)
		}

		// TODO go vet?

		s.Build(ctx, r)
		if err := r.Error(); err != nil {
			r.AddError(err, "go build")
			p.Error = err
			p.l.Error("failed to go build part", "err", err)
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "go build")
			r.Increment(1)
		}

		if err := p.saveStats(); err != nil {
			r.AddError(err, "save submission post-build state")
			p.Error = err
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save submission post-build state")
			r.Increment(1)
		}

		p.LastSolution = s
	}()

	return r
}

func (p *Part) RunSamples(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(5)

	go func() {
		defer r.SetDone()
		defer r.Increment(1)

		if p.Error != nil {
			r.AddError(ErrSeeOutputFiles, "pre-req check")
			return
		}

		p.ValidateDir(r)
		if r.Error() != nil {
			p.Error = ErrRunWithError
			return
		} else {
			r.Increment(1)
		}

		p.l.Info("after part validate", "s", r.Steps(), "t", r.Total())
		r.AddRow(result.EmojiCheckMark, "loaded dir")

		s := p.LastSolution

		if s == nil {
			r.AddError(fmt.Errorf("no Solutions to run"), "last solution")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "last solution", "", fmt.Sprintf("[%d]", s.Nr))
			r.Increment(1)
		}

		s.RunSamples(ctx, r)
		r.Increment(1)

		if err := p.saveStats(); err != nil {
			r.AddError(err, "save solution post-run state")
			p.Error = err
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save solution post-run state")
			r.Increment(1)
		}
	}()

	return r
}

func (p *Part) RunInput(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(6)

	go func() {
		defer r.SetDone()
		defer r.Increment(1)

		if p.Error != nil {
			r.AddError(ErrSeeOutputFiles, "pre-req check")
			return
		}

		p.ValidateDir(r)
		if r.Error() != nil {
			p.Error = ErrRunWithError
			return
		} else {
			r.Increment(1)
		}

		p.l.Info("after part validate", "s", r.Steps(), "t", r.Total())
		r.AddRow(result.EmojiCheckMark, "loaded dir")

		s := p.LastSolution

		if s == nil {
			r.AddError(fmt.Errorf("no Solutions to run"), "last solution")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "last solution", "", fmt.Sprintf("[%d]", s.Nr))
			r.Increment(1)
		}

		s.RunInput(ctx, r)
		r.Increment(1)

		rr := s.RunResult

		// TODO Star TS
		if rr.OK() {
			if rr.Runtime != nil && (p.FastestSolution == nil || *rr.Runtime < *p.FastestSolution.RunResult.Runtime) {
				p.FastestSolution = s

				emptyTiming := []byte("## Timing\n\n```\n\n```\n")
				filledTiming := []byte(fmt.Sprintf("## Timing\n\n```\n%s\n```\n", s.RunResult.Runtime.String()))

				b, err := os.ReadFile(path.Join(p.Path(), "README.md"))
				if err != nil {
					r.AddError(err, "read README.md")
					return
				} else {
					r.AddRow(result.EmojiCheckMark, "read README.md")
				}

				b2 := bytes.Replace(b, emptyTiming, filledTiming, 1)

				if err := os.WriteFile(path.Join(p.Path(), "README.md"), b2, 0660); err != nil {
					r.AddError(err, "save README.md")
				} else {
					r.AddRow(result.EmojiCheckMark, "save README.md")
				}
			}

			if p.FirstValidSolution == nil {
				p.FirstValidSolution = s
			}
		}

		r.Increment(1)

		if err := p.saveStats(); err != nil {
			r.AddError(err, "save submission post-run state")
			p.Error = err
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save submission post-run")
			r.Increment(1)
		}
	}()

	return r
}

func (p *Part) getSolution(nr int) *Solution {
	if p.Solutions == nil {
		return nil
	}

	for _, s := range p.Solutions {
		if s.Nr == nr {
			return s
		}
	}

	return nil
}

func (p *Part) RecordResult(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(6)

	go func() {
		defer r.SetDone()
		defer r.Increment(1)

		s := p.LastSolution
		if s == nil {
			r.AddError(fmt.Errorf("no Solutions to run"), "last solution")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "last solution")
			r.Increment(1)
		}

		if s.Status != SolutionStatusReview {
			r.AddError(fmt.Errorf("solution is not in review"), "solution status", s.Status.String())
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "solution status")
			r.Increment(1)
		}

		aNo := &question.Answer{Key: "no", Title: "No"}
		aYes := &question.Answer{Key: "yes", Title: "Yes"}

		q := question.New(ctx,
			"Is the submitted answer correct?",
			[]*question.Answer{aNo, aYes},
		)

		r.AddRow(result.EmojiInput, "asking result")

		r.Ask(q)

		if q.Answer == nil {
			r.AddError(fmt.Errorf("question aborted"), "receive answer")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "receive answer", "", q.Answer.Title)
			r.Increment(1)
		}

		switch q.Answer {
		case aNo:
			s.Status = SolutionStatusInvalid
		case aYes:
			s.Status = SolutionStatusValid
			p.updateValidSolutions(r)
		}

		if err := p.saveStats(); err != nil {
			r.AddError(err, "save submission post-result state")
			p.Error = err
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save submission post-result state")
			r.Increment(1)
		}

		if s != p.FirstValidSolution && s != p.FastestSolution {
			r.AddRow(result.EmojiNote, "save /output/answer.txt", "", "not the first valid solution")
			r.AddRow(result.EmojiNote, "update /README.md", "", "not the first valid solution or fastest solution")
			return
		}

		if s != p.FirstValidSolution {
			r.AddRow(result.EmojiNote, "save /output/answer.txt", "", "not the first valid solution")
		} else if s.Status != SolutionStatusValid {
			r.AddError(fmt.Errorf("status not valid"), "save /output/answer.txt")
			return
		} else if s.RunResult == nil {
			r.AddError(ErrNoRunResult, "save /output/answer.txt")
			return
		} else if len(s.RunResult.Answer) == 0 {
			r.AddError(ErrNoAnswer, "save /output/answer.txt")
			return
		} else if err := writeFile(p.Path(), "output/answer.txt", []byte(s.RunResult.Answer)); err != nil {
			r.AddError(err, "save /output/answer.txt")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "save /output/answer.txt")
			r.Increment(1)
		}

		emptyFinished := []byte("| Finished on |                         |")
		filledFinished := []byte(fmt.Sprintf("| Finished on | %s |", s.RunResult.FinishedAt.Format("2006-01-02 15:04:05 MST")))

		if s != p.FirstValidSolution {
			r.AddRow(result.EmojiNote, "save /README.md", "", "not the first valid solution")
		} else {
			b, err := os.ReadFile(path.Join(p.Path(), "README.md"))
			if err != nil {
				r.AddError(err, "read README.md")
				return
			} else {
				r.AddRow(result.EmojiCheckMark, "read README.md")
			}

			b2 := bytes.Replace(b, emptyFinished, filledFinished, 1)

			if err := os.WriteFile(path.Join(p.Path(), "README.md"), b2, 0660); err != nil {
				r.AddError(err, "save README.md")
			} else {
				r.AddRow(result.EmojiCheckMark, "save README.md")
			}
		}

		emptyTiming := []byte("## Timing\n\n```\n\n```\n")
		filledTiming := []byte(fmt.Sprintf("## Timing\n\n```\n%s\n```\n", s.RunResult.Runtime.String()))

		if s != p.FastestSolution {
			r.AddRow(result.EmojiNote, "update /README.md", "", "not the first valid solution or fastest solution")
		} else {
			b, err := os.ReadFile(path.Join(p.Path(), "README.md"))
			if err != nil {
				r.AddError(err, "read README.md")
				return
			} else {
				r.AddRow(result.EmojiCheckMark, "read README.md")
			}

			b2 := bytes.Replace(b, emptyTiming, filledTiming, 1)

			if err := os.WriteFile(path.Join(p.Path(), "README.md"), b2, 0660); err != nil {
				r.AddError(err, "save README.md")
			} else {
				r.AddRow(result.EmojiCheckMark, "save README.md")
			}
		}
	}()

	return r
}

func ValidatePuzzlePartDir(r result.Result, dir string) (result.Result, error) {
	files := []string{
		"input.txt",
		"PROBLEM.md",
		"README.md",
		"sample-expected-1.txt",
		"sample-input-1.txt",
		"solution.go",
		"solution_test.go",
		"stats.json",
	}

	r = result.OrNew(r, 2+len(files))
	defer r.Increment(1)

	if ok, err := tools.DirExists(dir); err != nil {
		return r.AddError(err, "dir should exist"), err
	} else if !ok {
		return r.AddError(os.ErrNotExist, "dir should exist"), os.ErrNotExist
	} else {
		r.AddRow(result.EmojiCheckMark, "dir should exist")
	}
	r.Increment(1)

	for _, f := range files {
		if ok, err := tools.FileExists(filepath.Join(dir, f)); err != nil {
			r.AddError(err, fmt.Sprintf("file %s should exist", f))
		} else if !ok {
			r.AddError(os.ErrNotExist, fmt.Sprintf("file %s should exist", f))
		} else {
			r.AddRow(result.EmojiCheckMark, fmt.Sprintf("file %s should exist", f))
		}
		r.Increment(1)
	}

	return r, r.Error()
}

func writeFile(dir, name string, data []byte) error {
	fp := path.Join(dir, name)

	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
