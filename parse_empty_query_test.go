package hapi

import "testing"

// A URL with no query string (or a stray "&") must not produce a phantom filter.
func TestParse_EmptyQueryNoPhantomFilter(t *testing.T) {
	for _, u := range []string{
		"http://x/users",
		"http://x/users?",
		"http://x/users?&name=John&",
	} {
		r, err := Parse(u, Options{DefaultPerPage: 10, MaxPerPage: 100})
		if err != nil {
			t.Fatalf("Parse(%q): %v", u, err)
		}
		for _, f := range r.Filters {
			if f.Field == "" {
				t.Errorf("Parse(%q) produced phantom empty-field filter: %#v", u, r.Filters)
			}
		}
	}
}
