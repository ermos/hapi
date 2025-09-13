package hapi

import (
	"testing"
)

func TestFilterOperatorValid(t *testing.T) {
	tests := []struct {
		name     string
		operator FilterOperator
		wantErr  bool
	}{
		{"Valid Equal", FilterOperatorEqual, false},
		{"Valid Not Equal", FilterOperatorNotEqual, false},
		{"Valid Like", FilterOperatorLike, false},
		{"Valid Not Like", FilterOperatorNotLike, false},
		{"Valid In", FilterOperatorIn, false},
		{"Valid Not In", FilterOperatorNotIn, false},
		{"Valid In Like", FilterOperatorInLike, false},
		{"Valid Not In Like", FilterOperatorNotInLike, false},
		{"Valid Greater Than", FilterOperatorGreaterThan, false},
		{"Valid Less Than", FilterOperatorLessThan, false},
		{"Valid Greater Or Equal", FilterOperatorGreaterOrEqual, false},
		{"Valid Less Or Equal", FilterOperatorLessOrEqual, false},
		{"Invalid operator", FilterOperator("invalid"), true},
		{"Empty operator", FilterOperator(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operator.Valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterOperator.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterOperatorIsList(t *testing.T) {
	tests := []struct {
		name     string
		operator FilterOperator
		want     bool
	}{
		{"In operator", FilterOperatorIn, true},
		{"Not In operator", FilterOperatorNotIn, true},
		{"In Like operator", FilterOperatorInLike, true},
		{"Not In Like operator", FilterOperatorNotInLike, true},
		{"Equal operator", FilterOperatorEqual, false},
		{"Not Equal operator", FilterOperatorNotEqual, false},
		{"Like operator", FilterOperatorLike, false},
		{"Not Like operator", FilterOperatorNotLike, false},
		{"Greater Than operator", FilterOperatorGreaterThan, false},
		{"Less Than operator", FilterOperatorLessThan, false},
		{"Greater Or Equal operator", FilterOperatorGreaterOrEqual, false},
		{"Less Or Equal operator", FilterOperatorLessOrEqual, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.IsList(); got != tt.want {
				t.Errorf("FilterOperator.IsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
