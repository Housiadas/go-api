package filters

import (
	"testing"
)

func TestCalculateMetadata(t *testing.T) {
	tests := []struct {
		name             string
		pageSize         int
		page             int
		totalRecords     int
		expectedLastPage int
	}{
		{
			name:             "Test last page 10",
			pageSize:         20,
			page:             1,
			totalRecords:     500,
			expectedLastPage: 25,
		},
		{
			name:             "Test last page 1",
			pageSize:         20,
			page:             1,
			totalRecords:     20,
			expectedLastPage: 1,
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := CalculateMetadata(tt.totalRecords, tt.page, tt.pageSize)
			if sut.LastPage != tt.expectedLastPage {
				t.Errorf("got %v expected %v", sut.LastPage, tt.expectedLastPage)
			}
		})
	}
}
