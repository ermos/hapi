package hapi

import (
	"reflect"
	"testing"
)

func TestNewOptions(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		opts := NewOptions()

		if opts.DefaultPerPage != 10 {
			t.Errorf("DefaultPerPage = %d, want 10", opts.DefaultPerPage)
		}
		if opts.MaxPerPage != 100 {
			t.Errorf("MaxPerPage = %d, want 100", opts.MaxPerPage)
		}
		if len(opts.AllowedSorts) != 0 {
			t.Errorf("AllowedSorts = %v, want empty slice", opts.AllowedSorts)
		}
		if len(opts.AllowedFilters) != 0 {
			t.Errorf("AllowedFilters = %v, want empty slice", opts.AllowedFilters)
		}
	})

	t.Run("with option functions", func(t *testing.T) {
		opts := NewOptions(
			WithDefaultPerPage(25),
			WithMaxPerPage(200),
			WithAllowedSorts([]string{"name", "created_at"}),
			WithAllowedFilters([]string{"status", "type"}),
		)

		if opts.DefaultPerPage != 25 {
			t.Errorf("DefaultPerPage = %d, want 25", opts.DefaultPerPage)
		}
		if opts.MaxPerPage != 200 {
			t.Errorf("MaxPerPage = %d, want 200", opts.MaxPerPage)
		}
		if !reflect.DeepEqual(opts.AllowedSorts, []string{"name", "created_at"}) {
			t.Errorf("AllowedSorts = %v, want [name created_at]", opts.AllowedSorts)
		}
		if !reflect.DeepEqual(opts.AllowedFilters, []string{"status", "type"}) {
			t.Errorf("AllowedFilters = %v, want [status type]", opts.AllowedFilters)
		}
	})

	t.Run("chaining option functions", func(t *testing.T) {
		opts := NewOptions(
			WithDefaultPerPage(15),
			WithDefaultPerPage(30), // Should override previous value
		)

		if opts.DefaultPerPage != 30 {
			t.Errorf("DefaultPerPage = %d, want 30 (last value should win)", opts.DefaultPerPage)
		}
	})
}

func TestOptionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		optFunc  OptionFunc
		check    func(*Options) bool
		expected string
	}{
		{
			name:    "WithDefaultPerPage",
			optFunc: WithDefaultPerPage(50),
			check:   func(o *Options) bool { return o.DefaultPerPage == 50 },
			expected: "DefaultPerPage should be 50",
		},
		{
			name:    "WithMaxPerPage",
			optFunc: WithMaxPerPage(500),
			check:   func(o *Options) bool { return o.MaxPerPage == 500 },
			expected: "MaxPerPage should be 500",
		},
		{
			name:    "WithAllowedSorts",
			optFunc: WithAllowedSorts([]string{"id", "name", "date"}),
			check: func(o *Options) bool {
				return reflect.DeepEqual(o.AllowedSorts, []string{"id", "name", "date"})
			},
			expected: "AllowedSorts should match",
		},
		{
			name:    "WithAllowedFilters",
			optFunc: WithAllowedFilters([]string{"active", "type", "category"}),
			check: func(o *Options) bool {
				return reflect.DeepEqual(o.AllowedFilters, []string{"active", "type", "category"})
			},
			expected: "AllowedFilters should match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{
				DefaultPerPage: 10,
				MaxPerPage:     100,
				AllowedSorts:   []string{},
				AllowedFilters: []string{},
			}

			tt.optFunc(opts)

			if !tt.check(opts) {
				t.Errorf("%s: check failed", tt.expected)
			}
		})
	}
}

func TestOptionsEdgeCases(t *testing.T) {
	t.Run("zero and negative values", func(t *testing.T) {
		opts := NewOptions(
			WithDefaultPerPage(0),
			WithMaxPerPage(-10),
		)

		// These should be accepted as-is, validation happens during parsing
		if opts.DefaultPerPage != 0 {
			t.Errorf("DefaultPerPage = %d, want 0", opts.DefaultPerPage)
		}
		if opts.MaxPerPage != -10 {
			t.Errorf("MaxPerPage = %d, want -10", opts.MaxPerPage)
		}
	})

	t.Run("nil slices", func(t *testing.T) {
		opts := NewOptions(
			WithAllowedSorts(nil),
			WithAllowedFilters(nil),
		)

		if opts.AllowedSorts != nil {
			t.Errorf("AllowedSorts = %v, want nil", opts.AllowedSorts)
		}
		if opts.AllowedFilters != nil {
			t.Errorf("AllowedFilters = %v, want nil", opts.AllowedFilters)
		}
	})

	t.Run("empty slices", func(t *testing.T) {
		empty := []string{}
		opts := NewOptions(
			WithAllowedSorts(empty),
			WithAllowedFilters(empty),
		)

		if !reflect.DeepEqual(opts.AllowedSorts, empty) {
			t.Errorf("AllowedSorts = %v, want empty slice", opts.AllowedSorts)
		}
		if !reflect.DeepEqual(opts.AllowedFilters, empty) {
			t.Errorf("AllowedFilters = %v, want empty slice", opts.AllowedFilters)
		}
	})
}