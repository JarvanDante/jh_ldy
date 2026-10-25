package conf

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Channel struct {
	AndroidURL string `json:"android_url"`
	IOSURL     string `json:"ios_url"`
	IOSClipURL string `json:"ios_clip_url"`
	Skin       string `json:"skin"`
}

type Config struct {
	Listen      string             `json:"listen"`
	SiteName    string             `json:"site_name"`
	Slogan      string             `json:"slogan"`
	LogoText    string             `json:"logo_text"`
	H5URL       string             `json:"h5_url"`
	DefaultSkin string             `json:"default_skin"`
	AndroidURL  string             `json:"android_url"`
	IOSURL      string             `json:"ios_url"`
	IOSOrg      string             `json:"ios_org"`
	IOSClipURL  string             `json:"ios_clip_url"`
	Channels    map[string]Channel `json:"channels"`
}

type Page struct {
	SiteName   string `json:"site_name"`
	Slogan     string `json:"slogan"`
	LogoText   string `json:"logo_text"`
	Source     string `json:"source"`
	Skin       string `json:"skin"`
	SkinCSS    string `json:"skin_css"`
	AndroidURL string `json:"android_url"`
	IOSURL     string `json:"ios_url"`
	H5URL      string `json:"h5_url"`
	Clipboard  string `json:"clipboard_value"`
	ThemeColor string `json:"theme_color"`
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err = json.Unmarshal(raw, cfg); err != nil {
		return nil, err
	}
	cfg.Normalize()
	return cfg, nil
}

func (c *Config) Normalize() {
	if c.Listen == "" {
		c.Listen = ":5175"
	}
	if c.SiteName == "" {
		c.SiteName = "JH"
	}
	if c.Slogan == "" {
		c.Slogan = "打开 App，马上开始"
	}
	if c.LogoText == "" {
		c.LogoText = "JH"
	}
	if c.DefaultSkin == "" {
		c.DefaultSkin = "default"
	}
	if c.IOSOrg == "" {
		c.IOSOrg = "app.jh.webclip"
	}
	if c.Channels == nil {
		c.Channels = map[string]Channel{}
	}
}

func (c *Config) FindChannel(source string) (Channel, bool) {
	if source == "" {
		return Channel{}, false
	}
	if ch, ok := c.Channels[source]; ok {
		return ch, true
	}
	want := strings.ToLower(source)
	for code, ch := range c.Channels {
		if strings.ToLower(code) == want {
			return ch, true
		}
	}
	return Channel{}, false
}

func (c *Config) Resolve(source, skinRoot string) Page {
	p := Page{
		SiteName:   c.SiteName,
		Slogan:     c.Slogan,
		LogoText:   c.LogoText,
		Source:     source,
		Skin:       c.DefaultSkin,
		AndroidURL: strings.TrimSpace(c.AndroidURL),
		IOSURL:     strings.TrimSpace(c.IOSURL),
		H5URL:      strings.TrimSpace(c.H5URL),
	}
	if ch, ok := c.FindChannel(source); ok {
		if v := strings.TrimSpace(ch.AndroidURL); v != "" {
			p.AndroidURL = v
		}
		if v := strings.TrimSpace(ch.IOSURL); v != "" {
			p.IOSURL = v
		}
		if v := strings.TrimSpace(ch.Skin); v != "" {
			p.Skin = v
		}
		if v := strings.TrimSpace(ch.IOSClipURL); v != "" && p.H5URL == "" {
			p.H5URL = v
		}
	}
	p.Skin = safeSkin(p.Skin, c.DefaultSkin, skinRoot)
	p.SkinCSS = "/skins/" + p.Skin + "/app.css?v=20261026"
	if p.Skin == "light" {
		p.ThemeColor = "#f6f2fa"
	} else {
		p.ThemeColor = "#0b0b0d"
	}
	return p
}

func safeSkin(name, fallback, root string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || strings.ContainsAny(name, `/\.`) {
		name = fallback
	}
	if root == "" {
		return name
	}
	if _, err := os.Stat(filepath.Join(root, name, "app.css")); err != nil {
		return fallback
	}
	return name
}
