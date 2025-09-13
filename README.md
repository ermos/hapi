# 🚀 HAPI - HTTP API Query Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/ermos/hapi.svg)](https://pkg.go.dev/github.com/ermos/hapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermos/hapi)](https://goreportcard.com/report/github.com/ermos/hapi)

A clean, elegant Go library for parsing HTTP API query parameters with support for filtering, sorting, and pagination. Transform complex query strings into structured data effortlessly.

## ✨ Features

- **Rich Filtering**: Support for 12 different operators (`eq`, `ne`, `gt`, `lt`, `ge`, `le`, `lk`, `nlk`, `in`, `nin`, `inlk`, `ninlk`)
- **Flexible Sorting**: Parse field and direction from query parameters
- **Pagination Support**: Built-in limit and offset handling
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
    url := "http://api.example.com/users?name[lk]=John%25&age[ge]=18&status[in]=active,pending&limit=25&offset=50&sort=created_at:desc"

    result, err := hapi.Parse(url)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d filters\n", len(result.Filters))
    fmt.Printf("Sort by: %s %s\n", result.Sort.Field, result.Sort.Direction)
    fmt.Printf("Pagination: limit=%d, offset=%d\n", result.Limit, result.Offset)
}
```

## 📖 Usage Examples

### Basic Filtering

```go
// URL: /users?name=John&age[gt]=25&status[in]=active,pending

result, _ := hapi.Parse(url)

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
    result, err := hapi.ParseFromRequest(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Use parsed result for database queries
    applyFiltersToQuery(result.Filters)
    applySortingToQuery(result.Sort)
    applyPaginationToQuery(result.Limit, result.Offset)
}
```

### Strict Mode

```go
// Use strict mode for validation
result, err := hapi.ParseStrict(url)
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
sort=field:direction
sort=name:asc
sort=created_at:desc
```

### Pagination
```
limit=25
offset=50
```

## 🏗️ API Reference

### Core Functions

```go
// Parse URL string (lenient mode)
func Parse(url string) (Result, error)

// Parse URL string (strict mode)
func ParseStrict(url string) (Result, error)

// Parse from HTTP request (lenient mode)
func ParseFromRequest(r *http.Request) (Result, error)

// Parse from HTTP request (strict mode)
func ParseFromRequestStrict(r *http.Request) (Result, error)
```

### Result Structure

```go
type Result struct {
    Filters Filters // Collection of filter conditions
    Sort    Sort    // Sorting configuration
    Limit   int     // Maximum results (0 = unlimited)
    Offset  int     // Skip results (1-based)
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
