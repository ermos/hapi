package hapi_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ermos/hapi"
)

func ExampleParse() {
	// Parse a URL with various query parameters
	url := "http://api.example.com/users?name[lk]=John%25&age[ge]=18&status[in]=active,pending&page=25&per_page=50&sort=created_at:desc"

	opts := hapi.Options{
		DefaultPerPage: 10,
		MaxPerPage:     100,
	}

	result, err := hapi.Parse(url, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Filters: %d\n", len(result.Filters))
	if len(result.Sorts) > 0 {
		fmt.Printf("Sort: %s %s\n", result.Sorts[0].Field, result.Sorts[0].Direction)
	}
	fmt.Printf("Page: %d, PerPage: %d\n", result.Page, result.PerPage)

	// Output:
	// Filters: 3
	// Sort: created_at desc
	// Page: 25, PerPage: 50
}

func ExampleParseFromRequest() {
	// In a real HTTP handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		opts := hapi.Options{
			DefaultPerPage: 10,
			MaxPerPage:     100,
		}

		result, err := hapi.ParseFromRequest(r, opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use the parsed result for database queries
		filters := result.Filters.GetFromField("name")
		if len(filters) > 0 {
			// Apply name filter to query
			fmt.Printf("Name filter: %v\n", filters[0].Values.First())
		}
	}

	// Example request
	req, _ := http.NewRequest("GET", "http://api.example.com/users?name=John", nil)
	handler(nil, req)

	// Output:
	// Name filter: John
}

func ExampleFilters_GetFromFields() {
	filters := hapi.Filters{
		{Field: "name", Operator: hapi.FilterOperatorEqual, Values: hapi.Values{"John"}},
		{Field: "age", Operator: hapi.FilterOperatorGreaterThan, Values: hapi.Values{"25"}},
		{Field: "name", Operator: hapi.FilterOperatorLike, Values: hapi.Values{"Jo%"}},
		{Field: "status", Operator: hapi.FilterOperatorEqual, Values: hapi.Values{"active"}},
	}

	// Get filters for specific fields
	nameAndAge := filters.GetFromFields([]string{"name", "age"})

	for _, filter := range nameAndAge {
		fmt.Printf("%s %s %v\n", filter.Field, filter.Operator, filter.Values.First())
	}

	// Output:
	// name eq John
	// name lk Jo%
	// age gt 25
}

func ExampleValue_conversion() {
	v := hapi.Value("123")

	fmt.Printf("String: %s\n", v.String())
	fmt.Printf("Int: %d\n", v.Int())
	fmt.Printf("Int64: %d\n", v.Int64())
	fmt.Printf("Float64: %f\n", v.Float64())

	// Output:
	// String: 123
	// Int: 123
	// Int64: 123
	// Float64: 123.000000
}

func ExampleNewOptions() {
	// Create options with custom configuration
	opts := hapi.NewOptions(
		hapi.WithDefaultPerPage(25),
		hapi.WithMaxPerPage(200),
		hapi.WithAllowedSorts([]string{"name", "created_at", "updated_at"}),
		hapi.WithAllowedFilters([]string{"status", "type", "name"}),
	)

	url := "http://api.example.com/users?name=John&sort=created_at:desc&per_page=50"
	result, err := hapi.Parse(url, *opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PerPage: %d (default was %d, requested 50)\n", result.PerPage, opts.DefaultPerPage)
	fmt.Printf("Allowed sorts: %v\n", opts.AllowedSorts)
	fmt.Printf("Allowed filters: %v\n", opts.AllowedFilters)

	// Output:
	// PerPage: 50 (default was 25, requested 50)
	// Allowed sorts: [name created_at updated_at]
	// Allowed filters: [status type name]
}

func ExampleParseStrict() {
	// Use strict mode with allowed fields
	opts := hapi.NewOptions(
		hapi.WithDefaultPerPage(10),
		hapi.WithMaxPerPage(100),
		hapi.WithAllowedSorts([]string{"name", "age"}),
		hapi.WithAllowedFilters([]string{"name", "status"}),
	)

	// This URL contains a disallowed filter field "salary"
	url := "http://api.example.com/users?name=John&salary=100000"

	_, err := hapi.ParseStrict(url, *opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Error: filtering by field "salary" is not allowed
}

func ExampleParse_multipleSorts() {
	// Parse a URL with multiple sorts using comma separation
	url := "http://api.example.com/users?sort=name:asc,age:desc,created_at:asc"

	opts := hapi.Options{
		DefaultPerPage: 10,
		MaxPerPage:     100,
	}

	result, err := hapi.Parse(url, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of sorts: %d\n", len(result.Sorts))
	for i, sort := range result.Sorts {
		fmt.Printf("Sort %d: %s %s\n", i+1, sort.Field, sort.Direction)
	}

	// Output:
	// Number of sorts: 3
	// Sort 1: name asc
	// Sort 2: age desc
	// Sort 3: created_at asc
}
