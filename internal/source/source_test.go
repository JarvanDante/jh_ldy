package source

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"  MYAAA  ", "MYAAA"},
		{"channel://MYAAA", "MYAAA"},
		{"share://abc", ""},
		{"official", ""},
		{"a", ""},
	}
	for _, c := range cases {
		if got := Parse(c.in); got != c.want {
			t.Fatalf("Parse(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestFromRequestPrefersQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/boot?source=MYAAA", nil)
	r.AddCookie(&http.Cookie{Name: CookieName, Value: "OTHER"})
	if got := FromRequest(r); got != "MYAAA" {
		t.Fatalf("got %q", got)
	}
}
