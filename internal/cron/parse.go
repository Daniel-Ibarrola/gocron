package cron

import (
	"strconv"
	"strings"
)

// FieldError is returned when a single cron field token cannot be parsed or is out of range.
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}

// fieldKind discriminates how a cron field was specified.
type fieldKind int

const (
	fieldWildcard fieldKind = iota
	fieldSingle
	fieldRange
)

// fieldSpec holds the parsed representation of one cron field.
type fieldSpec struct {
	kind   fieldKind
	lo, hi int // single: lo==hi; range: lo<=hi; wildcard: zero
}

// cronExpression holds all five parsed fields of a cron expression (minute hour dom month dow).
type cronExpression struct {
	minutes, hours, daysOfMonth, months, daysOfWeek fieldSpec
}

// cronField identifies one of the five positional fields in a cron expression.
type cronField int

const (
	minuteField cronField = iota
	hourField
	dayOfMonthField
	monthField
	dayOfWeekField
	cronFieldCount
)

// rangeSpec holds the inclusive valid numeric bounds for a given cronField.
type rangeSpec struct {
	lo, hi int
}

var validRanges = map[cronField]rangeSpec{
	minuteField:     {0, 59},
	hourField:       {0, 23},
	dayOfMonthField: {1, 31},
	monthField:      {1, 12},
	dayOfWeekField:  {0, 6},
}

// parseWildcard returns a wildcard fieldSpec.
func parseWildcard(field cronField) (fieldSpec, error) {
	return fieldSpec{kind: fieldWildcard}, nil
}

// parseSingle parses a bare integer token and validates it is within the field's allowed bounds.
func parseSingle(val string, field cronField) (fieldSpec, error) {
	valNum, err := strconv.Atoi(val)
	if err != nil {
		return fieldSpec{}, &FieldError{Field: val, Message: "Invalid integer value"}
	}
	if valNum < validRanges[field].lo || valNum > validRanges[field].hi {
		return fieldSpec{}, &FieldError{Field: val, Message: "value out of range"}
	}
	return fieldSpec{kind: fieldSingle, lo: valNum, hi: valNum}, nil
}

// parseRange parses an "a-b" token; a degenerate range where a==b is normalized to fieldSingle.
func parseRange(val string, field cronField) (fieldSpec, error) {
	split := strings.Split(val, "-")
	low, err := strconv.Atoi(split[0])
	if err != nil {
		return fieldSpec{}, &FieldError{Field: val, Message: "Invalid integer value"}
	}
	high, err := strconv.Atoi(split[1])
	if err != nil {
		return fieldSpec{}, &FieldError{Field: val, Message: "Invalid integer value"}
	}
	if low > high {
		return fieldSpec{}, &FieldError{Field: val, Message: "Invalid range"}
	}
	if low < validRanges[field].lo || high > validRanges[field].hi {
		return fieldSpec{}, &FieldError{Field: val, Message: "value out of range"}
	}
	if low == high {
		return fieldSpec{kind: fieldSingle, lo: low, hi: low}, nil
	}
	return fieldSpec{kind: fieldRange, lo: low, hi: high}, nil
}

// parseField dispatches to the appropriate parser based on the shape of the token.
func parseField(val string, field cronField) (fieldSpec, error) {
	switch {
	case val == "*":
		return parseWildcard(field)
	case strings.Contains(val, "-"):
		return parseRange(val, field)
	default:
		return parseSingle(val, field)
	}
}

// parseExpression splits expr into fields and parses each one in positional order.
func parseExpression(expr string) (*cronExpression, error) {
	fields := strings.Fields(expr)
	if len(fields) != int(cronFieldCount) {
		return nil, &ValidationError{Expression: expr, Message: "Invalid number of fields"}
	}

	parsedFields := make([]fieldSpec, len(fields))
	for i, field := range fields {
		parsed, err := parseField(field, cronField(i))
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
