package cron

import "testing"

func TestParseExpressionValid(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected cronExpression
	}{
		{
			name:     "all wildcards",
			expr:     "* * * * *",
			expected: cronExpression{minutes: -1, hours: -1, daysOfMonth: -1, months: -1, daysOfWeek: -1},
		},
		{
			name:     "all numeric values",
			expr:     "0 12 1 1 0",
			expected: cronExpression{minutes: 0, hours: 12, daysOfMonth: 1, months: 1, daysOfWeek: 0},
		},
		{
			name:     "field ordering preserved",
			expr:     "5 10 15 7 3",
			expected: cronExpression{minutes: 5, hours: 10, daysOfMonth: 15, months: 7, daysOfWeek: 3},
		},
		{
			name:     "mix of wildcards and numbers",
			expr:     "30 * * 6 *",
			expected: cronExpression{minutes: 30, hours: -1, daysOfMonth: -1, months: 6, daysOfWeek: -1},
		},
		{
			name:     "extra whitespace between fields",
			expr:     "0   12  1   1  0",
			expected: cronExpression{minutes: 0, hours: 12, daysOfMonth: 1, months: 1, daysOfWeek: 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cron, err := parseExpression(test.expr)
			if err != nil {
				t.Fatalf("expected expression %q to parse, got error: %v", test.expr, err)
			}
			if *cron != test.expected {
				t.Errorf("expected %+v, got %+v", test.expected, *cron)
			}
		})
	}
}

func TestParseExpressionInvalidNumberOfFields(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "empty expression",
			expr: "",
		},
		{
			name: "too few fields",
			expr: "* * * *",
		},
		{
			name: "too many fields",
			expr: "* * * * * *",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cron, err := parseExpression(test.expr)
			if err == nil {
				t.Fatalf("expected error for expression %q, got nil", test.expr)
			}
			if cron != nil {
				t.Errorf("expected nil cronExpression, got %+v", cron)
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("expected *ValidationError, got %T", err)
			}

			if validationErr.Expression != test.expr {
				t.Errorf("expected expression %q, got %q", test.expr, validationErr.Expression)
			}

			if validationErr.Message != "Invalid number of fields" {
				t.Errorf("expected message %q, got %q", "Invalid number of fields", validationErr.Message)
			}
		})
	}
}

func TestParseExpressionInvalidFieldValue(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		invalidPart string
	}{
		{
			name:        "letter in minute field",
			expr:        "a * * * *",
			invalidPart: "a",
		},
		{
			name:        "letter in last field",
			expr:        "* * * * x",
			invalidPart: "x",
		},
		{
			name:        "range expression not yet supported",
			expr:        "0 9-17 * * *",
			invalidPart: "9-17",
		},
		{
			name:        "step expression not yet supported",
			expr:        "*/5 * * * *",
			invalidPart: "*/5",
		},
		{
			name:        "list expression not yet supported",
			expr:        "0,15 * * * *",
			invalidPart: "0,15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cron, err := parseExpression(test.expr)
			if err == nil {
				t.Fatalf("expected error for expression %q, got nil", test.expr)
			}
			if cron != nil {
				t.Errorf("expected nil cronExpression, got %+v", cron)
			}

			fieldErr, ok := err.(*FieldError)
			if !ok {
				t.Fatalf("expected *FieldError, got %T", err)
			}

			if fieldErr.Field != test.invalidPart {
				t.Errorf("expected field %q, got %q", test.invalidPart, fieldErr.Field)
			}

			if fieldErr.Message != "Invalid integer value" {
				t.Errorf("expected message %q, got %q", "Invalid integer value", fieldErr.Message)
			}
		})
	}
}

func TestFieldErrorErrorReturnsMessage(t *testing.T) {
	err := &FieldError{
		Field:   "abc",
		Message: "Invalid integer value",
	}

	if err.Error() != err.Message {
		t.Errorf("expected Error() to return %q, got %q", err.Message, err.Error())
	}
}
