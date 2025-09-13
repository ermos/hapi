package hapi

import (
	"errors"
	"strings"
)

type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

func parseSortFromString(value string) (Sort, error) {
	part := strings.Split(value, ":")
	if len(part) != 2 {
		return Sort{}, errors.New("invalid sort")
	}

	direction := SortDirection(part[1])
	if err := direction.Valid(); err != nil {
		return Sort{}, err
	}

	return Sort{
		Field:     part[0],
		Direction: direction,
	}, nil
}
