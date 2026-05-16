package cron

import "testing"

func TestValidateValidCronExpressions(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "all wildcards",
			expr: "* * * * *",
		},
		{
			name: "fixed numeric values",
			expr: "0 12 1 1 0",
		},
		{
			name: "step expression",
			expr: "*/5 * * * *",
		},
		{
			name: "range expression",
			expr: "0 9-17 * * 1-5",
		},
		{
			name: "list expression",
			expr: "0,15,30,45 * * * *",
		},
		{
			name: "extra whitespace between fields",
			expr: "0   12  1   1  0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := validate(test.expr); err != nil {
				t.Fatalf("expected expression %q to be valid, got error: %v", test.expr, err)
			}
		})
	}
}

func TestValidateInvalidNumberOfFields(t *testing.T) {
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
			err := validate(test.expr)
			if err == nil {
				t.Fatalf("expected validation error for expression %q", test.expr)
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

func TestValidateInvalidCharacters(t *testing.T) {
	tests := []struct {
		name string
		expr string
	}{
		{
			name: "letters",
			expr: "a * * * *",
		},
		{
			name: "question mark",
			expr: "? * * * *",
		},
		{
			name: "at sign",
			expr: "@daily * * * *",
		},
		{
			name: "semicolon",
			expr: "0 12 * * ;",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validate(test.expr)
			if err == nil {
				t.Fatalf("expected validation error for expression %q", test.expr)
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("expected *ValidationError, got %T", err)
			}

			if validationErr.Expression != test.expr {
				t.Errorf("expected expression %q, got %q", test.expr, validationErr.Expression)
			}

			if validationErr.Message != "Invalid characters in expression" {
				t.Errorf("expected message %q, got %q", "Invalid characters in expression", validationErr.Message)
			}
		})
	}
}

func TestValidationErrorErrorReturnsMessage(t *testing.T) {
	err := &ValidationError{
		Expression: "* * * *",
		Message:    "Invalid number of fields",
	}

	if err.Error() != err.Message {
		t.Errorf("expected Error() to return %q, got %q", err.Message, err.Error())
	}
}
