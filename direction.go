package hapi

import "errors"

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

func (s SortDirection) Valid() error {
	switch s {
	case SortDirectionAsc, SortDirectionDesc:
		return nil
	}

	return errors.New("invalid sort")
}
