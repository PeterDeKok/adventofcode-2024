package remote

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/utils/testabletime"
	"sync"
	"time"
)

const (
	RateInput   string = "input.txt"
	RateProblem string = "PROBLEM.md"
	RateFun     string = "fun.txt"
)

type RateLimit struct {
	cnf *RateLimitConfig
	l   *log.Logger

	Last   map[string]time.Time     `json:"last"`
	Limits map[string]time.Duration `json:"limits"`

	sync.Mutex
}

type RateLimitConfig struct {
	Logger     *log.Logger
	PuzzlesDir string
}

var defaultLimits = map[string]time.Duration{
	RateInput:   5 * time.Minute,
	RateProblem: 5 * time.Minute,
	RateFun:     5 * time.Minute,
}

func NewRateLimit(cnf *RateLimitConfig) *RateLimit {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "remote:ratelimit")
	} else {
		l = l.With("type", "remote:ratelimit")
	}

	rl := &RateLimit{
		cnf: cnf,
		l:   l,

		Last:   make(map[string]time.Time),
		Limits: make(map[string]time.Duration),
	}

	rl.load()

	return rl
}

func (r *RateLimit) Try(cat string) (bool, time.Duration, error) {
	r.Lock()
	defer r.Unlock()

	r.load()

	last, ok := r.Last[cat]
	if !ok {
		last = time.Time{}
		r.Last[cat] = last
	}

	limit, ok := r.Limits[cat]
	if !ok {
		return false, 0, fmt.Errorf("rate limit category %s invalid", cat)
	}

	now := testabletime.Now()

	safe := last.Add(limit)
	until := safe.Sub(now)

	if safe.After(now) {
		return false, until, nil
	}

	r.Last[cat] = now

	return true, until, r.save()
}

func (r *RateLimit) load() {
	fp := path.Join(r.cnf.PuzzlesDir, "ratelimit.json")
	if ok, err := tools.FileExists(fp); err != nil {
		r.l.Error("failed to test if ratelimit.json exists", "err", err)
	} else if !ok {
		r.l.Info("reatelimit.json does not exist yet")
	} else if f, err := os.ReadFile(fp); err != nil {
		r.l.Error("failed to read ratelimit.json contents", "err", err)
	} else if err := json.Unmarshal(f, r); err != nil {
		r.l.Error("failed to parse ratelimit.json contents", "err", err)
	}

	// Set defaults if category not loaded
	for k, limit := range defaultLimits {
		if v, ok := r.Limits[k]; !ok {
			r.Limits[k] = limit
		} else if v != limit {
			r.l.Warn("default rate limit different", "category", k, "default", limit.String(), "json", v.String())
		}
	}
}

func (r *RateLimit) save() error {
	fp := path.Join(r.cnf.PuzzlesDir, "ratelimit.json")

	if b, err := json.Marshal(r); err != nil {
		r.l.Error("failed to marshal ratelimit.json contents", "err", err)
		return err
	} else if err = os.WriteFile(fp, b, 0660); err != nil {
		r.l.Error("failed to save ratelimit.json contents", "err", err)
		return err
	}

	return nil
}
