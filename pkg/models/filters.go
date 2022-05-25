package models

import "strings"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f Filters) Valid() bool {
	if f.Page < 1 || f.Page > 1_000_000 {
		return false
	}

	if f.PageSize < 1 || f.PageSize > 100 {
		return false
	}

	var sort string
	for _, val := range f.SortSafeList {
		if f.Sort == val {
			sort = val
		}
	}

	if sort == "" {
		return false
	}

	return true
}

func (f Filters) SortColumn() string {
	for _, val := range f.SortSafeList {
		if f.Sort == val {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort value: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
