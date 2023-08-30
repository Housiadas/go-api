package validator

import (
	"go-api/internal/filters"
	"testing"
)

func TestValidateFilters(t *testing.T) {
	type args struct {
		v *Validator
		f filters.Filters
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Invalid sort value",
			args: args{
				v: New(),
				f: filters.Filters{
					Page:         1,
					PageSize:     20,
					Sort:         "-id",
					SortSafelist: []string{"foo", "bar"},
				},
			},
			expected: false,
		},
		{
			name: "Valid filters",
			args: args{
				v: New(),
				f: filters.Filters{
					Page:         1,
					PageSize:     20,
					Sort:         "-id",
					SortSafelist: []string{"id", "title", "-id"},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateFilters(tt.args.v, tt.args.f)
			if tt.args.v.Valid() != tt.expected {
				t.Errorf("got %v expected %v", tt.args.v.Valid(), tt.expected)
			}
		})
	}
}
