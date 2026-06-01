package cron

import "testing"

func wild() fieldSpec        { return fieldSpec{kind: fieldWildcard} }
func single(v int) fieldSpec { return fieldSpec{kind: fieldSingle, lo: v, hi: v} }
func rng(lo, hi int) fieldSpec {
	return fieldSpec{kind: fieldRange, lo: lo, hi: hi}
}

func TestParseExpressionValid(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected cronExpression
	}{
		{
			name:     "all wildcards",
			expr:     "* * * * *",
			expected: cronExpression{minutes: wild(), hours: wild(), daysOfMonth: wild(), months: wild(), daysOfWeek: wild()},
		},
		{
			name:     "all numeric values",
			expr:     "0 12 1 1 0",
			expected: cronExpression{minutes: single(0), hours: single(12), daysOfMonth: single(1), months: single(1), daysOfWeek: single(0)},
		},
		{
			name:     "field ordering preserved",
			expr:     "5 10 15 7 3",
			expected: cronExpression{minutes: single(5), hours: single(10), daysOfMonth: single(15), months: single(7), daysOfWeek: single(3)},
		},
		{
			name:     "mix of wildcards and numbers",
			expr:     "30 * * 6 *",
			expected: cronExpression{minutes: single(30), hours: wild(), daysOfMonth: wild(), months: single(6), daysOfWeek: wild()},
		},
		{
			name:     "extra whitespace between fields",
			expr:     "0   12  1   1  0",
			expected: cronExpression{minutes: single(0), hours: single(12), daysOfMonth: single(1), months: single(1), daysOfWeek: single(0)},
		},
		{
			name:     "minute range",
			expr:     "5-10 * * * *",
			expected: cronExpression{minutes: rng(5, 10), hours: wild(), daysOfMonth: wild(), months: wild(), daysOfWeek: wild()},
		},
		{
			name:     "hour range",
			expr:     "* 9-17 * * *",
			expected: cronExpression{minutes: wild(), hours: rng(9, 17), daysOfMonth: wild(), months: wild(), daysOfWeek: wild()},
		},
		{
			name:     "day of month range",
			expr:     "* * 1-15 * *",
			expected: cronExpression{minutes: wild(), hours: wild(), daysOfMonth: rng(1, 15), months: wild(), daysOfWeek: wild()},
		},
		{
			name:     "month range",
			expr:     "* * * 6-8 *",
			expected: cronExpression{minutes: wild(), hours: wild(), daysOfMonth: wild(), months: rng(6, 8), daysOfWeek: wild()},
		},
		{
			name:     "day of week range",
			expr:     "* * * * 1-5",
			expected: cronExpression{minutes: wild(), hours: wild(), daysOfMonth: wild(), months: wild(), daysOfWeek: rng(1, 5)},
		},
		{
			name:     "range mixed with single and wildcard",
			expr:     "5 9-17 * * 1-5",
			expected: cronExpression{minutes: single(5), hours: rng(9, 17), daysOfMonth: wild(), months: wild(), daysOfWeek: rng(1, 5)},
		},
		{
			name:     "degenerate range normalized to single",
			expr:     "5-5 * * * *",
			expected: cronExpression{minutes: single(5), hours: wild(), daysOfMonth: wild(), months: wild(), daysOfWeek: wild()},
		},
		{
			name:     "boundary ranges accepted",
			expr:     "0-59 0-23 1-31 1-12 0-6",
			expected: cronExpression{minutes: rng(0, 59), hours: rng(0, 23), daysOfMonth: rng(1, 31), months: rng(1, 12), daysOfWeek: rng(0, 6)},
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

func TestParseExpressionInvalidRange(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		invalidPart string
	}{
		{
			name:        "inverted minute range",
			expr:        "5-1 * * * *",
			invalidPart: "5-1",
		},
		{
			name:        "inverted hour range",
			expr:        "* 10-2 * * *",
			invalidPart: "10-2",
		},
		{
			name:        "out-of-bounds minute range",
			expr:        "0-99 * * * *",
			invalidPart: "0-99",
		},
		{
			name:        "out-of-bounds hour range",
			expr:        "* 0-24 * * *",
			invalidPart: "0-24",
		},
		{
			name:        "out-of-bounds day of month range",
			expr:        "* * 0-32 * *",
			invalidPart: "0-32",
		},
		{
			name:        "out-of-bounds month range",
			expr:        "* * * 0-13 *",
			invalidPart: "0-13",
		},
		{
			name:        "out-of-bounds day of week range",
			expr:        "* * * * 0-7",
			invalidPart: "0-7",
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
		})
	}
}

func TestParseExpressionOutOfBoundsSingle(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		invalidPart string
	}{
		{
			name:        "minute too high",
			expr:        "99 * * * *",
			invalidPart: "99",
		},
		{
			name:        "hour too high",
			expr:        "* 24 * * *",
			invalidPart: "24",
		},
		{
			name:        "day of month zero",
			expr:        "* * 0 * *",
			invalidPart: "0",
		},
		{
			name:        "day of month too high",
			expr:        "* * 32 * *",
			invalidPart: "32",
		},
		{
			name:        "month zero",
			expr:        "* * * 0 *",
			invalidPart: "0",
		},
		{
			name:        "month too high",
			expr:        "* * * 13 *",
			invalidPart: "13",
		},
		{
			name:        "day of week 7 rejected",
			expr:        "* * * * 7",
			invalidPart: "7",
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
