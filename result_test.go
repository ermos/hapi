package hapi

import (
	"reflect"
	"testing"
)

func TestResultStruct(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   Result
	}{
		{
			name:   "Empty result",
			result: Result{},
			want: Result{
				Filters: nil,
				Sort:    Sort{},
				Limit:   0,
				Offset:  0,
			},
		},
		{
			name: "Result with filters only",
			result: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
				},
			},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
				},
				Sort:   Sort{},
				Limit:  0,
				Offset: 0,
			},
		},
		{
			name: "Result with sort only",
			result: Result{
				Sort: Sort{Field: "name", Direction: SortDirectionAsc},
			},
			want: Result{
				Filters: nil,
				Sort:    Sort{Field: "name", Direction: SortDirectionAsc},
				Limit:   0,
				Offset:  0,
			},
		},
		{
			name: "Result with pagination only",
			result: Result{
				Limit:  50,
				Offset: 100,
			},
			want: Result{
				Filters: nil,
				Sort:    Sort{},
				Limit:   50,
				Offset:  100,
			},
		},
		{
			name: "Complete result",
			result: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"John%"}},
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "pending"}},
				},
				Sort:   Sort{Field: "created_at", Direction: SortDirectionDesc},
				Limit:  25,
				Offset: 50,
			},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"John%"}},
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "pending"}},
				},
				Sort:   Sort{Field: "created_at", Direction: SortDirectionDesc},
				Limit:  25,
				Offset: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.result, tt.want) {
				t.Errorf("Result struct = %v, want %v", tt.result, tt.want)
			}
		})
	}
}

func TestResultZeroValues(t *testing.T) {
	var result Result

	if result.Filters != nil {
		t.Errorf("Expected nil Filters, got %v", result.Filters)
	}

	expectedSort := Sort{}
	if !reflect.DeepEqual(result.Sort, expectedSort) {
		t.Errorf("Expected empty Sort %v, got %v", expectedSort, result.Sort)
	}

	if result.Limit != 0 {
		t.Errorf("Expected Limit 0, got %d", result.Limit)
	}

	if result.Offset != 0 {
		t.Errorf("Expected Offset 0, got %d", result.Offset)
	}
}

func TestResultFieldAssignment(t *testing.T) {
	result := Result{}

	result.Filters = Filters{
		{Field: "test", Operator: FilterOperatorEqual, Values: Values{"value"}},
	}
	result.Sort = Sort{Field: "test", Direction: SortDirectionAsc}
	result.Limit = 10
	result.Offset = 20

	expectedFilters := Filters{
		{Field: "test", Operator: FilterOperatorEqual, Values: Values{"value"}},
	}
	if !reflect.DeepEqual(result.Filters, expectedFilters) {
		t.Errorf("Filters assignment failed: got %v, want %v", result.Filters, expectedFilters)
	}

	expectedSort := Sort{Field: "test", Direction: SortDirectionAsc}
	if !reflect.DeepEqual(result.Sort, expectedSort) {
		t.Errorf("Sort assignment failed: got %v, want %v", result.Sort, expectedSort)
	}

	if result.Limit != 10 {
		t.Errorf("Limit assignment failed: got %d, want %d", result.Limit, 10)
	}

	if result.Offset != 20 {
		t.Errorf("Offset assignment failed: got %d, want %d", result.Offset, 20)
	}
}
