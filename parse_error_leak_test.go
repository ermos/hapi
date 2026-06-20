package hapi

import "testing"

// In lenient mode, an invalid trailing parameter must be ignored, not returned
// as an error (regression for the leaked named return value).
func TestParse_LenientDoesNotLeakError(t *testing.T) {
	for _, u := range []string{
		"http://x/users?name[xx]=foo",
		"http://x/users?ok=1&name[xx]=foo",
	} {
		if _, err := Parse(u, Options{DefaultPerPage: 10, MaxPerPage: 100}); err != nil {
			t.Errorf("Parse(%q) should not error in lenient mode, got: %v", u, err)
		}
	}
}
