package hapi

import "testing"

// A "=" inside the value (e.g. base64 padding, tokens) must be preserved.
func TestParse_EqualsInValuePreserved(t *testing.T) {
	r, err := Parse("http://x/users?token=YWJjZA==", Options{DefaultPerPage: 10, MaxPerPage: 100})
	if err != nil {
		t.Fatal(err)
	}
	got := r.Filters.GetFirstFromField("token").Values.First().String()
	if got != "YWJjZA==" {
		t.Errorf("token value = %q, want %q", got, "YWJjZA==")
	}
}
