package filters

import (
	"testing"
)

func TestFilters_Limit(t *testing.T) {
	f := Filters{
		Page:     1,
		PageSize: 20,
	}
	sut := f.Limit()
	if sut == f.Page {
		t.Errorf("got %q; expected %q", sut, f.Page)
	}
}

func TestFilters_Offset(t *testing.T) {
	tests := []struct {
		name     string
		f        Filters
		expected int
	}{
		{
			name: "No offset",
			f: Filters{
				Page:     1,
				PageSize: 20,
			},
			expected: 0,
		},
		{
			name: "Offset 20 rows",
			f: Filters{
				Page:     2,
				PageSize: 20,
			},
			expected: 20,
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.Offset()
			if sut != tt.expected {
				t.Errorf("got %q; expected %q", sut, tt.expected)
			}
		})
	}
}

func TestFilters_SortColumn(t *testing.T) {
	tests := []struct {
		name     string
		f        Filters
		expected string
	}{
		{
			name: "Sort Column year ASC",
			f: Filters{
				Sort:         "year",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			expected: "year",
		},
		{
			name: "Sort Column title DESC",
			f: Filters{
				Sort:         "-title",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			expected: "title",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.SortColumn()
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}

func TestFilters_SortDirection(t *testing.T) {
	tests := []struct {
		name     string
		f        Filters
		expected string
	}{
		{
			name: "Sort Column year ASC",
			f: Filters{
				Sort:         "year",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			expected: "ASC",
		},
		{
			name: "Sort Column title DESC",
			f: Filters{
				Sort:         "-title",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			expected: "DESC",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.SortDirection()
			if sut != tt.expected {
				t.Errorf("got %v expected %v", sut, tt.expected)
			}
		})
	}
}
