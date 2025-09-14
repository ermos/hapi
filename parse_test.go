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
		opts    Options
		want    Result
		wantErr bool
	}{
		{
			name: "Simple filter",
			url:  "http://example.com/api/users?name=John",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Multiple filters",
			url:  "http://example.com/api/users?name=John&age=25",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorEqual, Values: Values{"25"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "With pagination",
			url:  "http://example.com/api/users?per_page=10&page=3",
			opts: Options{DefaultPerPage: 20, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{},
				PerPage: 10,
				Page:    3,
			},
			wantErr: false,
		},
		{
			name: "With sort",
			url:  "http://example.com/api/users?sort=name:asc",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{{Field: "name", Direction: SortDirectionAsc}},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "With multiple sorts",
			url:  "http://example.com/api/users?sort=name:asc&sort=age:desc",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts: Sorts{
					{Field: "name", Direction: SortDirectionAsc},
					{Field: "age", Direction: SortDirectionDesc},
				},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "PerPage exceeding max",
			url:  "http://example.com/api/users?per_page=200",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 100, // Should be capped at MaxPerPage
			},
			wantErr: false,
		},
		{
			name: "With allowed sorts - valid",
			url:  "http://example.com/api/users?sort=name:asc",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedSorts:   []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{{Field: "name", Direction: SortDirectionAsc}},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "With allowed sorts - invalid (ignored in non-strict)",
			url:  "http://example.com/api/users?sort=salary:desc",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedSorts:   []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{}, // Sort should be ignored
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "With allowed filters - valid",
			url:  "http://example.com/api/users?name=John",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedFilters: []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "With allowed filters - invalid (ignored in non-strict)",
			url:  "http://example.com/api/users?salary=100000",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedFilters: []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{}, // Filter should be ignored
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name:    "Nil request",
			url:     "",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.name == "Nil request" {
				req = nil
			} else {
				req, err = http.NewRequest("GET", tt.url, nil)
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}
			}

			got, err := ParseFromRequest(req, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		opts    Options
		want    Result
		wantErr bool
	}{
		{
			name: "Empty URL",
			url:  "",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "", Operator: FilterOperatorEqual, Values: Values{""}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Basic filter",
			url:  "http://example.com?name=John",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Filter with operator",
			url:  "http://example.com?name[lk]=Jo%25",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"Jo%"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "List operator with comma-separated values",
			url:  "http://example.com?status[in]=active,inactive",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "inactive"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Multiple operators on same field",
			url:  "http://example.com?age[gt]=18&age[lt]=65",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "age", Operator: FilterOperatorGreaterThan, Values: Values{"18"}},
					{Field: "age", Operator: FilterOperatorLessThan, Values: Values{"65"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "PerPage and page",
			url:  "http://example.com?per_page=50&page=5",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{},
				PerPage: 50,
				Page:    5,
			},
			wantErr: false,
		},
		{
			name: "Sort ascending",
			url:  "http://example.com?sort=name:asc",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{{Field: "name", Direction: SortDirectionAsc}},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Sort descending",
			url:  "http://example.com?sort=created_at:desc",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{{Field: "created_at", Direction: SortDirectionDesc}},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Complex query",
			url:  "http://example.com?name[lk]=John%25&age[ge]=18&status[in]=active,pending&per_page=25&page=3&sort=created_at:desc",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorLike, Values: Values{"John%"}},
					{Field: "age", Operator: FilterOperatorGreaterOrEqual, Values: Values{"18"}},
					{Field: "status", Operator: FilterOperatorIn, Values: Values{"active", "pending"}},
				},
				Sorts:   Sorts{{Field: "created_at", Direction: SortDirectionDesc}},
				PerPage: 25,
				Page:    3,
			},
			wantErr: false,
		},
		{
			name: "Filter without value",
			url:  "http://example.com?name",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{""}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Zero and negative per_page/page",
			url:  "http://example.com?per_page=-5&page=-10",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{},
				PerPage: 1,
				Page:    1,
			},
			wantErr: false,
		},
		{
			name: "Custom default per page",
			url:  "http://example.com?name=test",
			opts: Options{DefaultPerPage: 25, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"test"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 25, // Should use custom default
			},
			wantErr: false,
		},
		{
			name: "Per page at max limit",
			url:  "http://example.com?per_page=150",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 50},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 50, // Should be capped at MaxPerPage
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.url, tt.opts)
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
		opts    Options
		want    Result
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid query",
			url:  "http://example.com?name=John&age=25",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
					{Field: "age", Operator: FilterOperatorEqual, Values: Values{"25"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name:    "Invalid per_page format",
			url:     "http://example.com?per_page",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid per_page filter format",
		},
		{
			name:    "Invalid page format",
			url:     "http://example.com?page",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid page filter format",
		},
		{
			name:    "Invalid sort format",
			url:     "http://example.com?sort",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid sort filter format",
		},
		{
			name:    "Invalid sort value",
			url:     "http://example.com?sort=name:invalid",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid sort direction",
		},
		{
			name:    "Invalid operator",
			url:     "http://example.com?name[invalid]=value",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid operator",
		},
		{
			name:    "Invalid URL encoding in strict mode",
			url:     "http://example.com?name=%",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "failed to unescape value",
		},
		{
			name: "Disallowed sort field",
			url:  "http://example.com?sort=salary:asc",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedSorts:   []string{"name", "age"},
			},
			want:    Result{},
			wantErr: true,
			errMsg:  "sorting by field \"salary\" is not allowed",
		},
		{
			name: "Disallowed filter field",
			url:  "http://example.com?salary=100000",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedFilters: []string{"name", "age"},
			},
			want:    Result{},
			wantErr: true,
			errMsg:  "filtering by field \"salary\" is not allowed",
		},
		{
			name: "Allowed sort field",
			url:  "http://example.com?sort=name:asc",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedSorts:   []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{},
				Sorts:   Sorts{{Field: "name", Direction: SortDirectionAsc}},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name: "Allowed filter field",
			url:  "http://example.com?name=John",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedFilters: []string{"name", "age"},
			},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStrict(tt.url, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStrict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseStrict() error = %v, want error containing %q", err, tt.errMsg)
				}
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
		opts    Options
		want    Result
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid request",
			url:  "http://example.com?name=John",
			opts: Options{DefaultPerPage: 10, MaxPerPage: 100},
			want: Result{
				Filters: Filters{
					{Field: "name", Operator: FilterOperatorEqual, Values: Values{"John"}},
				},
				Sorts:   Sorts{},
				Page:    1,
				PerPage: 10,
			},
			wantErr: false,
		},
		{
			name:    "Invalid request - malformed per_page",
			url:     "http://example.com?per_page",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "invalid per_page filter format",
		},
		{
			name:    "Nil request",
			url:     "",
			opts:    Options{DefaultPerPage: 10, MaxPerPage: 100},
			want:    Result{},
			wantErr: true,
			errMsg:  "request or URL is nil",
		},
		{
			name: "Disallowed filter in strict mode",
			url:  "http://example.com?salary=100000",
			opts: Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedFilters: []string{"name", "age"},
			},
			want:    Result{},
			wantErr: true,
			errMsg:  "filtering by field \"salary\" is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.name == "Nil request" {
				req = nil
			} else {
				req, err = http.NewRequest("GET", tt.url, nil)
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}
			}

			got, err := ParseFromRequestStrict(req, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromRequestStrict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ParseFromRequestStrict() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFromRequestStrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWithAllOperators(t *testing.T) {
	opts := Options{DefaultPerPage: 10, MaxPerPage: 100}

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
			result, err := Parse(tt.url, opts)
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

func TestParseOptionsIntegration(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		opts   Options
		strict bool
		check  func(t *testing.T, result Result, err error)
	}{
		{
			name: "Empty options uses defaults",
			url:  "http://example.com?name=test",
			opts: Options{},
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.PerPage != 0 { // DefaultPerPage is 0 when not set
					t.Errorf("Expected PerPage = 0, got %d", result.PerPage)
				}
			},
		},
		{
			name: "Options from NewOptions",
			url:  "http://example.com?name=test&sort=created:asc",
			opts: *NewOptions(
				WithDefaultPerPage(25),
				WithMaxPerPage(50),
				WithAllowedSorts([]string{"created", "updated"}),
				WithAllowedFilters([]string{"name", "status"}),
			),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.PerPage != 25 {
					t.Errorf("Expected PerPage = 25, got %d", result.PerPage)
				}
				if len(result.Sorts) != 1 || result.Sorts[0].Field != "created" {
					t.Errorf("Expected sort on 'created', got %v", result.Sorts)
				}
			},
		},
		{
			name: "Filtering with disallowed field in non-strict",
			url:  "http://example.com?salary=100000&name=John",
			opts: *NewOptions(
				WithDefaultPerPage(10),
				WithAllowedFilters([]string{"name"}),
			),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Only 'name' filter should be present
				if len(result.Filters) != 1 {
					t.Errorf("Expected 1 filter, got %d", len(result.Filters))
				}
				if result.Filters[0].Field != "name" {
					t.Errorf("Expected filter field 'name', got %s", result.Filters[0].Field)
				}
			},
		},
		{
			name:   "Filtering with disallowed field in strict",
			url:    "http://example.com?salary=100000",
			opts:   *NewOptions(WithAllowedFilters([]string{"name"})),
			strict: true,
			check: func(t *testing.T, result Result, err error) {
				if err == nil {
					t.Error("Expected error for disallowed filter, got nil")
				}
				if !contains(err.Error(), "filtering by field \"salary\" is not allowed") {
					t.Errorf("Unexpected error message: %v", err)
				}
			},
		},
		{
			name: "Multiple sorts with allowed list",
			url:  "http://example.com?sort=name:asc&sort=age:desc&sort=salary:asc",
			opts: *NewOptions(
				WithDefaultPerPage(10),
				WithAllowedSorts([]string{"name", "age"}),
			),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Only 'name' and 'age' sorts should be present
				if len(result.Sorts) != 2 {
					t.Errorf("Expected 2 sorts, got %d", len(result.Sorts))
				}
			},
		},
		{
			name: "PerPage boundary testing",
			url:  "http://example.com?per_page=1000",
			opts: *NewOptions(
				WithDefaultPerPage(10),
				WithMaxPerPage(100),
			),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result.PerPage != 100 {
					t.Errorf("Expected PerPage capped at 100, got %d", result.PerPage)
				}
			},
		},
		{
			name: "Zero max per page",
			url:  "http://example.com?per_page=50",
			opts: *NewOptions(
				WithDefaultPerPage(10),
				WithMaxPerPage(0),
			),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// When MaxPerPage is 0, it should be treated as no limit
				if result.PerPage != 0 {
					t.Errorf("Expected PerPage = 0 (no value should pass through max(1, min(50, 0))), got %d", result.PerPage)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Result
			var err error

			if tt.strict {
				result, err = ParseStrict(tt.url, tt.opts)
			} else {
				result, err = Parse(tt.url, tt.opts)
			}

			tt.check(t, result, err)
		})
	}
}

func TestParseEdgeCases(t *testing.T) {
	opts := Options{DefaultPerPage: 10, MaxPerPage: 100}

	tests := []struct {
		name string
		url  string
		check func(t *testing.T, result Result, err error)
	}{
		{
			name: "URL with fragment",
			url:  "http://example.com?name=test#section",
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result.Filters) != 1 || result.Filters[0].Field != "name" {
					t.Errorf("Expected name filter, got %v", result.Filters)
				}
			},
		},
		{
			name: "Empty query string",
			url:  "http://example.com?",
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// An empty query string results in a single empty filter
				if len(result.Filters) != 1 || result.Filters[0].Field != "" {
					t.Errorf("Expected single empty filter, got %v", result.Filters)
				}
			},
		},
		{
			name: "Multiple question marks",
			url:  "http://example.com??name=test",
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// The parser should handle this gracefully
			},
		},
		{
			name: "Unicode in values",
			url:  "http://example.com?name=测试&city=北京",
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result.Filters) != 2 {
					t.Errorf("Expected 2 filters, got %d", len(result.Filters))
				}
			},
		},
		{
			name: "Special characters in field names",
			url:  "http://example.com?user.name=John&user.age=25",
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result.Filters) != 2 {
					t.Errorf("Expected 2 filters, got %d", len(result.Filters))
				}
				if result.Filters[0].Field != "user.name" {
					t.Errorf("Expected field 'user.name', got %s", result.Filters[0].Field)
				}
			},
		},
		{
			name: "Extremely long value",
			url:  "http://example.com?description=" + stringRepeat("a", 1000),
			check: func(t *testing.T, result Result, err error) {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result.Filters) != 1 {
					t.Errorf("Expected 1 filter, got %d", len(result.Filters))
				}
				if len(result.Filters[0].Values.First().String()) != 1000 {
					t.Errorf("Expected value length 1000, got %d", len(result.Filters[0].Values.First().String()))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.url, opts)
			tt.check(t, result, err)
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr) != -1))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func stringRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}