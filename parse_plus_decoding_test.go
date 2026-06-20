package hapi

import "testing"

// Query values use "+" for spaces; it must be decoded like %20.
func TestParse_PlusDecodedAsSpace(t *testing.T) {
	cases := map[string]string{
		"http://x/users?name=John+Doe":   "John Doe",
		"http://x/users?name=John%20Doe": "John Doe",
		"http://x/users?q=a+b%2Bc":       "a b+c",
	}
	for u, want := range cases {
		r, err := Parse(u, Options{DefaultPerPage: 10, MaxPerPage: 100})
		if err != nil {
			t.Fatalf("Parse(%q): %v", u, err)
		}
		if len(r.Filters) == 0 {
			t.Fatalf("Parse(%q): no filters", u)
		}
		got := r.Filters[0].Values.First().String()
		if got != want {
			t.Errorf("Parse(%q) value = %q, want %q", u, got, want)
		}
	}
}
