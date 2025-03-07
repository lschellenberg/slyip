package common

import (
	"net/http"
	"strconv"
	"strings"
	"yip/src/slyerrors"
)

const (
	SortingOrderASC        = "ASC"
	SortingOrderDESC       = "DESC"
	DefaultMaxPageSize     = 100
	DefaultDefaultPageSize = 50
	DefaultSortBy          = "createdAt"
	DefaultSortOrder       = SortingOrderDESC
)

type AllowedFilters map[string]struct{}

var DefaultAllowedFilters = map[string]struct{}{
	DefaultSortBy: {},
}

type QueryFilters struct {
	Filters AllowedFilters `json:"-" schema:"-"`
}

// swagger:parameters PaginationQuery
type PaginationQuery struct {
	QueryFilters

	// in:query
	// minimum: 1
	// default: 50
	// maximum: 100
	PageSize int `schema:"pageSize" json:"pageSize"`
	// in:query
	// minimum: 0
	// default: 0
	Offset int `schema:"offset" json:"offset"`
	// in:query
	// default: createdAt
	SortBy string `schema:"sortBy" json:"sortBy"`
	// in:query
	// enum: ASC,DESC
	// default: DESC
	SortOrder string `schema:"order" json:"order"`
}

func PaginationQueryFromRequest(r *http.Request) (*PaginationQuery, error) {
	q := r.URL.Query()
	pq := &PaginationQuery{}
	var err error

	pq.PageSize, err = emptyOrInt(q.Get("pageSize"))
	if err != nil {
		return nil, err
	}

	pq.Offset, err = emptyOrInt(q.Get("offset"))
	if err != nil {
		return nil, err
	}

	pq.SortOrder = q.Get("sortOrder")
	pq.SortOrder = q.Get("sortBy")

	if err := pq.Validate(); err != nil {
		return nil, err
	}

	return pq, nil
}

func (query *PaginationQuery) Validate() error {
	val := slyerrors.NewValidation("paginationQuery")
	if query.PageSize < 0 || query.PageSize == 0 {
		query.PageSize = DefaultDefaultPageSize
	} else if query.PageSize > DefaultMaxPageSize {
		query.PageSize = DefaultMaxPageSize
	}

	if query.Offset < 0 {
		val.Add("offset", slyerrors.ValidationCodeNumberTooSmall, strconv.Itoa(query.Offset))
	}

	query.SortOrder = strings.ToUpper(query.SortOrder)

	if query.SortOrder == "" {
		query.SortOrder = DefaultSortOrder
	} else if query.SortOrder != SortingOrderASC && query.SortOrder != SortingOrderDESC {
		val.Add("order", slyerrors.ValidationCodeUnexpectedValue, "must be %s or %s", SortingOrderASC, SortingOrderDESC)
	}

	if query.QueryFilters.Filters == nil {
		query.QueryFilters.Filters = DefaultAllowedFilters
	}

	if query.SortBy != "" {
		if _, ok := query.QueryFilters.Filters[query.SortBy]; !ok {
			val.Add("sortBy", slyerrors.ValidationCodeUnexpectedValue, query.SortBy)
		}
	} else {
		query.SortBy = DefaultSortBy
	}

	return val.Error()
}

// ValidatePanic is meant to ease testing. That is all.
// It works and returns a copy of the current object.
func (r PaginationQuery) ValidatePanic() PaginationQuery {
	if err := r.Validate(); err != nil {
		panic(err)
	}
	return r
}

// PaginationResponse provides the middleware elements of a paginated response.
// You MUST provide the items field yourself: Items []<concrete type> `json:"items"`
type PaginationResponse struct {
	// in:body
	PageSize int `schema:"pageSize" json:"pageSize"`
	// in:body
	Offset int `schema:"offset" json:"offset"`
	// in:body
	SortBy string `schema:"sortBy" json:"sortBy"`
	// in:body
	SortOrder string `schema:"order" json:"order"`
	// in:body
	Total uint64 `json:"total"`
	// in:body
	Count uint64 `json:"count"`
	// swagger:ignore
	// in:body
	Filter map[string]string `json:"filter"`
}

func (r PaginationResponse) From(query PaginationQuery) PaginationResponse {
	r.PageSize = query.PageSize
	r.Offset = query.Offset
	r.SortOrder = query.SortOrder
	r.SortBy = query.SortBy
	return r
}

func emptyOrInt(n string) (int, error) {
	if n == "" {
		return 0, nil
	}
	return strconv.Atoi(n)
}
