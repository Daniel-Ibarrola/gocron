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
