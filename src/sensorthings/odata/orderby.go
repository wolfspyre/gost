package odata

import (
	"net/http"
	"strings"
)

// QueryOrderBy is used to define the order of the results set be ascending (asc) or descending (desc) order.
// orderby Is used to specify which properties are used to order the collection of entities identified by the resource path.
type QueryOrderBy struct {
	QueryBase
	Suffix   string
	Property string
}

// Parse tries to parse the OrderBy query into a suffix and property value
// if the suffix is not "asc" or "desc" or no property is given then
// parse will return an error
func (q *QueryOrderBy) Parse(value string) error {
	q.RawQuery = value
	ob := strings.Split(value, " ")
	if len(ob) != 2 || (strings.ToLower(ob[1]) != "asc" && strings.ToLower(ob[1]) != "desc") || len(ob[0]) < 1 {
		return CreateQueryError(QueryOrderByInvalid, http.StatusBadRequest, value)
	}

	q.Property = ob[0]
	q.Suffix = ob[1]
	return nil
}

// IsValid checks if the given property value in the request is valid for the
// used endpoint, returns an error if not supported
// values = available properties of an entity
func (q *QueryOrderBy) IsValid(values []string) (bool, error) {
	for _, s := range values {
		if q.Property == s {
			return true, nil
		}
	}

	return false, CreateQueryError(QueryOrderByInvalid, http.StatusBadRequest, q.Property)
}

// GetQueryOptionType returns the QueryOptionType for QueryOrderBy
func (q *QueryOrderBy) GetQueryOptionType() QueryOptionType {
	return QueryOptionOrderBy
}

// IsNil checks if *QueryOrderBy is nil
func (q *QueryOrderBy) IsNil() bool {
	if q == nil {
		return true
	}

	return false
}
