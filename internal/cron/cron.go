// Package cron parses and explains 5-field cron expressions in plain English.
package cron

import (
	"fmt"
	"strings"
)

// explainMinutes returns the opening phrase of the explanation.
// When both minutes and hours are fixed values it formats as an HH:MM clock string.
func explainMinutes(minutes, hours fieldSpec) string {
	if minutes.kind == fieldSingle && hours.kind == fieldSingle {
		return fmt.Sprintf("at %02d:%02d", hours.lo, minutes.lo)
	}

	switch minutes.kind {
	case fieldWildcard:
		return "every minute"
	case fieldSingle:
		return fmt.Sprintf("at minute %d", minutes.lo)
	case fieldRange:
		return fmt.Sprintf("at minutes %d through %d", minutes.lo, minutes.hi)
	default:
		return ""
	}
}

// explainHours returns the "past hour(s)" clause, or empty when both fields are fixed (already captured as HH:MM).
func explainHours(minutes, hours fieldSpec) string {
	switch hours.kind {
	case fieldSingle:
		if minutes.kind == fieldWildcard || minutes.kind == fieldRange {
			return fmt.Sprintf("past hour %d", hours.lo)
		}
	case fieldRange:
		return fmt.Sprintf("past hours %d through %d", hours.lo, hours.hi)
	default:
		return ""
	}
	return ""
}

var dayOfWeekName = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

var monthName = [12]string{
	"January", "February", "March", "April",
	"May", "June", "July", "August",
	"September", "October", "November", "December",
}

// explainField dispatches a non-wildcard fieldSpec to single or rangeFn; returns empty for wildcards.
func explainField(field fieldSpec, single func(int) string, rangeFn func(int, int) string) string {
	switch field.kind {
	case fieldSingle:
		return single(field.lo)
	case fieldRange:
		return rangeFn(field.lo, field.hi)
	default:
		return ""
	}
}

// explainDaysOfMonth returns the day-of-month clause, or empty for wildcards.
func explainDaysOfMonth(cron cronExpression) string {
	return explainField(
		cron.daysOfMonth,
		func(day int) string {
			return fmt.Sprintf("on day of month %d", day)
		},
		func(start, end int) string {
			return fmt.Sprintf("on days of month %d through %d", start, end)
		},
	)
}

// explainDaysOfWeek returns the day-of-week clause, or empty for wildcards.
func explainDaysOfWeek(cron cronExpression) string {
	return explainField(
		cron.daysOfWeek,
		func(day int) string {
			return fmt.Sprintf("on %s", dayOfWeekName[day])
		},
		func(start, end int) string {
			return fmt.Sprintf("on %s through %s", dayOfWeekName[start], dayOfWeekName[end])
		},
	)
}

// explainMonths returns the month clause, or empty for wildcards.
func explainMonths(cron cronExpression) string {
	return explainField(
		cron.months,
		func(month int) string {
			return fmt.Sprintf("in %s", monthName[month-1])
		},
		func(start, end int) string {
			return fmt.Sprintf("in %s through %s", monthName[start-1], monthName[end-1])
		},
	)
}

// ExplainExpression validates, parses, and returns a plain-English description of a 5-field cron expression.
func ExplainExpression(expr string) (string, error) {
	if err := validate(expr); err != nil {
		return "", err
	}

	cron, err := parseExpression(expr)
	if err != nil {
		return "", err
	}

	explanation := []string{explainMinutes(cron.minutes, cron.hours)}

	if hoursExplanation := explainHours(cron.minutes, cron.hours); hoursExplanation != "" {
		explanation = append(explanation, hoursExplanation)
	}

	if daysExplanation := explainDaysOfMonth(*cron); daysExplanation != "" {
		explanation = append(explanation, daysExplanation)
	}

	if daysOfWeekExplanation := explainDaysOfWeek(*cron); daysOfWeekExplanation != "" {
		explanation = append(explanation, daysOfWeekExplanation)
	}

	if monthsExplanation := explainMonths(*cron); monthsExplanation != "" {
		explanation = append(explanation, monthsExplanation)
	}

	return strings.Join(explanation, " "), nil
}
