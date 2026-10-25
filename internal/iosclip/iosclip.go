package iosclip

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Build(templatePath, siteName, org, clipURL, host string) (string, error) {
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	id := newUUID()
	clipID := newUUID()
	if org == "" {
		org = "app.jh.webclip"
	}
	if siteName == "" {
		siteName = "JH"
	}
	short := host
	if u, e := url.Parse(clipURL); e == nil && u.Host != "" {
		short = u.Host
	}
	icon := loadIcon(templatePath)
	out := string(raw)
	repl := map[string]string{
		"{$logo}":           base64.StdEncoding.EncodeToString(icon),
		"{$systemName}":     xmlEscape(siteName),
		"{$orgName}":        xmlEscape(org),
		"{$id}":             id,
		"{$clipId}":         clipID,
		"{$ios_url}":        xmlEscape(clipURL),
		"{$systemShortUrl}": xmlEscape(short),
	}
	for k, v := range repl {
		out = strings.ReplaceAll(out, k, v)
	}
	return out, nil
}

func StampURL(raw, source string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		if source == "" {
			return raw
		}
		sep := "?"
		if strings.Contains(raw, "?") {
			sep = "&"
		}
		return raw + sep + "source=" + url.QueryEscape(source)
	}
	q := u.Query()
	if source != "" {
		q.Set("source", source)
	}
	q.Del("invite")
	u.RawQuery = q.Encode()
	return u.String()
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func loadIcon(templatePath string) []byte {
	p := filepath.Join(filepath.Dir(templatePath), "icon-180.png")
	if b, err := os.ReadFile(p); err == nil && len(b) > 0 {
		return b
	}
	return pinkPNG()
}

func pinkPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 180, 180))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 92, B: 147, A: 255}}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
