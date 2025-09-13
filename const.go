// Package hapi provides HTTP API query parameter parsing for filtering, sorting, and pagination.
// It supports various filter operators, sorting directions, and pagination controls.
package hapi

import "fmt"

// FilterOperator represents the comparison operator used in filter expressions.
type FilterOperator string

const (
	// FilterOperatorEqual is the default operator, not specifying it will use this
	FilterOperatorEqual          FilterOperator = "eq"
	FilterOperatorNotEqual       FilterOperator = "ne"
	FilterOperatorLike           FilterOperator = "lk"
	FilterOperatorNotLike        FilterOperator = "nlk"
	FilterOperatorIn             FilterOperator = "in"
	FilterOperatorNotIn          FilterOperator = "nin"
	FilterOperatorInLike         FilterOperator = "inlk"
	FilterOperatorNotInLike      FilterOperator = "ninlk"
	FilterOperatorGreaterThan    FilterOperator = "gt"
	FilterOperatorLessThan       FilterOperator = "lt"
	FilterOperatorGreaterOrEqual FilterOperator = "ge"
	FilterOperatorLessOrEqual    FilterOperator = "le"
)

// Valid checks if the filter operator is valid.
// Returns an error if the operator is not recognized.
func (o FilterOperator) Valid() error {
	switch o {
	case FilterOperatorEqual,
		FilterOperatorNotEqual,
		FilterOperatorLike,
		FilterOperatorNotLike,
		FilterOperatorIn,
		FilterOperatorNotIn,
		FilterOperatorInLike,
		FilterOperatorNotInLike,
		FilterOperatorGreaterThan,
		FilterOperatorLessThan,
		FilterOperatorGreaterOrEqual,
		FilterOperatorLessOrEqual:
		return nil
	}

	return fmt.Errorf("invalid operator: %q", o)
}

// IsList returns true if the operator expects multiple values (comma-separated).
func (o FilterOperator) IsList() bool {
	return o == FilterOperatorIn || o == FilterOperatorNotIn || o == FilterOperatorInLike || o == FilterOperatorNotInLike
}
