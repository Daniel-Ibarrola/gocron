package cron

import (
	"fmt"
	"strconv"
	"strings"
)

type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}

type cronExpression struct {
	minutes, hours, daysOfMonth, months, daysOfWeek int
}

// Parse a field of the cron expression into a slice of integers
func parseField(val string) (int, error) {
	if val == "*" {
		return -1, nil
	}
	valNum, err := strconv.Atoi(val)
	if err != nil {
		return -1, &FieldError{Field: val, Message: "Invalid integer value"}
	}
	return valNum, nil
}

const (
	minuteField = iota
	hourField
	dayOfMonthField
	monthField
	dayOfWeekField
	cronFieldCount
)

func parseExpression(expr string) (*cronExpression, error) {
	fields := strings.Fields(expr)
	if len(fields) != cronFieldCount {
		return nil, &ValidationError{Expression: expr, Message: "Invalid number of fields"}
	}

	parsedFields := make([]int, len(fields))
	for i, field := range fields {
		parsed, err := parseField(field)
		if err != nil {
			return nil, err
		}
		parsedFields[i] = parsed
	}

	cron := cronExpression{
		minutes:     parsedFields[minuteField],
		hours:       parsedFields[hourField],
		daysOfMonth: parsedFields[dayOfMonthField],
		months:      parsedFields[monthField],
		daysOfWeek:  parsedFields[dayOfWeekField],
	}
	return &cron, nil
}

// ExplainExpression provides a human-readable explanation of the cron expression
func ExplainExpression(expr string) (string, error) {
	if err := validate(expr); err != nil {
		return "", err
	}

	cron, err := parseExpression(expr)
	if err != nil {
		return "", err
	}

	var explanation []string
	var monthName = [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	var dayOfWeekName = [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	if cron.minutes > -1 && cron.hours > -1 {
		explanation = append(explanation, fmt.Sprintf("at %02d:%02d", cron.hours, cron.minutes))
	} else if cron.minutes == -1 {
		explanation = append(explanation, "every minute")
	} else {
		explanation = append(explanation, fmt.Sprintf("at minute %d", cron.minutes))
	}

	if cron.minutes == -1 && cron.hours > -1 {
		explanation = append(explanation, fmt.Sprintf("past hour %d", cron.hours))
	}

	if cron.daysOfMonth > -1 {
		explanation = append(explanation, fmt.Sprintf("on day of month %d", cron.daysOfMonth))
	}

	if cron.daysOfWeek > -1 {
		explanation = append(explanation, fmt.Sprintf("on %s", dayOfWeekName[cron.daysOfWeek]))
	}

	if cron.months > -1 {
		explanation = append(explanation, fmt.Sprintf("in %s", monthName[cron.months-1]))
	}

	return strings.Join(explanation, " "), nil
}
