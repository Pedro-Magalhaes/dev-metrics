package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type ScanResult struct {
	Processed int
	Skipped   int
}

type JSONLLineError struct {
	Line int
	Err  error
	Raw  string
}

func (e JSONLLineError) Error() string {
	snippet := strings.TrimSpace(e.Raw)
	const max = 200
	if len(snippet) > max {
		snippet = snippet[:max] + "…"
	}
	return fmt.Sprintf("jsonl parse error on line %d: %v (raw=%q)", e.Line, e.Err, snippet)
}

func (e JSONLLineError) Unwrap() error { return e.Err }

// ScanJSONL reads newline-delimited JSON (JSONL) from r.
//
// If strict is true, it stops at the first invalid JSON line and returns an error.
// If strict is false, invalid JSON lines are skipped (counted in ScanResult.Skipped).
func ScanJSONL(r io.Reader, strict bool, fn func(BuildMetric) error) (ScanResult, error) {
	var res ScanResult

	scanner := bufio.NewScanner(r)
	// Default scanner token limit (64K) can be too small for very long JSON lines.
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var m BuildMetric
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			res.Skipped++
			if strict {
				return res, JSONLLineError{Line: lineNo, Err: err, Raw: line}
			}
			continue
		}

		if err := fn(m); err != nil {
			return res, fmt.Errorf("processing jsonl line %d: %w", lineNo, err)
		}
		res.Processed++
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}
