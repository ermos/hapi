package hapi

import (
	"reflect"
	"testing"
)

func TestValuesFirst(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		want   Value
	}{
		{"Empty slice", Values{}, ""},
		{"Single value", Values{"test"}, "test"},
		{"Multiple values", Values{"first", "second", "third"}, "first"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.values.First(); got != tt.want {
				t.Errorf("Values.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValuesToString(t *testing.T) {
	tests := []struct {
		name   string
		values Values
		want   []string
	}{
		{"Empty slice", Values{}, nil},
		{"Single value", Values{"test"}, []string{"test"}},
		{"Multiple values", Values{"first", "second", "third"}, []string{"first", "second", "third"}},
		{"Numeric values", Values{"123", "456"}, []string{"123", "456"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.values.ToString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{"Empty value", Value(""), ""},
		{"String value", Value("test"), "test"},
		{"Numeric value", Value("123"), "123"},
		{"Special characters", Value("test@example.com"), "test@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Value.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueInt(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  int
	}{
		{"Valid positive integer", Value("123"), 123},
		{"Valid negative integer", Value("-456"), -456},
		{"Zero", Value("0"), 0},
		{"Invalid string", Value("abc"), 0},
		{"Empty string", Value(""), 0},
		{"Float string", Value("123.45"), 0},
		{"Large number", Value("999999999"), 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Int(); got != tt.want {
				t.Errorf("Value.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueInt64(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  int64
	}{
		{"Valid positive integer", Value("123"), 123},
		{"Valid negative integer", Value("-456"), -456},
		{"Zero", Value("0"), 0},
		{"Invalid string", Value("abc"), 0},
		{"Empty string", Value(""), 0},
		{"Float string", Value("123.45"), 0},
		{"Large number", Value("9223372036854775807"), 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Int64(); got != tt.want {
				t.Errorf("Value.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  float64
	}{
		{"Valid integer", Value("123"), 123.0},
		{"Valid float", Value("123.45"), 123.45},
		{"Valid negative float", Value("-456.78"), -456.78},
		{"Zero", Value("0"), 0.0},
		{"Zero float", Value("0.0"), 0.0},
		{"Invalid string", Value("abc"), 0.0},
		{"Empty string", Value(""), 0.0},
		{"Scientific notation", Value("1.23e2"), 123.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Float64(); got != tt.want {
				t.Errorf("Value.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueBool(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  bool
	}{
		{"True string", Value("true"), true},
		{"False string", Value("false"), false},
		{"True uppercase", Value("True"), true},
		{"False uppercase", Value("False"), false},
		{"True number", Value("1"), true},
		{"False number", Value("0"), false},
		{"Invalid string", Value("abc"), false},
		{"Empty string", Value(""), false},
		{"Yes", Value("yes"), false},
		{"No", Value("no"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.Bool(); got != tt.want {
				t.Errorf("Value.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
