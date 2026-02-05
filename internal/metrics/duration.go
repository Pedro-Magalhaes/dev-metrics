package metrics

import (
	"fmt"
	"math"
	"strings"
)

// DurationUnit define a unidade para exibição de durações.
type DurationUnit string

const (
	DurationAuto    DurationUnit = "auto"
	DurationSeconds DurationUnit = "s"
	DurationMinutes DurationUnit = "min"
	DurationHours   DurationUnit = "h"
)

// ParseDurationUnit converte a string da flag para DurationUnit.
func ParseDurationUnit(value string) (DurationUnit, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "", "auto":
		return DurationAuto, nil
	case "s", "sec", "seg", "secs", "second", "seconds", "segundos":
		return DurationSeconds, nil
	case "m", "min", "mins", "minute", "minutes", "minutos":
		return DurationMinutes, nil
	case "h", "hr", "hrs", "hour", "hours", "horas":
		return DurationHours, nil
	default:
		return "", fmt.Errorf("unidade inválida: %s (use auto|s|min|h)", value)
	}
}

// DurationUnitLabel retorna o rótulo padrão para a unidade.
func DurationUnitLabel(unit DurationUnit) string {
	switch unit {
	case DurationSeconds:
		return "s"
	case DurationMinutes:
		return "min"
	case DurationHours:
		return "h"
	default:
		return ""
	}
}

// AutoDurationUnit escolhe a unidade baseada no limiar definido.
func AutoDurationUnit(seconds float64) DurationUnit {
	if seconds >= 3600 {
		return DurationHours
	}
	if seconds >= 60 {
		return DurationMinutes
	}
	return DurationSeconds
}

// FormatDuration formata um valor em segundos para a unidade desejada.
// includeUnit controla se o sufixo de unidade é adicionado ao texto.
func FormatDuration(seconds float64, unit DurationUnit, includeUnit bool) string {
	u := unit
	if u == DurationAuto {
		u = AutoDurationUnit(seconds)
	}

	switch u {
	case DurationMinutes:
		totalSeconds := int(math.Round(seconds))
		mins := totalSeconds / 60
		secs := totalSeconds % 60
		return fmt.Sprintf("%dmin%02ds", mins, secs)
	case DurationHours:
		totalMinutes := int(math.Round(seconds / 60))
		hours := totalMinutes / 60
		mins := totalMinutes % 60
		return fmt.Sprintf("%dh%02dmin", hours, mins)
	case DurationSeconds:
		formatted := fmt.Sprintf("%.1f", seconds)
		if includeUnit {
			return fmt.Sprintf("%s s", formatted)
		}
		return formatted
	}

	formatted := fmt.Sprintf("%.1f", seconds)
	if includeUnit {
		return fmt.Sprintf("%s %s", formatted, DurationUnitLabel(u))
	}
	return formatted
}
