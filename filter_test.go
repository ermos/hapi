package hapi

import (
	"reflect"
	"testing"
)

func TestFiltersGetFromField(t *testing.T) {
	filters := Filters{
		{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
		{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
		{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
		{Field: "status", Operator: FilterOperatorEqual, Values: Values{"active"}},
	}

	tests := []struct {
		name  string
		field string
		want  []Filter
	}{
		{
			name:  "Field with multiple matches",
			field: "name",
			want: []Filter{
				{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
			},
		},
		{
			name:  "Field with single match",
			field: "age",
			want: []Filter{
				{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
			},
		},
		{
			name:  "Field with no matches",
			field: "nonexistent",
			want:  nil,
		},
		{
			name:  "Empty field name",
			field: "",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.GetFromField(tt.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filters.GetFromField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiltersGetFromFields(t *testing.T) {
	filters := Filters{
		{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
		{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
		{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
		{Field: "status", Operator: FilterOperatorEqual, Values: Values{"active"}},
	}

	tests := []struct {
		name   string
		fields []string
		want   []Filter
	}{
		{
			name:   "Multiple fields with matches",
			fields: []string{"name", "age"},
			want: []Filter{
				{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
				{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
			},
		},
		{
			name:   "Single field",
			fields: []string{"status"},
			want: []Filter{
				{Field: "status", Operator: FilterOperatorEqual, Values: Values{"active"}},
			},
		},
		{
			name:   "Fields with no matches",
			fields: []string{"nonexistent1", "nonexistent2"},
			want:   nil,
		},
		{
			name:   "Empty fields slice",
			fields: []string{},
			want:   nil,
		},
		{
			name:   "Mix of existing and non-existing fields",
			fields: []string{"name", "nonexistent"},
			want: []Filter{
				{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.GetFromFields(tt.fields)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filters.GetFromFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiltersGetFirstFromField(t *testing.T) {
	filters := Filters{
		{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
		{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
		{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
		{Field: "status", Operator: FilterOperatorEqual, Values: Values{"active"}},
	}

	tests := []struct {
		name  string
		field string
		want  Filter
	}{
		{
			name:  "Field with multiple matches - returns first",
			field: "name",
			want:  Filter{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
		},
		{
			name:  "Field with single match",
			field: "age",
			want:  Filter{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"25"}},
		},
		{
			name:  "Field with no matches - returns zero value",
			field: "nonexistent",
			want:  Filter{},
		},
		{
			name:  "Empty field name",
			field: "",
			want:  Filter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filters.GetFirstFromField(tt.field)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filters.GetFirstFromField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiltersGetFirstFromFieldEmptyFilters(t *testing.T) {
	var emptyFilters Filters

	got := emptyFilters.GetFirstFromField("anyfield")
	want := Filter{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Empty Filters.GetFirstFromField() = %v, want %v", got, want)
	}
}
