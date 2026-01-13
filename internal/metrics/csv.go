package metrics

import (
	"encoding/csv"
	"io"
	"strconv"
)

func CSVHeader() []string {
	return []string{
		"timestamp",
		"user",
		"hostname",
		"os",
		"project",
		"branch",
		"commit",
		"command",
		"duration_sec",
		"returncode",
		"cpus",
		"status",
	}
}

func BuildMetricCSVRow(m BuildMetric) []string {
	return []string{
		m.Timestamp,
		m.User,
		m.Hostname,
		m.OS,
		m.Project,
		m.Branch,
		m.Commit,
		m.Command,
		strconv.FormatFloat(m.DurationSec, 'f', -1, 64),
		strconv.Itoa(m.ReturnCode),
		strconv.Itoa(m.CPUs),
		m.Status,
	}
}

// ExportCSVFromJSONL converts a JSONL stream to CSV.
//
// The CSV header is always written as the first row.
func ExportCSVFromJSONL(r io.Reader, w io.Writer, strict bool) (ScanResult, error) {
	csvw := csv.NewWriter(w)
	if err := csvw.Write(CSVHeader()); err != nil {
		return ScanResult{}, err
	}

	res, err := ScanJSONL(r, strict, func(m BuildMetric) error {
		return csvw.Write(BuildMetricCSVRow(m))
	})
	csvw.Flush()
	if werr := csvw.Error(); werr != nil {
		return res, werr
	}
	if err != nil {
		return res, err
	}
	return res, nil
}
