package hapi

import (
	"testing"
)

func TestSortDirectionValid(t *testing.T) {
	tests := []struct {
		name      string
		direction SortDirection
		wantErr   bool
	}{
		{"Valid ascending", SortDirectionAsc, false},
		{"Valid descending", SortDirectionDesc, false},
		{"Invalid direction", SortDirection("invalid"), true},
		{"Empty direction", SortDirection(""), true},
		{"Case sensitive - ASC", SortDirection("ASC"), true},
		{"Case sensitive - DESC", SortDirection("DESC"), true},
		{"Partial match", SortDirection("as"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.direction.Valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("SortDirection.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
