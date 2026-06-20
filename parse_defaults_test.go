package hapi

import "testing"

// A zero-value Options must fall back to package defaults instead of yielding
// PerPage 0 (min(x, MaxPerPage=0) == 0).
func TestParse_ZeroOptionsFallBackToDefaults(t *testing.T) {
	r, err := Parse("http://x/users?per_page=50", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if r.PerPage != 50 {
		t.Errorf("PerPage = %d, want 50 (capped by default max %d)", r.PerPage, defaultMaxPerPage)
	}

	r, _ = Parse("http://x/users?per_page=9999", Options{})
	if r.PerPage != defaultMaxPerPage {
		t.Errorf("PerPage = %d, want default max %d", r.PerPage, defaultMaxPerPage)
	}
}
