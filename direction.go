package hapi

import "fmt"

// SortDirection represents the direction of sorting (ascending or descending).
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// Valid checks if the sort direction is valid.
// Returns an error if the direction is not recognized.
func (s SortDirection) Valid() error {
	switch s {
	case SortDirectionAsc, SortDirectionDesc:
		return nil
	}

	return fmt.Errorf("invalid sort direction: %q", s)
}
