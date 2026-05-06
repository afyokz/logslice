package pipeline

import (
	"bufio"
	"io"

	"github.com/logslice/logslice/internal/exporter"
	"github.com/logslice/logslice/internal/matcher"
	"github.com/logslice/logslice/internal/scanner"
	"github.com/logslice/logslice/internal/timeparser"
)

// Config holds all parameters needed to run a pipeline.
type Config struct {
	Reader    io.Reader
	Writer    io.Writer
	From      string
	To        string
	Includes  []string
	Excludes  []string
	Numbered  bool
	OutputFile string
}

// Result summarises what the pipeline processed.
type Result struct {
	Scanned  int
	Matched  int
	Exported int
}

// Run executes the full slice pipeline: scan → match → export.
func Run(cfg Config) (Result, error) {
	from, err := timeparser.ParseTimestamp(cfg.From)
	if err != nil {
		return Result{}, err
	}
	to, err := timeparser.ParseTimestamp(cfg.To)
	if err != nil {
		return Result{}, err
	}

	m, err := matcher.New(cfg.Includes, cfg.Excludes)
	if err != nil {
		return Result{}, err
	}

	s := scanner.New(bufio.NewReader(cfg.Reader), from, to)
	lines, scanErr := s.Scan()

	var filtered []string
	for _, l := range lines {
		if m.Match(l) {
			filtered = append(filtered, l)
		}
	}

	ex, err := exporter.New(cfg.Writer, cfg.OutputFile, cfg.Numbered)
	if err != nil {
		return Result{}, err
	}
	exported, expErr := ex.Export(filtered)

	if scanErr != nil {
		return Result{}, scanErr
	}
	if expErr != nil {
		return Result{}, expErr
	}

	return Result{
		Scanned:  len(lines),
		Matched:  len(filtered),
		Exported: exported,
	}, nil
}
