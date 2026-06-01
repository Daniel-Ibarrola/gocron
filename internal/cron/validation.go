package cron

import (
	"regexp"
	"strings"
)

// ValidationError is returned when the expression fails structural validation before field parsing.
type ValidationError struct {
	Expression string
	Message    string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var cronSafeRegex = regexp.MustCompile(`^[0-9*/\-, ]+$`)

// validate checks that expr has exactly five fields and contains only allowed characters (0-9 * / - , space).
func validate(expr string) error {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return &ValidationError{Expression: expr, Message: "Invalid number of fields"}
	}
	if !cronSafeRegex.MatchString(expr) {
		return &ValidationError{Expression: expr, Message: "Invalid characters in expression"}
	}
	return nil
}
