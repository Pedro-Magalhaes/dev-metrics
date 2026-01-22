package metrics

import (
	"fmt"
	"io"
	"sort"
	"time"
)

type reportKey struct {
	Project string
	Year    int
	Week    int
}

// GenerateReport processa o log e retorna os dados estruturados
func GenerateReport(r io.Reader) (*FullReport, error) {
	// 1. Estruturas temporárias para acumulação (Mapas)
	// Map: [Projeto - Semana - ano] -> Stats
	tempData := make(map[reportKey]*BuildStats)

	// 2. Scan e Acumulação
	_, err := ScanJSONL(r, false, func(m BuildMetric) error {
		t, err := time.Parse(time.RFC3339, m.Timestamp)
		if err != nil {
			return nil // Ignora erro de parse pontual ou retorna erro
		}

		year, week := t.ISOWeek()

		key := reportKey{
			Project: m.Project,
			Year:    year,
			Week:    week,
		}

		if _, ok := tempData[key]; !ok {
			tempData[key] = &BuildStats{}
		}

		stats := tempData[key]
		stats.TotalDuration += m.DurationSec
		stats.Count++
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Transformação de Mapas para Slices (Struct Final)
	report := &FullReport{}
	projectMap := make(map[string]*ProjectSummary)

	for k, stat := range tempData {
		// Busca ou cria o projeto no mapa auxiliar
		if _, ok := projectMap[k.Project]; !ok {
			projectMap[k.Project] = &ProjectSummary{Name: k.Project}
		}
		proj := projectMap[k.Project]

		weekLabel := fmt.Sprintf("%d-W%02d", k.Year, k.Week)

		summary := WeeklySummary{
			WeekLabel:   weekLabel,
			BuildStats:  *stat,
			AvgDuration: 0,
		}
		if stat.Count > 0 {
			summary.AvgDuration = stat.TotalDuration / float64(stat.Count)
		}

		proj.Weeks = append(proj.Weeks, summary)
		proj.TotalDuration += stat.TotalDuration
		proj.TotalBuilds += stat.Count
	}

	// 3. Transformar mapa auxiliar em Slice final e Ordenar
	for _, proj := range projectMap {
		// Ordena semanas
		sort.Slice(proj.Weeks, func(i, j int) bool {
			return proj.Weeks[i].WeekLabel < proj.Weeks[j].WeekLabel
		})
		report.Projects = append(report.Projects, *proj)
		report.GlobalDuration += proj.TotalDuration
		report.GlobalBuilds += proj.TotalBuilds
	}

	// Ordena projetos
	sort.Slice(report.Projects, func(i, j int) bool {
		return report.Projects[i].Name < report.Projects[j].Name
	})

	return report, nil
}
