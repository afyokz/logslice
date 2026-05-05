# logslice

Fast log file slicer that filters and exports time-range segments from large log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

```bash
logslice [flags] <logfile>
```

### Examples

Extract logs between two timestamps:
```bash
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" app.log
```

Output to a file:
```bash
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" --out slice.log app.log
```

Filter by log level within a time range:
```bash
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" --level ERROR app.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Start of time range (RFC3339 or common log formats) | required |
| `--to` | End of time range | required |
| `--out` | Output file path | stdout |
| `--level` | Filter by log level (INFO, WARN, ERROR) | all |
| `--format` | Timestamp format hint | auto-detect |

## Building from Source

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build ./...
```

## License

MIT — see [LICENSE](LICENSE) for details.