package cron

import (
	"regexp"
	"strings"
)

type ValidationError struct {
	Expression string
	Message    string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Define a regex that matches characters NOT in our allowed list
// Allowed: 0-9, *, /, -, comma, and space
var cronSafeRegex = regexp.MustCompile(`^[0-9*/\-, ]+$`)

// Validate that the cron expression contains all fields and only valid characters
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
