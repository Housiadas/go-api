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
		t.Errorf("got %q; want %q", sut, f.Page)
	}
}

func TestFilters_Offset(t *testing.T) {
	tests := []struct {
		name string
		f    Filters
		want int
	}{
		{
			name: "No offset",
			f: Filters{
				Page:     1,
				PageSize: 20,
			},
			want: 0,
		},
		{
			name: "Offset 20 rows",
			f: Filters{
				Page:     2,
				PageSize: 20,
			},
			want: 20,
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.Offset()
			if sut != tt.want {
				t.Errorf("got %q; want %q", sut, tt.want)
			}
		})
	}
}

func TestFilters_SortColumn(t *testing.T) {
	tests := []struct {
		name string
		f    Filters
		want string
	}{
		{
			name: "Sort Column year ASC",
			f: Filters{
				Sort:         "year",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			want: "year",
		},
		{
			name: "Sort Column title DESC",
			f: Filters{
				Sort:         "-title",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			want: "title",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.SortColumn()
			if sut != tt.want {
				t.Errorf("got %q; want %q", sut, tt.want)
			}
		})
	}
}

func TestFilters_SortDirection(t *testing.T) {
	tests := []struct {
		name string
		f    Filters
		want string
	}{
		{
			name: "Sort Column year ASC",
			f: Filters{
				Sort:         "year",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			want: "ASC",
		},
		{
			name: "Sort Column title DESC",
			f: Filters{
				Sort:         "-title",
				SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
			},
			want: "DESC",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := tt.f.SortDirection()
			if sut != tt.want {
				t.Errorf("got %q; want %q", sut, tt.want)
			}
		})
	}
}
