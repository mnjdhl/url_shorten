package main

import (
	"testing"
	"time"
)

var urlTab = []URLData{URLData{"www.abc.com/longlongurlkjjsss", "ti.ny/10c2398f", time.Time{}}, URLData{"www.xyz.com/longlongurl2nmmmsssakkssstyy", "ti.ny/f17fc388", time.Time{}}, URLData{"www.efg.com/longlongurl3nhgjkkkssslkkhhaa", "ti.ny/0de0929b", time.Time{}}}

func TestGetShortURL(t *testing.T) {
	for _, lu := range urlTab {
		surl := GetShortURL(lu.LongURL)
		if surl != lu.ShortURL {
			t.Errorf("Failed to get short url")
		} else {
			surl2 := GetShortURL(lu.LongURL)
			if surl != surl2 {
				t.Errorf("Failed to match short url with previous one")
			} else {
				t.Logf("Success:Got short url = %s", surl)
			}
		}
	}
}

func TestInMemURLs(t *testing.T) {
	for _, lu := range urlTab {
		if URLTable[lu.LongURL].ShortURL != lu.ShortURL {
			t.Errorf("Failed to match short url with in-memory URL Table")
		}
	}
}

func TestInFileURLs(t *testing.T) {
	inFileURLs := LoadURLTable()
	for _, lu := range urlTab {
		if inFileURLs[lu.LongURL].ShortURL != lu.ShortURL {
			t.Errorf("Failed to match short url with in-file URL Table")
		}
	}
}
