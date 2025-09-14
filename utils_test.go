package hapi

import "testing"

func TestContainsString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		str      string
		expected bool
	}{
		{
			name:     "string exists in slice",
			slice:    []string{"apple", "banana", "orange"},
			str:      "banana",
			expected: true,
		},
		{
			name:     "string does not exist in slice",
			slice:    []string{"apple", "banana", "orange"},
			str:      "grape",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			str:      "anything",
			expected: false,
		},
		{
			name:     "nil slice",
			slice:    nil,
			str:      "anything",
			expected: false,
		},
		{
			name:     "empty string in non-empty slice",
			slice:    []string{"apple", "", "orange"},
			str:      "",
			expected: true,
		},
		{
			name:     "case sensitive comparison",
			slice:    []string{"Apple", "Banana", "Orange"},
			str:      "apple",
			expected: false,
		},
		{
			name:     "single element slice - match",
			slice:    []string{"only"},
			str:      "only",
			expected: true,
		},
		{
			name:     "single element slice - no match",
			slice:    []string{"only"},
			str:      "other",
			expected: false,
		},
		{
			name:     "duplicate values in slice",
			slice:    []string{"apple", "banana", "apple", "orange"},
			str:      "apple",
			expected: true,
		},
		{
			name:     "special characters",
			slice:    []string{"field[eq]", "field[ne]", "field[gt]"},
			str:      "field[eq]",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsString(tt.slice, tt.str)
			if result != tt.expected {
				t.Errorf("containsString(%v, %q) = %v, want %v",
					tt.slice, tt.str, result, tt.expected)
			}
		})
	}
}

func BenchmarkContainsString(b *testing.B) {
	slice := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta", "eta", "theta", "iota", "kappa"}

	b.Run("found at beginning", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			containsString(slice, "alpha")
		}
	})

	b.Run("found at end", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			containsString(slice, "kappa")
		}
	})

	b.Run("not found", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			containsString(slice, "notfound")
		}
	})
}