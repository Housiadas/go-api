package validator

import (
	"regexp"
	"testing"
)

func TestValidator_Valid(t *testing.T) {
	tests := []struct {
		name     string
		v        Validator
		expected bool
	}{
		{
			name: "Test validator no errors",
			v: Validator{
				Errors: map[string]string{},
			},
			expected: true,
		},
		{
			name: "Test validator has errors",
			v: Validator{
				Errors: map[string]string{
					"error": "error",
				},
			},
			expected: false,
		},
	}
	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.v.Valid()
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		matches  []string
		expected bool
	}{
		{
			name:     "Test validator In method valid",
			value:    "test",
			matches:  []string{"foo", "bar"},
			expected: false,
		},
		{
			name:     "Test validator In invalid",
			value:    "test",
			matches:  []string{"foo", "test", "bar"},
			expected: true,
		},
	}
	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := In(tt.value, tt.matches...)
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}

func TestMatchesEmailRegExp(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		regExp   *regexp.Regexp
		expected bool
	}{
		{
			name:     "Test email regex #1",
			value:    "test",
			regExp:   EmailRX,
			expected: false,
		},
		{
			name:     "Test email regex #2",
			value:    "test@test.com",
			regExp:   EmailRX,
			expected: true,
		},
		{
			name:     "Test email regex #4",
			value:    "@@test",
			regExp:   EmailRX,
			expected: false,
		},
		{
			name:     "Test email regex #5",
			value:    "test@test.",
			regExp:   EmailRX,
			expected: false,
		},
	}
	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := Matches(tt.value, tt.regExp)
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{
			name:     "Test has unique values",
			values:   []string{"test", "foo", "bar"},
			expected: true,
		},
		{
			name:     "Test has not values",
			values:   []string{"test", "foo", "test"},
			expected: false,
		},
	}
	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := Unique(tt.values)
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}
