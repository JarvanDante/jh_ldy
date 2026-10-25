package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/JarvanDante/jh_ldy/internal/conf"
)

func TestBootUsesChannelAPK(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "web", "skins", "default"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "web", "skins", "default", "app.css"), []byte("x"), 0o644)
	cfg := &conf.Config{
		AndroidURL:  "https://cdn.example.com/official.apk",
		DefaultSkin: "default",
		Channels: map[string]conf.Channel{
			"MYAAA": {AndroidURL: "https://cdn.example.com/myaaa.apk"},
		},
	}
	h := New(cfg, root).Handler()
	req := httptest.NewRequest(http.MethodGet, "/api/boot?source=MYAAA", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	var page conf.Page
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if page.AndroidURL != "https://cdn.example.com/myaaa.apk" {
		t.Fatalf("apk %s", page.AndroidURL)
	}
	if page.IOSURL != "/download/ios.mobileconfig?source=MYAAA" {
		t.Fatalf("ios %s", page.IOSURL)
	}
	if page.Clipboard != "channel://MYAAA" {
		t.Fatalf("clip %s", page.Clipboard)
	}
}
