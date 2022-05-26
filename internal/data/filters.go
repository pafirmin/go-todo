package data

import (
	"math"
	"strings"

	"github.com/pafirmin/go-todo/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type MetaData struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func (f *Filters) Validate(v *validator.Validator) {
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page < 1_000_000, "page", "must be less than 1,000,000")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 1000, "page_size", "must be 1000 or lower")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort key")
}

func CalculateMetadata(totalRecords, page, pageSize int) MetaData {
	if totalRecords == 0 {
		return MetaData{}
	}

	return MetaData{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
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
