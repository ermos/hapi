package hapi

import (
	"fmt"
	"strings"
)

// Sort represents a sorting configuration with field and direction.
type Sort struct {
	Field     string        `json:"field"`     // The field to sort by
	Direction SortDirection `json:"direction"` // The sort direction (asc or desc)
}

// parseSortFromString parses a sort string in the format "field:direction".
func parseSortFromString(value string) (Sort, error) {
	if value == "" {
		return Sort{}, fmt.Errorf("sort value cannot be empty")
	}

	part := strings.Split(value, ":")
	if len(part) != 2 {
		return Sort{}, fmt.Errorf("invalid sort format: expected 'field:direction', got %q", value)
	}

	// Allow empty field for backward compatibility (though not recommended)
	direction := SortDirection(part[1])
	if err := direction.Valid(); err != nil {
		return Sort{}, err
	}

	return Sort{
		Field:     part[0],
		Direction: direction,
	}, nil
}
