package manage

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/puzzle"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/remote"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/exit"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/utils/testabletime"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	cnf     *AdventOfCodeConfig
	puzzles []*puzzle.Puzzle
	lookup  map[string]*puzzle.Puzzle

	funLines []string
}

type AdventOfCodeConfig struct {
	Ctx    context.Context
	Logger *log.Logger
	Remote *remote.Client

	PuzzlesDir string
	FirstDay   string
	NrDays     int
	TZ         string
}

func Create(cnf *AdventOfCodeConfig) *Manager {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "manager")
	} else {
		l = l.With("type", "manager")
	}

	puzzles := make([]*puzzle.Puzzle, cnf.NrDays)
	lookup := make(map[string]*puzzle.Puzzle, cnf.NrDays)

	for i, d := range DaysGenerator(cnf.FirstDay, cnf.TZ, cnf.NrDays) {
		p := puzzle.Create(&puzzle.Config{
			Ctx:        cnf.Ctx,
			Logger:     l,
			Remote:     cnf.Remote,
			Day:        d,
			PuzzlesDir: cnf.PuzzlesDir,
		})

		puzzles[i] = p
		lookup[d.Format(time.DateOnly)] = p
	}

	return &Manager{
		cnf:     cnf,
		puzzles: puzzles,
		lookup:  lookup,

		funLines: []string{},
	}
}

func (m *Manager) LoadLocal() result.Result {
	r := result.New()
	r.AddTotal(2)
	defer r.Increment(1)

	l := log.With("puzzles-dir", m.cnf.PuzzlesDir)

	dir, err := os.ReadDir(m.cnf.PuzzlesDir)
	if err != nil {
		l.Error("failed to read puzzles dir", "err", err)
		r.AddError(fmt.Errorf("failed to read puzzles dir: %v", err))
		return r
	}
	r.Increment(1)

	r.AddTotal(3 * len(dir))

	for _, f := range dir {
		if !f.IsDir() {
			continue
		}

		p, ok := m.lookup[f.Name()]
		if !ok {
			l.With("name", f.Name()).Warn("directory not recognized")
			continue
		}

		if ok, err := tools.DirExists(f.Name()); err != nil {
			l.Error("load local dir err", "err", err)
			r.AddError(err, "puzzle dir error", f.Name())
		} else if ok {
			r.AddRow(result.EmojiCheckMark, "validate dir", "", f.Name())
			r.Increment(1)

			if ok, err := tools.DirExists(p.Part1.Path()); err != nil {
				l.Error("load local dir part1 err", "err", err)
				r.AddError(err, "part1 dir exists")
			} else if ok {
				p.Part1.ValidateDir(r)
				r.AddRow(result.EmojiCheckMark, "validated part 1", "", f.Name())
				r.Increment(1)
			}

			if ok, err := tools.DirExists(p.Part2.Path()); err != nil {
				l.Error("load local dir part2 err", "err", err)
				r.AddError(err, "part2 dir exists")
			} else if ok {
				p.Part2.ValidateDir(r)
				r.AddRow(result.EmojiCheckMark, "validated part 2", "", f.Name())
				r.Increment(1)
			}
		} else {
			r.Increment(3)
		}
	}

	if build.DEBUG {
		fmt.Println(r)
	}

	m.loadFunLines(r)

	return r
}

// MustLoadLocal loads the local files and stats into memory.
// If loading fails, it will output an error to stdout and exit the program.
func (m *Manager) MustLoadLocal() {
	r := m.LoadLocal()

	if r.Error() != nil {
		fmt.Printf("unable to load the local puzzle files and stats: %v\n", r)

		panic(exit.ErrExitManager)
	}
}

func (m *Manager) Watch() {

}

func (m *Manager) Start() {

}

func (m *Manager) Path() string {
	return m.cnf.PuzzlesDir
}

func (m *Manager) Puzzles() []*puzzle.Puzzle {
	return m.puzzles
}

func (m *Manager) HasFuturePuzzle() bool {
	for _, d := range m.puzzles {
		if d.Day.After(testabletime.Now()) {
			return true
		}
	}

	return false
}

func (m *Manager) NextPuzzle() (index int, next *puzzle.Puzzle, in time.Duration) {
	for i, d := range m.puzzles {
		if !d.Day.After(testabletime.Now()) {
			continue
		}

		if next == nil {
			next = d
			index = i
			continue
		}

		if d.Day.Before(next.Day) {
			next = d
			index = i
		}
	}

	if next != nil {
		in = next.Day.Sub(testabletime.Now())
	}

	return index, next, in
}

func (m *Manager) FunLines() []string {
	return m.funLines
}

func (m *Manager) loadFunLines(r result.Result) result.Result {
	if r == nil {
		r = result.New()
	}
	r.AddTotal(1)

	fp := path.Join(m.Path(), "fun.txt")

	if ok, err := tools.FileExists(fp); err != nil {
		return r.AddError(err, "fun file exists")
	} else if !ok {
		r.AddRow(result.EmojiCheckMark, "fun file does not exist")
		r.Increment(1)
		return r
	} else {
		r.AddRow(result.EmojiCheckMark, "fun file exists")
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		return r.AddError(err, "read fun file")
	} else if len(b) == 0 {
		r.AddRow(result.EmojiCheckMark, "fun file does not contain lines")
		r.Increment(1)
		return r
	} else {
		r.AddRow(result.EmojiCheckMark, "read fun file")
	}

	m.funLines = strings.Split(string(b), "\n")

	r.Increment(1)

	return r
}
func (m *Manager) LoadRemoteFunLines(ctx context.Context) result.Result {
	r := result.New()
	r.SetTotal(1)

	go func(r result.Result) {
		if ok, err := tools.DirExists(m.Path()); err != nil {
			r.AddError(err, "puzzles dir exists")
			return
		} else if !ok {
			r.AddError(os.ErrNotExist, "puzzles dir exists")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "puzzles dir exists")
		}

		if m.cnf.Remote == nil {
			r.AddError(fmt.Errorf("remote not set"), "select remote")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "select remote")
		}

		if len(m.puzzles) == 0 {
			r.AddError(fmt.Errorf("not found"), "first puzzle")
			return
		}
		year := strconv.Itoa(m.puzzles[0].Day.Year())
		r.AddRow(result.EmojiCheckMark, "first puzzle", "", year)

		uri := fmt.Sprintf("%s", year)

		body, err := m.cnf.Remote.Get(ctx, uri, &remote.RequestOptions{
			RateLimitCategory: remote.RateFun,
		})
		if err != nil {
			r.AddError(err, "retrieve content from remote", remote.RateFun)
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "retrieve content from remote", "", remote.RateFun)
		}
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(body)

		lines, err := ParseFunLines(body)
		if err != nil {
			r.AddError(err, "parse fun lines")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "parse fun lines")
		}

		if err := overwriteFile(m.Path(), "fun.txt", []byte(lines)); err != nil {
			r.AddError(err, "write fun lines to file", "fun.txt")
			return
		} else {
			r.AddRow(result.EmojiCheckMark, "write fun lines to file", "", "fun.txt")
		}

		r.Increment(1)

		m.loadFunLines(r)
	}(r)

	return r
}

func overwriteFile(dir, name string, data []byte) error {
	fp := path.Join(dir, name)

	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
