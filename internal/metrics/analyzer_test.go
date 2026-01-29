package metrics

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGenerateReport(t *testing.T) {
	// Definição dos cenários de teste
	tests := []struct {
		name    string
		input   string // Simula o conteúdo do arquivo JSONL
		options ReportOptions
		want    *FullReport
		wantErr bool
	}{
		{
			name: "Basic Aggregation",
			input: `
{"project": "backend", "timestamp": "2024-01-03T10:00:00Z", "duration_sec": 10}
{"project": "backend", "timestamp": "2024-01-04T12:00:00Z", "duration_sec": 20}
`,
			// 2024-01-03 é Semana 01 de 2024
			want: &FullReport{
				Projects: []ProjectSummary{
					{
						Name:          "backend",
						TotalDuration: 30,
						TotalBuilds:   2,
						Weeks: []WeeklySummary{
							{
								WeekLabel: "2024-W01",
								BuildStats: BuildStats{
									TotalDuration: 30,
									Count:         2,
								},
								AvgDuration: 15, // (10+20)/2
							},
						},
					},
				},
				GlobalDuration: 30,
				GlobalBuilds:   2,
			},
			wantErr: false,
		},
		{
			name:    "Basic Aggregation with until filter",
			options: ReportOptions{Until: parseTime(t, "2024-01-04T00:00:00Z")},
			input: `
{"project": "backend", "timestamp": "2024-01-03T10:00:00Z", "duration_sec": 10}
{"project": "backend", "timestamp": "2024-01-04T00:00:00Z", "duration_sec": 20}
{"project": "backend", "timestamp": "2024-01-04T12:00:00Z", "duration_sec": 200}
{"project": "backend", "timestamp": "2024-01-05T12:00:00Z", "duration_sec": 200}
`,
			// 2024-01-03 é Semana 01 de 2024
			want: &FullReport{
				Projects: []ProjectSummary{
					{
						Name:          "backend",
						TotalDuration: 30,
						TotalBuilds:   2,
						Weeks: []WeeklySummary{
							{
								WeekLabel: "2024-W01",
								BuildStats: BuildStats{
									TotalDuration: 30,
									Count:         2,
								},
								AvgDuration: 15, // (10+20)/2
							},
						},
					},
				},
				GlobalDuration: 30,
				GlobalBuilds:   2,
			},
			wantErr: false,
		},
		{
			name: "Sorting: Projects and Weeks",
			input: `
{"project": "zeta-service", "timestamp": "2024-01-01T00:00:00Z", "duration_sec": 5}
{"project": "alpha-service", "timestamp": "2024-01-08T00:00:00Z", "duration_sec": 10}
{"project": "alpha-service", "timestamp": "2024-01-01T00:00:00Z", "duration_sec": 10}
`,
			// Esperado: alpha-service antes de zeta-service
			// Esperado (alpha): Semana 01 antes da Semana 02
			want: &FullReport{
				Projects: []ProjectSummary{
					{
						Name:          "alpha-service",
						TotalDuration: 20,
						TotalBuilds:   2,
						Weeks: []WeeklySummary{
							{
								WeekLabel:   "2024-W01",
								BuildStats:  BuildStats{TotalDuration: 10, Count: 1},
								AvgDuration: 10,
							},
							{
								WeekLabel:   "2024-W02",
								BuildStats:  BuildStats{TotalDuration: 10, Count: 1},
								AvgDuration: 10,
							},
						},
					},
					{
						Name:          "zeta-service",
						TotalDuration: 5,
						TotalBuilds:   1,
						Weeks: []WeeklySummary{
							{
								WeekLabel:   "2024-W01",
								BuildStats:  BuildStats{TotalDuration: 5, Count: 1},
								AvgDuration: 5,
							},
						},
					},
				},
				GlobalDuration: 25,
				GlobalBuilds:   3,
			},
			wantErr: false,
		},
		{
			name:    "Empty Input",
			input:   "",
			want:    &FullReport{}, // Retorna struct vazia, mas inicializada
			wantErr: false,
		},
		{
			name: "Invalid Timestamp (Should Skip)",
			input: `
{"project": "A", "timestamp": "DATA-INVALIDA", "duration_sec": 10}
{"project": "A", "timestamp": "2024-01-01T10:00:00Z", "duration_sec": 5}
`,
			// Deve ignorar a primeira linha e processar a segunda
			want: &FullReport{
				Projects: []ProjectSummary{
					{
						Name:          "A",
						TotalDuration: 5,
						TotalBuilds:   1,
						Weeks: []WeeklySummary{
							{
								WeekLabel:   "2024-W01",
								BuildStats:  BuildStats{TotalDuration: 5, Count: 1},
								AvgDuration: 5,
							},
						},
					},
				},
				GlobalDuration: 5,
				GlobalBuilds:   1,
			},
			wantErr: false,
		},
		{
			name: "Broken JSON lines should be skipped",
			input: `
{"project": "... ERROR
{"project": "A", "time... ERROR
{"project": "backend", "timestamp": "2024-01-03T10:00:00Z", "duration_sec": 10}
{"project": "backend", "timestamp": "2024-01-04T12:00:00Z", "duration_sec": 20}
`,
			want: &FullReport{
				Projects: []ProjectSummary{
					{
						Name:          "backend",
						TotalDuration: 30,
						TotalBuilds:   2,
						Weeks: []WeeklySummary{
							{
								WeekLabel: "2024-W01",
								BuildStats: BuildStats{
									TotalDuration: 30,
									Count:         2,
								},
								AvgDuration: 15, // (10+20)/2
							},
						},
					},
				},
				GlobalDuration: 30,
				GlobalBuilds:   2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cria um Reader a partir da string de input
			r := strings.NewReader(tt.input)

			got, err := GenerateReport(r, tt.options)

			// Verifica se o erro ocorreu conforme esperado
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Se esperamos erro, não precisamos comparar o retorno
			if tt.wantErr {
				return
			}

			// Inicializa slices vazios no 'got' para evitar erro de comparação com nil vs []
			// (O DeepEqual diferencia slice nil de slice vazio)
			if got.Projects == nil {
				got.Projects = nil
			}

			// Comparação profunda das estruturas
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateReport() = \n%+v, \nwant \n%+v", got, tt.want)
			}
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestGenerateReportWithInvalidReader(t *testing.T) {
	// Simula um reader que retorna erro
	r := &errorReader{}

	_, err := GenerateReport(r, ReportOptions{})
	if err == nil {
		t.Errorf("GenerateReport() expected error, got nil")
	}
}

// Funções auxiliares de teste

// transforma data em string para time.Time
func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Fatalf("failed to parse time %s: %v", timeStr, err)
	}
	return parsed
}
