# 🚀 HAPI - HTTP API Query Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/ermos/hapi.svg)](https://pkg.go.dev/github.com/ermos/hapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermos/hapi)](https://goreportcard.com/report/github.com/ermos/hapi)

A clean, elegant Go library for parsing HTTP API query parameters with support for filtering, sorting, and pagination. Transform complex query strings into structured data effortlessly.

## ✨ Features

- **Rich Filtering**: Support for 12 different operators (`eq`, `ne`, `gt`, `lt`, `ge`, `le`, `lk`, `nlk`, `in`, `nin`, `inlk`, `ninlk`)
- **Flexible Configuration**: Options system with validation and field restrictions
- **Sort Support**: Parse single and multiple sort parameters with directions (`asc`, `desc`)
- **Pagination Support**: Built-in per page and page handling with configurable limits
- **Type Conversion**: Automatic conversion to common Go types (int, float64, string)
- **Strict Mode**: Optional strict parsing with comprehensive error handling
- **Zero Dependencies**: Pure Go implementation with only standard library

## 📦 Installation

```bash
go get github.com/ermos/hapi
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/ermos/hapi"
)

func main() {
    url := "http://api.example.com/users?name[lk]=John%25&age[ge]=18&status[in]=active,pending&page=25&per_page=50&sort=created_at:desc,name:asc"

    opts := hapi.Options{
        DefaultPerPage: 10,
        MaxPerPage:     100,
    }

    result, err := hapi.Parse(url, opts)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d filters\n", len(result.Filters))
    fmt.Printf("Found %d sorts\n", len(result.Sorts))
    if len(result.Sorts) > 0 {
        fmt.Printf("Primary sort: %s %s\n", result.Sorts[0].Field, result.Sorts[0].Direction)
    }
    fmt.Printf("Pagination: page=%d, per_page=%d\n", result.Page, result.PerPage)
}
```

## 📖 Usage Examples

### Basic Filtering

```go
// URL: /users?name=John&age[gt]=25&status[in]=active,pending

opts := hapi.Options{DefaultPerPage: 20, MaxPerPage: 100}
result, _ := hapi.Parse(url, opts)

// Access filters by field
nameFilters := result.Filters.GetFromField("name")
if len(nameFilters) > 0 {
    fmt.Println("Name:", nameFilters[0].Values.First().String())
}

// Get first filter for a field
ageFilter := result.Filters.GetFirstFromField("age")
fmt.Printf("Age filter: %s %v\n", ageFilter.Operator, ageFilter.Values.First().Int())
```

### HTTP Request Parsing

```go
func handler(w http.ResponseWriter, r *http.Request) {
    opts := hapi.Options{
        DefaultPerPage: 20,
        MaxPerPage:     100,
    }

    result, err := hapi.ParseFromRequest(r, opts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Use parsed result for database queries
    applyFiltersToQuery(result.Filters)
    applySortingToQuery(result.Sorts)
    applyPaginationToQuery(result.Page, result.PerPage)
}
```

### Strict Mode

```go
// Use strict mode for validation
opts := hapi.Options{
    DefaultPerPage: 10,
    MaxPerPage:     50,
    AllowedSorts:   []string{"name", "created_at"},
    AllowedFilters: []string{"name", "status", "age"},
}

result, err := hapi.ParseStrict(url, opts)
if err != nil {
    // Handle validation errors
    fmt.Printf("Invalid query parameters: %v\n", err)
}
```

## 🔧 Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `eq` | Equal (default) | `name=John` or `name[eq]=John` |
| `ne` | Not equal | `status[ne]=inactive` |
| `gt` | Greater than | `age[gt]=18` |
| `lt` | Less than | `price[lt]=100` |
| `ge` | Greater or equal | `rating[ge]=4` |
| `le` | Less or equal | `score[le]=95` |
| `lk` | Like (substring) | `name[lk]=Jo%` |
| `nlk` | Not like | `email[nlk]=%spam%` |
| `in` | In list | `status[in]=active,pending` |
| `nin` | Not in list | `role[nin]=admin,super` |
| `inlk` | In like (any match) | `tags[inlk]=tech,go` |
| `ninlk` | Not in like | `categories[ninlk]=old,deprecated` |

## 📊 Query Structure

### Filters
```
field[operator]=value
field[operator]=value1,value2  // for list operators
```

### Sorting
```
# Basic sorting
sort=field:direction
sort=name:asc
sort=created_at:desc

# Multiple sorts in one parameter
sort=name:asc,age:desc,created_at:desc

# Multiple sorts with multiple parameters
sort=name:asc&sort=age:desc&sort=created_at:desc

# Mixed approach
sort=name:asc,age:desc&sort=status:asc
```

### Pagination
```
page=25
per_page=50
```

## ⚙️ Configuration

### Options System

HAPI provides a flexible options system for configuring parsing behavior:

```go
// Create options with defaults
opts := hapi.Options{
    DefaultPerPage: 20,   // Default items per page
    MaxPerPage:     100,  // Maximum allowed items per page
}

// Or use the builder pattern
opts := hapi.NewOptions(
    hapi.WithDefaultPerPage(25),
    hapi.WithMaxPerPage(200),
    hapi.WithAllowedSorts([]string{"name", "created_at", "updated_at"}),
    hapi.WithAllowedFilters([]string{"status", "type", "name"}),
)
```

### Field Validation

Restrict which fields can be used for sorting and filtering:

```go
opts := hapi.NewOptions(
    hapi.WithAllowedSorts([]string{"name", "age", "created_at"}),
    hapi.WithAllowedFilters([]string{"status", "role", "department"}),
)

// In non-strict mode: disallowed fields are ignored
result, _ := hapi.Parse(url, *opts)

// In strict mode: disallowed fields cause errors
result, err := hapi.ParseStrict(url, *opts)
if err != nil {
    // Handle validation error
}
```

### Complete Example

Here's a comprehensive example combining all features:

```go
package main

import (
    "fmt"
    "log"
    "github.com/ermos/hapi"
)

func main() {
    // Complex query with multiple filters, sorts, and pagination
    url := "http://api.example.com/users?name[lk]=John%&age[ge]=18&status[in]=active,pending&department=engineering&sort=name:asc,created_at:desc&page=2&per_page=25"

    // Configure options with validation
    opts := hapi.NewOptions(
        hapi.WithDefaultPerPage(20),
        hapi.WithMaxPerPage(100),
        hapi.WithAllowedSorts([]string{"name", "created_at", "age"}),
        hapi.WithAllowedFilters([]string{"name", "age", "status", "department"}),
    )

    // Parse in strict mode for validation
    result, err := hapi.ParseStrict(url, *opts)
    if err != nil {
        log.Fatal("Validation failed:", err)
    }

    // Display results
    fmt.Printf("=== Parsing Results ===\n")
    fmt.Printf("Filters: %d\n", len(result.Filters))
    for _, filter := range result.Filters {
        fmt.Printf("  %s [%s] = %v\n", filter.Field, filter.Operator, filter.Values.First())
    }

    fmt.Printf("Sorts: %d\n", len(result.Sorts))
    for i, sort := range result.Sorts {
        fmt.Printf("  %d. %s %s\n", i+1, sort.Field, sort.Direction)
    }

    fmt.Printf("Pagination: page=%d, per_page=%d\n", result.Page, result.PerPage)

    // Use with database queries
    query := buildDatabaseQuery(result)
    fmt.Printf("Generated SQL: %s\n", query)
}

func buildDatabaseQuery(result hapi.Result) string {
    // Example of how you might use the parsed result
    query := "SELECT * FROM users WHERE 1=1"

    // Add filters
    for _, filter := range result.Filters {
        switch filter.Operator {
        case hapi.FilterOperatorLike:
            query += fmt.Sprintf(" AND %s LIKE '%s'", filter.Field, filter.Values.First())
        case hapi.FilterOperatorGreaterOrEqual:
            query += fmt.Sprintf(" AND %s >= %d", filter.Field, filter.Values.First().Int())
        case hapi.FilterOperatorIn:
            query += fmt.Sprintf(" AND %s IN (%s)", filter.Field, filter.Values.ToString())
        default:
            query += fmt.Sprintf(" AND %s = '%s'", filter.Field, filter.Values.First())
        }
    }

    // Add sorting
    if len(result.Sorts) > 0 {
        query += " ORDER BY"
        for i, sort := range result.Sorts {
            if i > 0 {
                query += ","
            }
            query += fmt.Sprintf(" %s %s", sort.Field, sort.Direction)
        }
    }

    // Add pagination
    offset := (result.Page - 1) * result.PerPage
    query += fmt.Sprintf(" LIMIT %d OFFSET %d", result.PerPage, offset)

    return query
}
```

## 🏗️ API Reference

### Core Functions

```go
// Parse URL string (lenient mode)
func Parse(url string, opts Options) (Result, error)

// Parse URL string (strict mode)
func ParseStrict(url string, opts Options) (Result, error)

// Parse from HTTP request (lenient mode)
func ParseFromRequest(r *http.Request, opts Options) (Result, error)

// Parse from HTTP request (strict mode)
func ParseFromRequestStrict(r *http.Request, opts Options) (Result, error)
```

### Result Structure

```go
type Result struct {
    Filters Filters // Collection of filter conditions
    Sorts   Sorts   // Collection of sort configurations
    Page    int     // Current page number (1-based)
    PerPage int     // Number of items per page
}

type Sort struct {
    Field     string        // The field to sort by
    Direction SortDirection // The sort direction (asc or desc)
}

type Options struct {
    DefaultPerPage int      // Default number of items per page
    MaxPerPage     int      // Maximum allowed items per page
    AllowedSorts   []string // Allowed fields for sorting (empty = all allowed)
    AllowedFilters []string // Allowed fields for filtering (empty = all allowed)
}
```

### Filter Operations

```go
// Get all filters for a field
filters := result.Filters.GetFromField("name")

// Get first filter for a field
filter := result.Filters.GetFirstFromField("age")

// Get filters for multiple fields
filters := result.Filters.GetFromFields([]string{"name", "age"})
```

### Value Conversion

```go
value := filter.Values.First()

str := value.String()
num := value.Int()
bigNum := value.Int64()
decimal := value.Float64()
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
