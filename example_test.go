package hapi_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ermos/hapi"
)

func ExampleParse() {
	// Parse a URL with various query parameters
	url := "http://api.example.com/users?name[lk]=John%25&age[ge]=18&status[in]=active,pending&limit=25&offset=50&sort=created_at:desc"

	result, err := hapi.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Filters: %d\n", len(result.Filters))
	fmt.Printf("Sort: %s %s\n", result.Sort.Field, result.Sort.Direction)
	fmt.Printf("Limit: %d, Offset: %d\n", result.Limit, result.Offset)

	// Output:
	// Filters: 3
	// Sort: created_at desc
	// Limit: 25, Offset: 50
}

func ExampleParseFromRequest() {
	// In a real HTTP handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		result, err := hapi.ParseFromRequest(r)
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

