package iosclip

import "testing"

func TestStampURLKeepsSourceDropsInvite(t *testing.T) {
	got := StampURL("https://h5.example.com/?invite=abc", "MYAAA")
	if got != "https://h5.example.com/?source=MYAAA" {
		t.Fatalf("got %s", got)
	}
}
