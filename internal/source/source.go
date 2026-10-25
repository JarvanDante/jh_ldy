package source

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	CookieName = "ldy_source"
	Prefix     = "channel://"
	CookieAge  = 5 * 24 * 3600
)

var codeRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{1,31}$`)

var reserved = map[string]struct{}{
	"official": {},
	"404":      {},
	"invite":   {},
	"source":   {},
	"channel":  {},
	"share":    {},
	"agent":    {},
}

func Parse(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if strings.Contains(s, "share://") || strings.Contains(s, "agent://") {
		return ""
	}
	if i := strings.Index(s, Prefix); i >= 0 {
		s = strings.TrimSpace(s[i+len(Prefix):])
	}
	if i := strings.IndexAny(s, "?&#/ \t\n"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if !codeRe.MatchString(s) {
		return ""
	}
	if _, ok := reserved[strings.ToLower(s)]; ok {
		return ""
	}
	return s
}

func Clipboard(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	return Prefix + code
}

func FromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	if s := Parse(r.URL.Query().Get("source")); s != "" {
		return s
	}
	if c, err := r.Cookie(CookieName); err == nil {
		if s := Parse(c.Value); s != "" {
			return s
		}
	}
	return ""
}

func SetCookie(w http.ResponseWriter, code string) {
	code = Parse(code)
	if code == "" {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    code,
		Path:     "/",
		MaxAge:   CookieAge,
		SameSite: http.SameSiteLaxMode,
	})
}
