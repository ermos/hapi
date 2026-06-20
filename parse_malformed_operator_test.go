package hapi

import "testing"

// Malformed operator brackets must not panic (regression for "[:-1]" slice).
func TestParse_MalformedOperatorNoPanic(t *testing.T) {
	for _, u := range []string{
		"http://x/users?name[=foo",
		"http://x/users?name[",
		"http://x/users?name[gt=foo",
	} {
		if _, err := Parse(u, Options{DefaultPerPage: 10, MaxPerPage: 100}); err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", u, err)
		}
	}

	if _, err := ParseStrict("http://x/users?name[=foo", Options{DefaultPerPage: 10, MaxPerPage: 100}); err == nil {
		t.Error("ParseStrict: expected error for malformed operator, got nil")
	}
}
