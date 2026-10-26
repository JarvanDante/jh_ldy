package httpserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/JarvanDante/jh_ldy/internal/conf"
	"github.com/JarvanDante/jh_ldy/internal/iosclip"
	"github.com/JarvanDante/jh_ldy/internal/source"
)

type Server struct {
	cfg     *conf.Config
	webDir  string
	skinDir string
	tplPath string
	apkPath string
}

func New(cfg *conf.Config, root string) *Server {
	if cfg != nil {
		cfg.Normalize()
	}
	return &Server{
		cfg:     cfg,
		webDir:  filepath.Join(root, "web"),
		skinDir: filepath.Join(root, "web", "skins"),
		tplPath: filepath.Join(root, "templates", "ios.mobileconfig"),
		apkPath: filepath.Join(root, "downloads", "app-release.apk"),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/boot", s.boot)
	mux.HandleFunc("/download/android.apk", s.androidApk)
	mux.HandleFunc("/download/ios.mobileconfig", s.iosConfig)
	mux.HandleFunc("/s/", s.pageWithPath)
	mux.Handle("/skins/", http.StripPrefix("/skins/", http.FileServer(http.Dir(s.skinDir))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(s.webDir, "js")))))
	mux.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(s.webDir, "favicon.svg"))
	})
	mux.HandleFunc("/", s.pageRoot)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) boot(w http.ResponseWriter, r *http.Request) {
	code := source.FromRequest(r)
	if code != "" {
		source.SetCookie(w, code)
	}
	page := s.pageFor(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(page)
}

func (s *Server) pageFor(code string) conf.Page {
	page := s.cfg.Resolve(code, s.skinDir)
	page.Clipboard = source.Clipboard(code)
	page.H5URL = iosclip.StampURL(page.H5URL, code)
	if strings.TrimSpace(page.AndroidURL) == "" {
		page.AndroidURL = "/download/android.apk?v=3"
	}
	if strings.TrimSpace(page.IOSURL) == "" {
		iosPath := "/download/ios.mobileconfig"
		if code != "" {
			iosPath += "?source=" + url.QueryEscape(code)
		}
		page.IOSURL = iosPath
	}
	return page
}

func (s *Server) androidApk(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(s.apkPath); err != nil {
		http.Error(w, "android apk missing", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.android.package-archive")
	w.Header().Set("Content-Disposition", `attachment; filename="jh.apk"`)
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, s.apkPath)
}

func (s *Server) iosConfig(w http.ResponseWriter, r *http.Request) {
	code := source.FromRequest(r)
	if code != "" {
		source.SetCookie(w, code)
	}
	clip := strings.TrimSpace(s.cfg.IOSClipURL)
	if ch, ok := s.cfg.FindChannel(code); ok && strings.TrimSpace(ch.IOSClipURL) != "" {
		clip = strings.TrimSpace(ch.IOSClipURL)
	}
	if clip == "" {
		clip = strings.TrimSpace(s.cfg.H5URL)
	}
	if clip == "" {
		clip = absOrigin(r) + "/"
		if code != "" {
			clip += "?source=" + url.QueryEscape(code)
		}
	} else {
		clip = iosclip.StampURL(clip, code)
	}
	body, err := iosclip.Build(s.tplPath, s.cfg.SiteName, s.cfg.IOSOrg, clip)
	if err != nil {
		http.Error(w, "ios template missing", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	w.Header().Set("Content-Disposition", `attachment; filename="jh-ios.mobileconfig"`)
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(body))
}

func (s *Server) pageWithPath(w http.ResponseWriter, r *http.Request) {
	code := source.Parse(strings.Trim(strings.TrimPrefix(r.URL.Path, "/s/"), "/"))
	if code != "" {
		source.SetCookie(w, code)
	}
	s.serveIndex(w, r)
}

func (s *Server) pageRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if code := source.FromRequest(r); code != "" {
		source.SetCookie(w, code)
	}
	s.serveIndex(w, r)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	index := filepath.Join(s.webDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		http.Error(w, "index missing", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, index)
}

func absOrigin(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
