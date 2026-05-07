// Package followpipeline wires a TailReader into the logslice processing
// stack so that live log lines are filtered, matched, and exported in
// follow (tail -f) mode.
package followpipeline

import (
	"context"
	"time"

	"github.com/user/logslice/internal/exporter"
	"github.com/user/logslice/internal/matcher"
	"github.com/user/logslice/internal/tailreader"
	"github.com/user/logslice/internal/timeparser"
)

// Config holds the options for a follow pipeline run.
type Config struct {
	FilePath     string
	PollInterval time.Duration
	From         string // optional lower-bound timestamp
	Include      []string
	Exclude      []string
	Format       string // timestamp layout for parsing
	ExportCfg    exporter.Config
}

// Run starts tailing FilePath and processes each line through the matcher
// and exporter until ctx is cancelled.
func Run(ctx context.Context, cfg Config) error {
	tr, err := tailreader.New(cfg.FilePath, cfg.PollInterval)
	if err != nil {
		return err
	}
	defer tr.Stop()

	m, err := matcher.New(cfg.Include, cfg.Exclude)
	if err != nil {
		return err
	}

	exp, err := exporter.New(cfg.ExportCfg)
	if err != nil {
		return err
	}
	defer exp.Close()

	var fromTime time.Time
	if cfg.From != "" {
		fromTime, err = timeparser.ParseTimestamp(cfg.From)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-tr.Lines():
			if !ok {
				return nil
			}
			if !fromTime.IsZero() {
				ts, perr := timeparser.ParseWithFormat(line, cfg.Format)
				if perr == nil && ts.Before(fromTime) {
					continue
				}
			}
			if !m.Match(line) {
				continue
			}
			if werr := exp.Write(line); werr != nil {
				return werr
			}
		}
	}
}
