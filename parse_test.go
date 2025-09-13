package hapi

import (
	"net/http"
	"reflect"
	"testing"
)

func TestParseFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    Result
		wantErr bool
	}{
		{
			name: "Simple filter",
			url:  "http://example.com/api/users?name=John",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple filters",
			url:  "http://example.com/api/users?name=John&age=25",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorEqual, Values: Values{"25"}},
				},
			},
			wantErr: false,
		},
		{
			name: "With pagination",
			url:  "http://example.com/api/users?per_page=10&page=3",
			want: Result{
				PerPage: 10,
				Page:    3,
			},
			wantErr: false,
		},
		{
			name: "With sort",
			url:  "http://example.com/api/users?sort=name:asc",
			want: Result{
				Sort: Sort{Field: "name", Direction: SortDirectionAsc},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			got, err := ParseFromRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    Result
		wantErr bool
	}{
		{
			name: "Empty URL",
			url:  "",
			want: Result{
				Filters: Filters{
					{Field: "", Operator: FilterOperatorEqual, Values: Values{""}},
				},
			},
			wantErr: false,
		},
		{
			name: "Basic filter",
			url:  "http://example.com?name=John",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Filter with operator",
			url:  "http://example.com?name[lk]=Jo%25",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
				},
			},
			wantErr: false,
		},
		{
			name: "List operator with comma-separated values",
			url:  "http://example.com?status[in]=active,inactive",
			want: Result{
				Filters: Filters{
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "inactive"}},
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple operators on same field",
			url:  "http://example.com?age[gt]=18&age[lt]=65",
			want: Result{
				Filters: Filters{
					{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"18"}},
					{Field: "age", Operator: FilterOperatorLessThan, Values: Values{"65"}},
				},
			},
			wantErr: false,
		},
		{
			name: "PerPage and page",
			url:  "http://example.com?per_page=50&page=5",
			want: Result{
				PerPage: 50,
				Page:    5,
			},
			wantErr: false,
		},
		{
			name: "Sort ascending",
			url:  "http://example.com?sort=name:asc",
			want: Result{
				Sort: Sort{Field: "name", Direction: SortDirectionAsc},
			},
			wantErr: false,
		},
		{
			name: "Sort descending",
			url:  "http://example.com?sort=created_at:desc",
			want: Result{
				Sort: Sort{Field: "created_at", Direction: SortDirectionDesc},
			},
			wantErr: false,
		},
		{
			name: "Complex query",
			url:  "http://example.com?name[lk]=John%25&age[ge]=18&status[in]=active,pending&per_page=25&page=3&sort=created_at:desc",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"John%"}},
					{Field: "age", Operator: FilterOperatorGreaterOrEqual, Values: Values{"18"}},
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "pending"}},
				},
				Sort:    Sort{Field: "created_at", Direction: SortDirectionDesc},
				PerPage: 25,
				Page:    3,
			},
			wantErr: false,
		},
		{
			name: "Filter without value",
			url:  "http://example.com?name",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{""}},
				},
			},
			wantErr: false,
		},
		{
			name: "Zero and negative per_page/page",
			url:  "http://example.com?per_page=-5&page=-10",
			want: Result{
				PerPage: 1,
				Page:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStrict(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    Result
		wantErr bool
	}{
		{
			name: "Valid query",
			url:  "http://example.com?name=John&age=25",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorEqual, Values: Values{"25"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid per_page format",
			url:     "http://example.com?per_page",
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "Invalid page format",
			url:     "http://example.com?page",
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "Invalid sort format",
			url:     "http://example.com?sort",
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "Invalid sort value",
			url:     "http://example.com?sort=name:invalid",
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "Invalid operator",
			url:     "http://example.com?name[invalid]=value",
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "Invalid URL encoding in strict mode",
			url:     "http://example.com?name=%",
			want:    Result{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStrict(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStrict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFromRequestStrict(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    Result
		wantErr bool
	}{
		{
			name: "Valid request",
			url:  "http://example.com?name=John",
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid request - malformed per_page",
			url:     "http://example.com?per_page",
			want:    Result{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			got, err := ParseFromRequestStrict(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromRequestStrict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromRequestStrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWithAllOperators(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		operator FilterOperator
		value    string
		expected Values
	}{
		{"Equal operator", "http://example.com?name[eq]=John", FilterOperatorEqual, "John", Values{"John"}},
		{"Not equal operator", "http://example.com?name[ne]=John", FilterOperatorNotEqual, "John", Values{"John"}},
		{"Like operator", "http://example.com?name[lk]=John%25", FilterOperatorLike, "John%", Values{"John%"}},
		{"Not like operator", "http://example.com?name[nlk]=John%25", FilterOperatorNotLike, "John%", Values{"John%"}},
		{"In operator", "http://example.com?status[in]=active,inactive", FilterOperatorIn, "active,inactive", Values{"active", "inactive"}},
		{"Not in operator", "http://example.com?status[nin]=active,inactive", FilterOperatorNotIn, "active,inactive", Values{"active", "inactive"}},
		{"In like operator", "http://example.com?name[inlk]=John,Jane", FilterOperatorInLike, "John,Jane", Values{"John", "Jane"}},
		{"Not in like operator", "http://example.com?name[ninlk]=John,Jane", FilterOperatorNotInLike, "John,Jane", Values{"John", "Jane"}},
		{"Greater than operator", "http://example.com?age[gt]=18", FilterOperatorGreaterThan, "18", Values{"18"}},
		{"Less than operator", "http://example.com?age[lt]=65", FilterOperatorLessThan, "65", Values{"65"}},
		{"Greater or equal operator", "http://example.com?age[ge]=18", FilterOperatorGreaterOrEqual, "18", Values{"18"}},
		{"Less or equal operator", "http://example.com?age[le]=65", FilterOperatorLessOrEqual, "65", Values{"65"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.url)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if len(result.Filters) != 1 {
				t.Errorf("Expected 1 filter, got %d", len(result.Filters))
				return
			}

			filter := result.Filters[0]
			if filter.Operator != tt.operator {
				t.Errorf("Expected operator %v, got %v", tt.operator, filter.Operator)
			}

			if !reflect.DeepEqual(filter.Values, tt.expected) {
				t.Errorf("Expected values %v, got %v", tt.expected, filter.Values)
			}
		})
	}
}
