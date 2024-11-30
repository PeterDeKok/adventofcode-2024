package remote

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	cnf *ClientConfig
	l   *log.Logger
	c   *http.Client

	RateLimit *RateLimit
}

type ClientConfig struct {
	Logger               *log.Logger
	PuzzlesDir           string
	Base                 string
	SessionCookieValue   string
	SessionCookieExpires time.Time
}

type RequestOptions struct {
	DisableCache      bool
	RateLimitCategory string
}

func CreateClient(cnf *ClientConfig) *Client {
	l := cnf.Logger

	if l == nil {
		l = log.With("type", "remote:ratelimit")
	} else {
		l = l.With("type", "remote:ratelimit")
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic("failed to create client: failed to create cookiejar")
	}

	website, err := url.Parse(cnf.Base)
	if err != nil {
		panic("failed to create client: failed to parse base url")
	}

	jar.SetCookies(website, []*http.Cookie{
		{
			Name:     "session",
			Value:    cnf.SessionCookieValue,
			Domain:   ".adventofcode.com",
			Path:     "/",
			Expires:  cnf.SessionCookieExpires,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		},
	})

	c := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	return &Client{
		cnf: cnf,
		l:   l,
		RateLimit: NewRateLimit(&RateLimitConfig{
			Logger:     cnf.Logger,
			PuzzlesDir: cnf.PuzzlesDir,
		}),
		c: c,
	}
}

func (c *Client) Get(ctx context.Context, path string, options *RequestOptions) (io.ReadCloser, error) {
	l := c.l.With("method", "GET", "path", path)

	if options == nil || len(options.RateLimitCategory) == 0 {
		l.Error("missing rate limit category")
		return nil, fmt.Errorf("failed to create GET request: missing rate limit category")
	}

	if ok, until, err := c.RateLimit.Try(options.RateLimitCategory); err != nil {
		l.Error("rate limit error", "cat", options.RateLimitCategory, "err", err)
		return nil, fmt.Errorf("failed to create GET request: rate limit error: %v", err)
	} else if !ok {
		l.Error("rate limit blocking", "cat", options.RateLimitCategory, "until", until)
		return nil, fmt.Errorf("failed to create GET request: rate limit blocking")
	} else {
		l.Info("rate limit safe", "cat", options.RateLimitCategory, "margin", until)
	}

	u, err := c.url(path)
	if err != nil {
		l.Error("url invalid")
		return nil, fmt.Errorf("failed to create GET request: url invalid: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		l.Error("failed to create request")
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	req.Header.Add("Accept", "text/plain,text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Referer", "https://github.com/PeterDeKok/adventofcode-2024")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Priority", "u=6")

	l.Info("about to perform request")

	if resp, err := c.c.Do(req); err != nil {
		l.Error("failed to process request")
		return nil, fmt.Errorf("failed to process GET request: %v", err)
	} else {
		l.Info("response returned", "responseheaders", resp.Header, "status", resp.StatusCode, "requestheaders", resp.Request.Header)
		return resp.Body, nil
	}
}

func (c *Client) url(path string) (string, error) {
	return url.JoinPath(c.cnf.Base, path)
}
