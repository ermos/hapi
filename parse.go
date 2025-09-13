package hapi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func ParseFromRequest(r *http.Request) (result Result, err error) {
	return parseFromURL(r.URL.String(), false)
}

func ParseFromRequestStrict(r *http.Request) (result Result, err error) {
	return parseFromURL(r.URL.String(), true)
}

func Parse(url string) (result Result, err error) {
	return parseFromURL(url, false)
}

func ParseStrict(url string) (result Result, err error) {
	return parseFromURL(url, true)
}

func parseFromURL(u string, strict bool) (result Result, err error) {
	res, err := url.Parse(u)
	if err != nil {
		return
	}

	q := res.RawQuery

	for _, filter := range strings.Split(q, "&") {
		var values Values

		parts := strings.Split(filter, "=")
		if parts[0] == "limit" {
			if len(parts) != 2 {
				if !strict {
					continue
				}

				err = fmt.Errorf("invalid limit filter format: %s", filter)
				return
			}

			result.Limit = max(0, Value(parts[1]).Int())
			continue
		} else if parts[0] == "offset" {
			if len(parts) != 2 {
				if !strict {
					continue
				}

				err = fmt.Errorf("invalid offset filter format: %s", filter)
				return
			}

			result.Offset = max(1, Value(parts[1]).Int())
			continue
		} else if parts[0] == "sort" {
			if len(parts) != 2 {
				if !strict {
					continue
				}

				err = fmt.Errorf("invalid sort filter format: %s", filter)
				return
			}

			result.Sort, err = parseSortFromString(parts[1])
			if err != nil {
				if !strict {
					continue
				}

				return
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
			if !strict {
				continue
			}
			return
		}

		if operator.IsList() {
			for _, v := range strings.Split(value, ",") {
				var unescaped string

				unescaped, err = url.PathUnescape(v)
				if err != nil {
					if !strict {
						continue
					}
					return
				}

				values = append(values, Value(unescaped))
			}
		} else {
			var unescaped string

			unescaped, err = url.PathUnescape(value)
			if err != nil {
				if !strict {
					continue
				}
				return
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
