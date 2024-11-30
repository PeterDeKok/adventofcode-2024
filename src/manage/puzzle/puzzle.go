package puzzle

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"math"
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/color"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/styles"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/utils/testabletime"
	"time"
)

type Puzzle struct {
	cnf *Config
	l   *log.Logger

	Day   time.Time
	Error error

	Part1 *Part
	Part2 *Part
}

type Config struct {
	Ctx        context.Context
	Logger     *log.Logger
	Remote     *remote.Client
	PuzzlesDir string
	Day        time.Time
}

func Create(cnf *Config) *Puzzle {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "Puzzle", "day", cnf.Day.Format(time.DateOnly))
	} else {
		l = l.With("type", "Puzzle", "day", cnf.Day.Format(time.DateOnly))
	}

	p := &Puzzle{
		cnf: cnf,
		l:   l,

		Day: cnf.Day,
	}

	p.Part1 = CreatePart(&PartConfig{
		Ctx:    cnf.Ctx,
		Logger: p.l,
		Remote: cnf.Remote,
		Puzzle: p,
		Nr:     1,
		SkipValidationOnCreate: true,
	})
	p.Part2 = CreatePart(&PartConfig{
		Ctx:    cnf.Ctx,
		Logger: p.l,
		Remote: cnf.Remote,
		Puzzle: p,
		Nr:     2,
		SkipValidationOnCreate: true,
	})

	return p
}

func (pz *Puzzle) FilterValue() string {
	return pz.Day.Format(time.DateOnly)
}

func (pz *Puzzle) Title() string {
	return pz.Day.Format(time.DateOnly)
}

func (pz *Puzzle) Description() string {
	if pz.Error != nil {
		return pz.Error.Error()
	} else if pz.Part1.Error != nil {
		return pz.Part1.Error.Error()
	} else if pz.Part2.Error != nil {
		return pz.Part2.Error.Error()
	}

	if pz.Day.After(testabletime.Now()) {
		return pz.until()
	}

	str := ""
	str += pz.Part1.RenderStar()
	str += pz.Part2.RenderStar()
	str += "    "

	if pz.Part2.OK() && pz.Part2.FastestSolution != nil && pz.Part2.FastestSolution.RunResult != nil {
		str += lipgloss.NewStyle().Foreground(styles.DimmedColor).Render("p2: " + pz.Part2.FastestSolution.RunResult.Runtime.String())
	} else if pz.Part1.OK() && pz.Part1.FastestSolution != nil && pz.Part1.FastestSolution.RunResult != nil {
		str += lipgloss.NewStyle().Foreground(styles.DimmedColor).Render("p1: " + pz.Part1.FastestSolution.RunResult.Runtime.String())
	}

	return str
}

func (pz *Puzzle) Path() string {
	return path.Join(
		pz.cnf.PuzzlesDir,
		pz.Day.Format(time.DateOnly),
	)
}

func (pz *Puzzle) until() string {
	until := pz.Day.Sub(testabletime.Now())
	h := int(until.Hours())

	if h > 48 {
		return pz.Day.Format(time.DateOnly)
	}

	m := int(math.Mod(until.Minutes(), 60))
	s := int(math.Mod(until.Seconds(), 60))

	if h > 23 {
		return fmt.Sprintf("%dh %dm", h, m)
	}

	if h > 1 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}

	if m > 1 {
		return fmt.Sprintf("%dm %ds", m, s)
	}

	return fmt.Sprintf("%s%ds%s !!", color.Green, s, color.Reset)
}

func (pz *Puzzle) ValidateDir(r result.Result) result.Result {
	r = result.OrNew(r, 1)
	defer r.Increment(1)

	if _, err := ValidatePuzzleDir(pz.Path()); err != nil {
		pz.Error = err
		r.AddError(err, "validate Puzzle dir")
		return r
	}

	return r
}

func ValidatePuzzleDir(dir string) (bool, error) {
	if ok, err := tools.DirExists(dir); !ok || err != nil {
		return ok, err
	}

	return true, nil
}
