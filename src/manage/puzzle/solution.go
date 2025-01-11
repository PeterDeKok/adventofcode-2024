package puzzle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	plugin2 "peterdekok.nl/adventofcode/twentytwentyfour/src/manage/plugin"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"plugin"
	"strconv"
	"strings"
	"time"
)

var (
	ErrAlreadyRan   = errors.New("failed to build: solution was already run")
	ErrNoAnswer     = errors.New("no answer available")
	ErrNoInputRun   = errors.New("no input run")
	ErrNoRunResult  = errors.New("no run result")
	ErrNoSamples    = errors.New("failed to run: no samples loaded")
	ErrRunWithError = errors.New("failed to run: solution has an error")
)

type Solution struct {
	cnf *SolutionConfig
	l   *log.Logger

	part      *Part
	RunLogger *logger.IterationLogger `json:"-"`

	Nr int `json:"nr"`

	Status           SolutionStatus     `json:"status"`
	RunResult        *RunResult         `json:"run_result"`
	SampleRunResults []*SampleRunResult `json:"sample_run_results"`
}

type SolutionConfig struct {
	Ctx    context.Context
	Logger *log.Logger
	Part   *Part
	Nr     int
}

type SolutionFn func(ctx context.Context, l *logger.IterationLogger, in io.Reader, out io.Writer) error

func CreateSolution(cnf *SolutionConfig) *Solution {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "solution", "Nr", cnf.Nr)
	} else {
		l = l.With("type", "solution", "Nr", cnf.Nr)
	}

	return &Solution{
		cnf: cnf,
		l:   l,

		part:      cnf.Part,
		RunLogger: logger.CreateIterationLogger(cnf.Ctx),

		Nr: cnf.Nr,

		Status:           SolutionStatusCreated,
		SampleRunResults: make([]*SampleRunResult, 0, 5),
	}
}

func (s *Solution) Path() string {
	return path.Join(s.part.Path(), s.RelPath())
}

func (s *Solution) RelPath() string {
	return path.Join("output", fmt.Sprintf("%04d-solution", s.cnf.Nr))
}

func (s *Solution) mergeJson(sJson *Solution, r result.Result) result.Result {
	r = result.OrNew(r, 1)

	if s.Nr != sJson.Nr {
		r.AddError(fmt.Errorf("invalid solution nr"), "stats.json solution nr match")
		return r
	} else {
		r.AddRow(result.EmojiCheckMark, "stats.json solution nr match")
	}

	s.Status = sJson.Status

	if sJson.RunResult == nil {
		s.RunResult = nil
	} else {
		if s.RunResult == nil {
			s.RunResult = sJson.RunResult
		} else {
			s.RunResult.mergeJson(sJson.RunResult, r)
		}

		if s.Status == SolutionStatusReview || s.Status == SolutionStatusValid {
			s.loadAnswer(r)
		}
		s.loadExpected(r)
	}
	r.AddRow(result.EmojiCheckMark, "solution run result updated")

	if sJson.SampleRunResults == nil {
		s.SampleRunResults = make([]*SampleRunResult, 0, 5)
	} else {
		s.SampleRunResults = sJson.SampleRunResults

		for _, srr := range s.SampleRunResults {
			if srr.ErrorStr != nil && srr.Error == nil {
				srr.Error = errors.New(*srr.ErrorStr)
			}
		}
	}
	r.AddRow(result.EmojiCheckMark, "solution sample run results updated", "", fmt.Sprintf("[%d x]", len(sJson.SampleRunResults)))
	r.Increment(1)

	return r.Increment(1)
}

func (s *Solution) Fmt(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 3)
	defer r.Increment(1)

	cctx, cancel := context.WithCancel(ctx)
	done := context.AfterFunc(s.cnf.Ctx, cancel)
	defer done()

	stdoutFile, stderrFile, cls, err := s.pipeDescriptors("fmt")
	defer cls() // Intentional before err check
	if err != nil {
		return r.AddError(err, "create pipe descriptors")
	} else {
		r.AddRow(result.EmojiCheckMark, "create pipe descriptors")
		r.Increment(1)
	}

	err = build.GoFmt(&build.Config{
		Ctx: cctx,

		WorkingDir: s.part.Path(),

		Stdout: stdoutFile,
		Stderr: stderrFile,
	})

	if err != nil {
		s.l.With("err", err).Error("failed to format solution")
		return r.AddError(err, "format solution")
	} else {
		r.AddRow(result.EmojiCheckMark, "format solution")
		r.Increment(1)
	}

	return r
}

func (s *Solution) Build(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 5)
	defer r.Increment(1)

	if s.Status.IsAfter(SolutionStatusCreated) {
		return r.AddError(ErrAlreadyRan, "solution build status check", s.Status.String())
	} else {
		r.AddRow(result.EmojiCheckMark, "solution build status check", "", s.Status.String())
	}

	cctx, cancel := context.WithCancel(ctx)
	done := context.AfterFunc(s.cnf.Ctx, cancel)
	defer done()

	if _, err := os.Stat(s.Path()); !os.IsNotExist(err) {
		return r.AddError(os.ErrExist, "check existing file")
	} else {
		r.AddRow(result.EmojiCheckMark, "check existing file")
		r.Increment(1)
	}

	stdoutFile, stderrFile, cls, err := s.pipeDescriptors("build")
	defer cls() // Intentional before err check
	if err != nil {
		return r.AddError(err, "create pipe descriptors")
	} else {
		r.AddRow(result.EmojiCheckMark, "create pipe descriptors")
		r.Increment(1)
	}

	_, tmpDir := s.copySolutionCode(r)
	if r.Error() != nil {
		return r.AddError(fmt.Errorf("failed to copy solution"), "copy solution")
	} else {
		r.AddRow(result.EmojiCheckMark, "copy solution")
		r.Increment(1)
	}

	err = build.GoPlugin(&build.Config{
		Ctx: cctx,

		WorkingDir: tmpDir,
		OutputFile: s.Path(),

		Stdout: stdoutFile,
		Stderr: stderrFile,
	})

	if err != nil {
		s.Status = SolutionStatusBuildFailed
		s.l.With("err", err).Error("failed to build solution")
		return r.AddError(err, "build solution")
	} else {
		r.AddRow(result.EmojiCheckMark, "build solution")
		r.Increment(1)
	}

	s.Status = SolutionStatusBuild

	return r
}

func (s *Solution) RunSamples(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 3)
	defer r.Increment(1)

	if s.Status.IsBefore(SolutionStatusBuild) {
		return r.AddError(fmt.Errorf("solution is not build"), "solution status", s.Status.String())
	} else if s.Status.IsAfter(SolutionStatusSamplesInvalid) {
		// TODO Probably should make sure every individual sample file is checked
		// Instead of 'any'. (And only run the non-run ones)
		return r.AddError(fmt.Errorf("solution already ran samples"), "solution status", s.Status.String())
	} else {
		r.AddRow(result.EmojiCheckMark, "solution status", "", s.Status.String())
	}

	if _, err := os.Stat(s.Path()); err != nil {
		return r.AddError(err, "check solution build file")
	} else {
		r.AddRow(result.EmojiCheckMark, "check solution build file")
		r.Increment(1)
	}

	_, sampleFiles := s.GetSampleFiles(r)

	sampleCount := len(sampleFiles)
	if sampleCount == 0 {
		s.Status = SolutionStatusSamplesInvalid
		return r.AddError(ErrNoSamples, "load samples")
	} else {
		r.AddRow(result.EmojiCheckMark, "load samples", "", fmt.Sprintf("[%dx]", sampleCount))
		r.AddTotal(sampleCount)
	}

	for i, sf := range sampleFiles {
		desc := fmt.Sprintf("%d: %s", i+1, sf.Name)

		r.AddRow(result.EmojiCheckMark, "Running solution against sample input", "", desc)

		if sf.Error != nil {
			r.AddError(sf.Error, "sample files read", desc)
			continue
		} else {
			r.AddRow(result.EmojiCheckMark, "sample files read", "", desc)
		}

		_, rr, err := s.run(ctx, r, sf)

		if rr == nil {
			r.AddError(err, "sample run failed to start", desc)
			continue
		}

		srr := rr.ToSampleRunResult(sf.Name, sf.Expected, err)
		s.SampleRunResults = append(s.SampleRunResults, srr)

		if err != nil {
			r.AddError(err, "sample run failed", desc)
		} else if srr.OK() {
			r.AddRow(result.EmojiCheckMark, "answer OK", "", srr.Answer)
		} else {
			r.AddRow(result.EmojiCross, "answer WRONG", srr.Answer, srr.Expected)
		}

		r.Increment(1)
	}

	hasInvalid := false
	for _, srr := range s.SampleRunResults {
		if srr.Error != nil || !srr.OK() {
			hasInvalid = true
		}
	}

	if r.Error() != nil || len(s.SampleRunResults) != len(sampleFiles) || hasInvalid {
		s.Status = SolutionStatusSamplesInvalid
	} else {
		s.Status = SolutionStatusSamplesValid
	}

	r.AddRow(result.EmojiCheckMark, "all samples checked", "", fmt.Sprintf("[%dx]", sampleCount))

	return r
}

func (s *Solution) RunInput(ctx context.Context, r result.Result) result.Result {
	r = result.OrNew(r, 4)
	defer r.Increment(1)

	if s.Status.IsBefore(SolutionStatusSamplesValid) {
		return r.AddError(fmt.Errorf("solution not verified against samples"), "solution status", s.Status.String())
	} else if s.Status.IsAfter(SolutionStatusSamplesValid) {
		// TODO Probably should make sure every individual sample file is checked
		// Instead of 'any'. (And only run the non-run ones)
		return r.AddError(fmt.Errorf("solution already ran against input"), "solution status", s.Status.String())
	} else {
		r.AddRow(result.EmojiCheckMark, "solution status", "", s.Status.String())
	}

	if _, err := os.Stat(s.Path()); err != nil {
		return r.AddError(err, "check solution build file")
	} else {
		r.AddRow(result.EmojiCheckMark, "check solution build file")
		r.Increment(1)
	}

	_, sf := s.GetInputfile(r)

	if r.Error() != nil {
		s.Status = SolutionStatusInvalid
		return r
	} else if sf == nil {
		s.Status = SolutionStatusInvalid
		return r.AddError(fmt.Errorf("failed to get input file"), "get input file")
	} else {
		r.AddRow(result.EmojiCheckMark, "get input file")
		r.Increment(1)
	}

	if sf.Error != nil {
		s.Status = SolutionStatusInvalid
		return r.AddError(sf.Error, "input file read", sf.Name)
	} else {
		r.AddRow(result.EmojiCheckMark, "input file read", "", sf.Name)
		r.Increment(1)
	}

	r.AddRow(result.EmojiCheckMark, "Running solution against input", "", sf.Name)

	_, rr, err := s.run(ctx, r, sf)

	if rr == nil {
		s.Status = SolutionStatusInvalid
		return r.AddError(err, "run failed to start")
	}

	rr.InputFile = strings.TrimLeft(strings.TrimPrefix(sf.Name, s.part.Path()), "/")
	rr.Expected = sf.Expected
	rr.Error = err
	if err != nil {
		es := err.Error()
		rr.ErrorStr = &es
	}
	s.RunResult = rr

	if err != nil {
		s.Status = SolutionStatusInvalid
		return r.AddError(err, "run failed")
	} else if rr.OK() {
		s.Status = SolutionStatusValid
		r.AddRow(result.EmojiCheckMark, "answer OK", "", rr.Answer)
		s.saveAnswer(r)
	} else if len(sf.Expected) == 0 {
		s.Status = SolutionStatusReview
		r.AddRow(result.EmojiInput, "answer in review", "", rr.Answer)
		s.saveAnswer(r)
	} else {
		s.Status = SolutionStatusInvalid
		r.AddRow(result.EmojiCross, "answer WRONG", rr.Answer, rr.Expected)
	}

	return r
}

func (s *Solution) run(ctx context.Context, r result.Result, in *SolutionFilepair) (result.Result, *RunResult, error) {
	r = result.OrNew(r, 7)
	defer r.Increment(1)

	if in == nil {
		return r.AddError(fmt.Errorf("missing"), "check input file pair"), nil, r.Error()
	} else {
		r.AddRow(result.EmojiCheckMark, "check input file pair")
		r.Increment(1)
	}

	cctx, cancel := context.WithCancel(ctx)
	done := context.AfterFunc(s.cnf.Ctx, cancel)
	defer done()

	stderrFile, cls, err := s.pipeDescriptorStderr("run")
	defer cls() // Intentional before err check
	if err != nil {
		return r.AddError(err, "create stderr pipe"), nil, r.Error()
	} else {
		r.AddRow(result.EmojiCheckMark, "create stderr pipe")
		r.Increment(1)
	}

	p, err := plugin.Open(s.Path())
	if err != nil {
		return r.AddError(err, "open solution as plugin"), nil, r.Error()
	} else {
		r.AddRow(result.EmojiCheckMark, "open solution as plugin")
		r.Increment(1)
	}

	if pPtr, err := p.Lookup("Pre"); err != nil {
		// Not exisiting (or loadable) is OK
		r.AddRow(result.EmojiEmpty, "pre-solution func from plugin", "", "Pre func symbol not found")
		r.Increment(1)
	} else if pre, ok := pPtr.(func(run string)); !ok {
		// Existing, but not casting is a failure
		return r.AddError(err, "pre-solution func from plugin"), nil, r.Error()
	} else {
		pre(in.Name)

		r.AddRow(result.EmojiCheckMark, "pre-solution func from plugin")
		r.Increment(1)
	}

	sPtr, err := p.Lookup("Solution")
	if err != nil {
		return r.AddError(err, "lookup solution func from plugin"), nil, r.Error()
	} else {
		r.AddRow(result.EmojiCheckMark, "lookup solution func from plugin")
		r.Increment(1)
	}

	solution := sPtr.(func(ctx context.Context, l *logger.IterationLogger, input io.Reader, w io.Writer) error)

	// TODO Register s.RunLogger listener????

	// TODO When another reporting structure is created,
	// Replace this with an 'in-progress' item.
	r.AddRow(result.EmojiRunning, "running solution")

	runResult, err := solutionWrapper(solution, cctx, s.RunLogger, in.In)
	if err != nil {
		if _, err := fmt.Fprintf(stderrFile, "%s", err); err != nil {
			s.l.Error("Failed to write to stderr log file")
		}

		return r.AddError(err, "run solution"), runResult, r.Error()
	} else {
		r.AddRow(result.EmojiCheckMark, "run solution")
		r.Increment(1)
	}

	return r, runResult, nil
}

func (s *Solution) GetSampleFiles(r result.Result) (result.Result, []*SolutionFilepair) {
	dir, err := os.ReadDir(s.part.Path())

	if err != nil {
		r.AddError(err, "read dir for samples")
		return r, []*SolutionFilepair{}
	} else {
		r.AddRow(result.EmojiCheckMark, "read dir for samples")
	}

	files := make([]*SolutionFilepair, 0, 10)

	for _, de := range dir {
		if de.IsDir() {
			continue
		}

		sf := &SolutionFilepair{
			Name: de.Name(),
		}

		sid := ""

		if noExt, ok := strings.CutSuffix(de.Name(), ".txt"); !ok {
			continue
		} else if sid, ok = strings.CutPrefix(noExt, "sample-input-"); !ok || len(sid) == 0 {
			continue
		}

		ifp := path.Base(de.Name())
		efp := fmt.Sprintf("sample-expected-%s.txt", sid)

		if inpB, err := os.ReadFile(path.Join(s.part.Path(), de.Name())); err != nil {
			sf.Error = err
			r.AddError(err, "load input for sample", fmt.Sprintf("%s: %s", sid, ifp))
		} else if expB, err := os.ReadFile(path.Join(s.part.Path(), efp)); err != nil {
			sf.Error = err
			r.AddRow(result.EmojiCheckMark, "load input for sample", "", fmt.Sprintf("%s: %s", sid, ifp))
			r.AddError(err, "load expected for sample", fmt.Sprintf("%s: %s", sid, efp))
		} else {
			sf.In = bytes.NewBuffer(inpB)
			sf.Expected = string(expB)
			r.AddRow(result.EmojiCheckMark, "load input for sample", "", fmt.Sprintf("%s: %s", sid, ifp))
			r.AddRow(result.EmojiCheckMark, "load expected for sample", "", fmt.Sprintf("%s: %s", sid, efp))
		}

		files = append(files, sf)
	}

	return r, files
}

func (s *Solution) GetInputfile(r result.Result) (result.Result, *SolutionFilepair) {
	ifp := path.Join(s.part.Path(), "input.txt")
	ifpRel := strings.TrimLeft(strings.TrimPrefix(ifp, s.part.Path()), "/")
	efp := path.Join(path.Dir(s.Path()), "answer.txt")
	efpRel := strings.TrimLeft(strings.TrimPrefix(efp, s.part.Path()), "/")

	sf := &SolutionFilepair{
		Name: "input.txt",
	}

	if inpB, err := os.ReadFile(ifp); err != nil {
		sf.Error = err
		r.AddError(err, "load input", ifpRel)
		return r, nil
	} else {
		sf.In = bytes.NewBuffer(inpB)
		r.AddRow(result.EmojiCheckMark, "load input", "", ifpRel)
	}

	if expB, err := os.ReadFile(efp); err != nil && !os.IsNotExist(err) {
		sf.Error = err
		r.AddError(err, "load previous answer", efpRel)
	} else if err == nil {
		sf.Expected = string(expB)
		r.AddRow(result.EmojiCheckMark, "load previous answer", "", efpRel)
	} else {
		r.AddRow(result.EmojiEmpty, "no previous answer", "", efpRel)
	}

	return r, sf
}

func (s *Solution) pipeDescriptors(t string) (stdout, stderr io.Writer, cls func(), err error) {
	stdOutCls, stdErrCls := func() {
		/* No-op */
	}, func() {
		/* No-op */
	}

	cls = func() {
		defer stdErrCls()
		defer stdOutCls()
	}

	stdoutFile, err := os.OpenFile(s.Path()+"."+t+".stdout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		s.l.With("err", err).Error("failed to open stdout " + t + "log file")
		return nil, nil, cls, err
	}
	stdOutCls = func() {
		if err := stdoutFile.Close(); err != nil {
			s.l.With("err", err).Error("failed to close stdout " + t + "log file")
		}
	}

	stderr, stdErrCls, err = s.pipeDescriptorStderr(t)
	if err != nil {
		return nil, nil, cls, err
	}

	return stdoutFile, stderr, cls, nil
}

func (s *Solution) pipeDescriptorStderr(t string) (stderr io.Writer, cls func(), err error) {
	cls = func() {
		/* No-op */
	}

	stderrFile, err := os.OpenFile(s.Path()+"."+t+".stderr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		s.l.With("err", err).Error("failed to open stderr " + t + "log file")
		return nil, cls, err
	}
	cls = func() {
		if err := stderrFile.Close(); err != nil {
			s.l.With("err", err).Error("failed to close stderr " + t + "log file")
		}
	}

	return stderrFile, cls, nil
}

func (s *Solution) saveAnswer(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if s.Status.IsBefore(SolutionStatusReview) {
		return r.AddError(ErrNoInputRun, "status before save", s.Status.String())
	} else {
		r.AddRow(result.EmojiCheckMark, "status before save", "", s.Status.String())
	}

	if s.RunResult == nil {
		return r.AddError(ErrNoRunResult, "get answer for save")
	} else if len(s.RunResult.Answer) == 0 {
		return r.AddError(ErrNoAnswer, "get answer for save")
	} else {
		r.AddRow(result.EmojiCheckMark, "get answer for save")
	}

	if err := os.WriteFile(s.Path()+".answer.txt", []byte(s.RunResult.Answer), 0660); err != nil {
		return r.AddError(err, "save answer")
	} else {
		r.AddRow(result.EmojiCheckMark, "save answer")
	}

	return r
}

func (s *Solution) loadAnswer(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if s.Status.IsBefore(SolutionStatusReview) {
		return r.AddError(ErrNoInputRun, "status before load", s.Status.String())
	} else {
		r.AddRow(result.EmojiCheckMark, "status before load", "", s.Status.String())
	}

	if s.RunResult == nil {
		return r.AddError(ErrNoRunResult, "get run result")
	} else {
		r.AddRow(result.EmojiCheckMark, "get run result")
	}

	afp := s.Path()+".answer.txt"
	afpRel := strings.TrimLeft(strings.TrimPrefix(afp, s.part.Path()), "/")

	if b, err := os.ReadFile(afp); err != nil && !os.IsNotExist(err) {
		r.AddError(err, "load answer", afpRel)
	} else if err == nil {
		s.RunResult.Answer = string(b)
		r.AddRow(result.EmojiCheckMark, "load answer", "", afpRel)
	} else {
		r.AddRow(result.EmojiEmpty, "no answer", "", afpRel)
	}

	return r
}

func (s *Solution) loadExpected(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if s.RunResult == nil {
		return r.AddError(ErrNoRunResult, "get run result")
	} else {
		r.AddRow(result.EmojiCheckMark, "get run result")
	}

	efp := path.Join(path.Dir(s.Path()), "answer.txt")
	efpRel := strings.TrimLeft(strings.TrimPrefix(efp, s.part.Path()), "/")

	if expB, err := os.ReadFile(efp); err != nil && !os.IsNotExist(err) {
		r.AddError(err, "load previous answer", efpRel)
	} else if err == nil {
		s.RunResult.Expected = string(expB)
		r.AddRow(result.EmojiCheckMark, "load previous answer", "", efpRel)
	} else {
		r.AddRow(result.EmojiEmpty, "no previous answer", "", efpRel)
	}

	return r
}

func (s *Solution) NrStr() string {
	return strconv.Itoa(s.Nr)
}

func (s *Solution) UnmarshalJSON(b []byte) error {
	var v map[string]json.RawMessage
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v["nr"], &s.Nr); err != nil {
		return err
	}
	if err := json.Unmarshal(v["run_result"], &s.RunResult); err != nil {
		return err
	} else if s.RunResult != nil && s.RunResult.ErrorStr != nil {
		s.RunResult.Error = errors.New(*s.RunResult.ErrorStr)
	}
	if err := json.Unmarshal(v["sample_run_results"], &s.SampleRunResults); err != nil {
		return err
	} else if s.SampleRunResults != nil {
		for _, srr := range s.SampleRunResults {
			if srr.ErrorStr != nil {
				srr.Error = errors.New(*srr.ErrorStr)
			}
		}
	}

	var status string
	if err := json.Unmarshal(v["status"], &status); err != nil {
		return err
	}
	s.Status = solutionStatus(status)

	return nil
}

func (s *Solution) copySolutionCode(r result.Result) (result.Result, string) {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	uniquePackage := s.UniquePackageName()

	base := s.part.Path()
	trgt := s.Path()+"tmp"

	if err := plugin2.CopyGoPackage(&plugin2.DupConfig{
        Base:        base,
        PackageName: uniquePackage,
        TargetDir:   trgt,
    }); err != nil {
		r.AddError(err, "copy solution code", "", uniquePackage)
	}

	return r, trgt
}

func (s *Solution) UniquePackageName() string {
	return fmt.Sprintf("d%2dp%2ds%2dsolution", s.part.Puzzle.Day.Day(), s.part.Nr, s.Nr)
}

func solutionWrapper(fn SolutionFn, ctx context.Context, l *logger.IterationLogger, in io.Reader) (runResult *RunResult, panerr error) {
	defer func() {
		if r := recover(); r != nil {
			panerr = fmt.Errorf("%v", r)
		}
	}()

	out := &strings.Builder{}

	start := time.Now()
	err := fn(ctx, l, in, out)
	runtime := time.Since(start)
	finishedAt := Timestamp(start.Add(runtime))

	return &RunResult{
		Answer:     out.String(),
		Runtime:    &runtime,
		FinishedAt: &finishedAt,
	}, err
}

type solutionStatus string
type SolutionStatus interface {
	IsAfter(ss2 SolutionStatus) bool
	IsBefore(ss2 SolutionStatus) bool
	String() string
}

var _ SolutionStatus = solutionStatus("")

const (
	SolutionStatusCreated        = solutionStatus("created")
	SolutionStatusBuildFailed    = solutionStatus("build failed")
	SolutionStatusBuild          = solutionStatus("build")
	SolutionStatusSamplesInvalid = solutionStatus("invalid samples result")
	SolutionStatusSamplesValid   = solutionStatus("valid samples result")
	SolutionStatusReview         = solutionStatus("result in review")
	SolutionStatusInvalid        = solutionStatus("invalid result")
	SolutionStatusValid          = solutionStatus("valid result")
)

var solutionStatusOrder = []solutionStatus{
	SolutionStatusCreated,
	SolutionStatusBuildFailed,
	SolutionStatusBuild,
	SolutionStatusSamplesInvalid,
	SolutionStatusSamplesValid,
	SolutionStatusReview,
	SolutionStatusInvalid,
	SolutionStatusValid,
}

func (ss solutionStatus) IsAfter(ss2 SolutionStatus) bool {
	if ss == ss2 {
		return false
	}

	seenSS2 := false
	if ss2ss, ok := ss2.(solutionStatus); ok && ss2ss == "" {
		seenSS2 = true
	}

	for _, sso := range solutionStatusOrder {
		if sso == ss2 {
			seenSS2 = true
			continue
		}

		if sso == ss {
			return seenSS2
		}
	}

	panic(fmt.Errorf("invalid Solutions status, not in isafter order list: [%v][%v]", ss, ss2))
}

func (ss solutionStatus) IsBefore(ss2 SolutionStatus) bool {
	if ss == ss2 {
		return false
	}

	if ss2ss, ok := ss2.(solutionStatus); ok && ss2ss == "" {
		return false
	}

	seenSS2 := false

	for _, sso := range solutionStatusOrder {
		if sso == ss {
			return !seenSS2
		}

		if sso == ss2 {
			seenSS2 = true
		}
	}

	panic(fmt.Errorf("invalid Solutions status, not in isbefore order list: [%v][%v]", ss, ss2))
}

func (ss solutionStatus) String() string {
	return string(ss)
}
