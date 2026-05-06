package pipeline

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/logslice/internal/cli"
	"github.com/yourorg/logslice/internal/deduplicator"
	"github.com/yourorg/logslice/internal/exporter"
	"github.com/yourorg/logslice/internal/limiter"
	"github.com/yourorg/logslice/internal/matcher"
	"github.com/yourorg/logslice/internal/offsetter"
	"github.com/yourorg/logslice/internal/progress"
	"github.com/yourorg/logslice/internal/sampler"
	"github.com/yourorg/logslice/internal/scanner"
)

// Run executes the full log-slicing pipeline according to cfg.
func Run(cfg *cli.Config, out io.Writer) error {
	f, err := os.Open(cfg.Input)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer f.Close()

	off := offsetter.New(cfg.TZOffset)
	from := off.ShiftFrom(cfg.From)
	to := off.ShiftFrom(cfg.To)

	m, err := matcher.New(cfg.Include, cfg.Exclude)
	if err != nil {
		return fmt.Errorf("matcher: %w", err)
	}

	sc := scanner.New(bufio.NewReader(f), from, to)
	exp := exporter.New(out, cfg.Output, cfg.Numbered)
	prog := progress.New(cfg.Verbose, cfg.Total, os.Stderr)
	smp := sampler.New(cfg.Step)
	lim := limiter.New(cfg.MaxLines)
	dedup := deduplicator.New(cfg.DeduplicateWindow)

	var lines []string
	for sc.Scan() {
		prog.IncProcessed()
		line := sc.Text()

		if dedup.IsDuplicate(line) {
			prog.IncSkipped()
			continue
		}

		if !m.Match(line) {
			prog.IncSkipped()
			continue
		}

		if !smp.Keep() {
			prog.IncSkipped()
			continue
		}

		if lim.Reached() {
			break
		}

		lim.Allow()
		prog.IncMatched()
		lines = append(lines, line)
	}

	if err := sc.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	prog.Print()
	return exp.Export(lines)
}
