package validator

import (
	"go-api/internal/models"
	"testing"
)

func TestValidateMovie(t *testing.T) {
	type args struct {
		v     *Validator
		movie *models.Movie
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Empty title",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title: "",
				},
			},
			expected: false,
		},
		{
			name: "Empty year",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title: "top gun",
					Year:  0,
				},
			},
			expected: false,
		},
		{
			name: "Year a long time ago",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title: "top gun",
					Year:  1700,
				},
			},
			expected: false,
		},
		{
			name: "Runtime negative",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title:   "top gun",
					Year:    2005,
					Runtime: -100,
				},
			},
			expected: false,
		},
		{
			name: "Not unique genres",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title:   "top gun",
					Year:    2005,
					Runtime: 100,
					Genres:  []string{"action", "action"},
				},
			},
			expected: false,
		},
		{
			name: "Valid movie input",
			args: args{
				v: New(),
				movie: &models.Movie{
					Title:   "top gun",
					Year:    2005,
					Runtime: 100,
					Genres:  []string{"action"},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateMovie(tt.args.v, tt.args.movie)
			if tt.args.v.Valid() != tt.expected {
				t.Errorf("got %v expected %v", tt.args.v.Valid(), tt.expected)
			}
		})
	}
}
