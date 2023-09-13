package validator

import "go-api/internal/filters"

func ValidateFilters(v *Validator, f filters.Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
