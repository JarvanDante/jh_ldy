package conf

import "testing"

func TestResolveChannelAPKAndSkin(t *testing.T) {
	cfg := &Config{
		SiteName:    "漫隐",
		DefaultSkin: "default",
		AndroidURL:  "https://cdn.example.com/official.apk",
		Channels: map[string]Channel{
			"MYAAA": {AndroidURL: "https://cdn.example.com/myaaa.apk", Skin: "light"},
		},
	}
	cfg.Normalize()
	off := cfg.Resolve("", "")
	if off.AndroidURL != cfg.AndroidURL || off.Skin != "default" {
		t.Fatalf("default got apk=%s skin=%s", off.AndroidURL, off.Skin)
	}
	ch := cfg.Resolve("MYAAA", "")
	if ch.AndroidURL != "https://cdn.example.com/myaaa.apk" {
		t.Fatalf("channel apk %s", ch.AndroidURL)
	}
	if ch.Skin != "light" {
		t.Fatalf("channel skin %s", ch.Skin)
	}
}
