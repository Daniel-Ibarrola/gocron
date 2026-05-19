package cron

import "testing"

func TestExplainSimpleCronExpressions(t *testing.T) {
	// Test expression that only contain wildcard characters and/or fixed values
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"every minute", "* * * * *", "every minute"},
		{"at minute", "5 * * * *", "at minute 5"},
		{"every minute past an hour", "* 4 * * *", "every minute past hour 4"},
		{"minute and hour", "30 6 * * *", "at 06:30"},
		{"every minute at day of month", "* * 2 * *", "every minute on day of month 2"},
		{"at minute at day of month", "20 * 10 * *", "at minute 20 on day of month 10"},
		{"at minute and hour at day of month", "20 4 5 * *", "at 04:20 on day of month 5"},
		{"at minute in month", "20 * * 5 *", "at minute 20 in May"},
		{"at minute and hour in month", "20 4 * 9 *", "at 04:20 in September"},
		{"at minute and hour in day of month at month", "20 4 5 8 *", "at 04:20 on day of month 5 in August"},
		{"every minute at day of week", "* * * * 1", "every minute on Monday"},
		{"at minute on day of week", "20 * * * 3", "at minute 20 on Wednesday"},
		{"at minute and hour on day of week", "20 4 * * 5", "at 04:20 on Friday"},
		{"at minute and hour on day of month and day of week", "20 4 4 * 5", "at 04:20 on day of month 4 on Friday"},
		{"at minute and hour on day of month and on month and day of week", "20 4 4 10 5", "at 04:20 on day of month 4 on Friday in October"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			explanation, err := ExplainExpression(test.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if explanation != test.expected {
				t.Errorf("expected explanation %q, got %q", test.expected, explanation)
			}
		})
	}
}

func TestInvalidCronExpressions(t *testing.T) {
	// Test expressions that are invalid like missing fields or completely wrong
	tests := []string{
		"* * * *",
		"* * * * * *",
		"hi how are you pal",
	}
	for _, expr := range tests {
		t.Run(expr, func(t *testing.T) {
			_, err := ExplainExpression(expr)
			if err == nil {
				t.Errorf("expected error for invalid expression %q, got nil", expr)
			}
		})
	}
}

func TestExplainRangeExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "minute and hour ranges",
			expr:     "5-10 9-17 * * *",
			expected: "at minutes 5 through 10 past hours 9 through 17",
		},
		{
			name:     "fixed minute with hour range",
			expr:     "5 9-17 * * *",
			expected: "at minute 5 past hours 9 through 17",
		},
		{
			name:     "minute range with fixed hour",
			expr:     "5-10 4 * * *",
			expected: "at minutes 5 through 10 past hour 4",
		},
		{
			name:     "month range",
			expr:     "* * * 6-8 *",
			expected: "every minute in June through August",
		},
		{
			name:     "day of week range",
			expr:     "* * * * 1-5",
			expected: "every minute on Monday through Friday",
		},
		{
			name:     "day of month range",
			expr:     "* * 1-15 * *",
			expected: "every minute on days of month 1 through 15",
		},
		{
			name:     "degenerate range renders like single",
			expr:     "5-5 * * * *",
			expected: "at minute 5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			explanation, err := ExplainExpression(test.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if explanation != test.expected {
				t.Errorf("expected explanation %q, got %q", test.expected, explanation)
			}
		})
	}
}

func TestExplainFieldErrors(t *testing.T) {
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
			name:        "out-of-bounds day of month range",
			expr:        "* * 0-32 * *",
			invalidPart: "0-32",
		},
		{
			name:        "out-of-bounds day of week range",
			expr:        "* * * * 0-7",
			invalidPart: "0-7",
		},
		{
			name:        "out-of-bounds month range",
			expr:        "* * * 0-13 *",
			invalidPart: "0-13",
		},
		{
			name:        "out-of-bounds single minute",
			expr:        "99 * * * *",
			invalidPart: "99",
		},
		{
			name:        "out-of-bounds single hour",
			expr:        "* 24 * * *",
			invalidPart: "24",
		},
		{
			name:        "out-of-bounds single day of month",
			expr:        "* * 32 * *",
			invalidPart: "32",
		},
		{
			name:        "out-of-bounds single month",
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
			explanation, err := ExplainExpression(test.expr)
			if err == nil {
				t.Fatalf("expected error for expression %q, got nil (explanation: %q)", test.expr, explanation)
			}
			if explanation != "" {
				t.Errorf("expected empty explanation on error, got %q", explanation)
			}

			fieldErr, ok := err.(*FieldError)
			if !ok {
				t.Fatalf("expected *FieldError, got %T (%v)", err, err)
			}
			if fieldErr.Field != test.invalidPart {
				t.Errorf("expected field %q, got %q", test.invalidPart, fieldErr.Field)
			}
		})
	}
}
