package filters

import (
	"strings"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // Add a SortSafe list field to hold the supported sort values.
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// SortColumn Check that the client-provided Sort field matches one of the entries in our safe list
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	//in theory this
	//shouldn’t happen — the Sort value should have already been checked by calling the
	//ValidateFilters() function — but this is a sensible failsafe to help stop a SQL injection
	//attack occurring.
	panic("unsafe sort parameter: " + f.Sort)
}

// SortDirection Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
