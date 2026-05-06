// Package exporter provides functionality for writing filtered log lines
// to various output destinations in configurable formats.
//
// Basic usage:
//
//	e, err := exporter.New(exporter.Options{
//		Format:     exporter.FormatRaw,
//		OutputPath: "/tmp/slice.log", // omit for stdout
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer e.Close()
//
//	if err := e.Export(filteredLines); err != nil {
//		log.Fatal(err)
//	}
//
// Supported formats:
//
//	- FormatRaw      — writes each line verbatim followed by a newline.
//	- FormatNumbered — prefixes each line with its 1-based index.
package exporter
