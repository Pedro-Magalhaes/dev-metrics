package metrics

// BuildMetric representa os dados coletados de uma execução de build
type BuildMetric struct {
	Timestamp   string  `json:"timestamp"`
	User        string  `json:"user"`
	Hostname    string  `json:"hostname"`
	OS          string  `json:"os"`
	Project      string  `json:"project"`
	Branch      string  `json:"branch"`
	Commit      string  `json:"commit"`
	Command     string  `json:"command"`
	DurationSec float64 `json:"duration_sec"`
	ReturnCode  int     `json:"returncode"`
	CPUs        int     `json:"cpus"`
	Status      string  `json:"status"`
}
