package hapi

import "strconv"

// Values represents a collection of Value items.
type Values []Value

// First returns the first value in the collection, or empty string if the collection is empty.
func (v Values) First() Value {
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// ToString converts all values to a string slice.
func (v Values) ToString() []string {
	if len(v) == 0 {
		return nil
	}

	// Pre-allocate the exact size needed
	res := make([]string, len(v))
	for i, value := range v {
		res[i] = value.String()
	}
	return res
}

// Value represents a single string value that can be converted to various types.
type Value string

// String returns the string representation of the value.
func (v Value) String() string {
	return string(v)
}

// Int converts the value to an integer, returning 0 if conversion fails.
func (v Value) Int() int {
	if v == "" {
		return 0
	}
	res, err := strconv.Atoi(string(v))
	if err != nil {
		return 0
	}
	return res
}

// Int64 converts the value to an int64, returning 0 if conversion fails.
func (v Value) Int64() int64 {
	if v == "" {
		return 0
	}
	res, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return 0
	}
	return res
}

// Float64 converts the value to a float64, returning 0 if conversion fails.
func (v Value) Float64() float64 {
	if v == "" {
		return 0
	}
	res, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0
	}
	return res
}

// Bool converts the value to a boolean, returning false if conversion fails.
// Recognizes "true", "false", "1", "0", "t", "f", "T", "F", "TRUE", "FALSE".
func (v Value) Bool() bool {
	if v == "" {
		return false
	}
	res, err := strconv.ParseBool(string(v))
	if err != nil {
		return false
	}
	return res
}
