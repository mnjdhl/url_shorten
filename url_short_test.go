package main

import (
	"testing"
)

func TestGetShortURL(t *testing.T) {
	expurl := "ti.ny/10c2398f"
	surl := GetShortURL("www.abc.com/longlongurlkjjsss")
	if surl != expurl {
		t.Errorf("Failed to get short url")
	} else {
		surl2 := GetShortURL("www.abc.com/longlongurlkjjsss")
		if surl != surl2 {
			t.Errorf("Failed to match short url with previous one")
		} else {
			t.Logf("Success:Got short url = %s", surl)
		}
	}
}
