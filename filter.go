package hapi

// Filters represents a collection of Filter conditions.
type Filters []Filter

// Filter represents a single filter condition with a field, operator, and values.
type Filter struct {
	Field    string         // The field name to filter on
	Operator FilterOperator // The comparison operator
	Values   Values         // The values to compare against
}

// GetFromField returns all filters that match the specified field name.
func (f Filters) GetFromField(field string) []Filter {
	if field == "" {
		return nil
	}

	// Pre-allocate with reasonable capacity to reduce allocations
	list := make([]Filter, 0, len(f)/4)
	for _, filter := range f {
		if filter.Field == field {
			list = append(list, filter)
		}
	}

	if len(list) == 0 {
		return nil
	}
	return list
}

// GetFromFields returns all filters that match any of the specified field names.
// The filters are returned in the order they appear for each field.
func (f Filters) GetFromFields(fields []string) []Filter {
	if len(fields) == 0 {
		return nil
	}

	// Pre-allocate with reasonable capacity
	list := make([]Filter, 0, len(f))

	// Maintain order by iterating through fields first
	for _, field := range fields {
		for _, filter := range f {
			if filter.Field == field {
				list = append(list, filter)
			}
		}
	}

	if len(list) == 0 {
		return nil
	}
	return list
}

// GetFirstFromField returns the first filter that matches the specified field name.
// Returns a zero-value Filter if no match is found.
func (f Filters) GetFirstFromField(field string) Filter {
	if field == "" {
		return Filter{}
	}

	for _, filter := range f {
		if filter.Field == field {
			return filter
		}
	}

	return Filter{}
}
