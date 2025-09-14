package hapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ParseFromRequest parses query parameters from an HTTP request.
// Invalid parameters are silently ignored.
func ParseFromRequest(r *http.Request, opts Options) (Result, error) {
	if r == nil || r.URL == nil {
		return Result{}, fmt.Errorf("request or URL is nil")
	}
	return parseFromURL(r.URL.String(), opts, false)
}

// ParseFromRequestStrict parses query parameters from an HTTP request.
// Returns an error for any invalid parameters.
func ParseFromRequestStrict(r *http.Request, opts Options) (Result, error) {
	if r == nil || r.URL == nil {
		return Result{}, fmt.Errorf("request or URL is nil")
	}
	return parseFromURL(r.URL.String(), opts, true)
}

// Parse parses query parameters from a URL string.
// Invalid parameters are silently ignored.
func Parse(url string, opts Options) (Result, error) {
	return parseFromURL(url, opts, false)
}

// ParseStrict parses query parameters from a URL string.
// Returns an error for any invalid parameters.
func ParseStrict(url string, opts Options) (Result, error) {
	return parseFromURL(url, opts, true)
}

func parseFromURL(u string, opts Options, strict bool) (result Result, err error) {
	result = Result{
		PerPage: opts.DefaultPerPage,
		Page:    1,
		Sorts:   make(Sorts, 0),
		Filters: make(Filters, 0),
	}

	res, err := url.Parse(u)
	if err != nil {
		return
	}

	q := res.RawQuery

	for _, filter := range strings.Split(q, "&") {
		var values Values

		parts := strings.Split(filter, "=")
		if parts[0] == "per_page" {
			if len(parts) != 2 {
				if strict {
					return Result{}, fmt.Errorf("invalid per_page filter format: %s", filter)
				}
				continue
			}

			result.PerPage = min(max(1, Value(parts[1]).Int()), opts.MaxPerPage)
			continue
		} else if parts[0] == "page" {
			if len(parts) != 2 {
				if strict {
					return Result{}, fmt.Errorf("invalid page filter format: %s", filter)
				}
				continue
			}

			result.Page = max(1, Value(parts[1]).Int())
			continue
		} else if parts[0] == "sort" {
			var sort Sort
			if len(parts) != 2 {
				if strict {
					return Result{}, fmt.Errorf("invalid sort filter format: %s", filter)
				}
				continue
			}

			sort, err = parseSortFromString(parts[1])
			if err != nil {
				if strict {
					return Result{}, err
				}
				continue
			}

			if len(opts.AllowedSorts) > 0 && !containsString(opts.AllowedSorts, sort.Field) {
				if strict {
					return Result{}, fmt.Errorf("sorting by field %q is not allowed", sort.Field)
				}

				continue
			}

			result.Sorts = append(result.Sorts, sort)

			continue
		}

		if len(opts.AllowedFilters) > 0 && !containsString(opts.AllowedFilters, parts[0]) {
			if strict {
				return Result{}, fmt.Errorf("filtering by field %q is not allowed", parts[0])
			}

			continue
		}

		if len(parts) != 2 {
			result.Filters = append(result.Filters, Filter{
				Field:    parts[0],
				Operator: FilterOperatorEqual,
				Values:   Values{""},
			})
			continue
		}

		field := parts[0]
		operator := FilterOperatorEqual
		value := parts[1]

		if strings.Contains(field, "[") {
			fieldParts := strings.Split(field, "[")
			field = fieldParts[0]
			operator = FilterOperator(fieldParts[1][:len(fieldParts[1])-1])
		}

		if err = operator.Valid(); err != nil {
			if strict {
				return Result{}, err
			}
			continue
		}

		if operator.IsList() {
			for _, v := range strings.Split(value, ",") {
				var unescaped string

				unescaped, err = url.PathUnescape(v)
				if err != nil {
					if strict {
						return Result{}, fmt.Errorf("failed to unescape value %q: %w", v, err)
					}
					continue
				}

				values = append(values, Value(unescaped))
			}
		} else {
			var unescaped string

			unescaped, err = url.PathUnescape(value)
			if err != nil {
				if strict {
					return Result{}, fmt.Errorf("failed to unescape value %q: %w", value, err)
				}
				continue
			}

			values = append(values, Value(unescaped))
		}

		result.Filters = append(result.Filters, Filter{
			Field:    field,
			Operator: operator,
			Values:   values,
		})
	}

	return
}
