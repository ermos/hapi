package hapi

import "fmt"

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

	return fmt.Errorf("invalid operator: %s", o)
}

func (o FilterOperator) IsList() bool {
	return o == FilterOperatorIn || o == FilterOperatorNotIn || o == FilterOperatorInLike || o == FilterOperatorNotInLike
}
