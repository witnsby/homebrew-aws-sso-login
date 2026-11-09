package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReleaseURL(t *testing.T) {
	tests := []struct {
		tag  string
		want string
	}{
		{"", "https://api.github.com/repos/witnsby/aws-sso-login/releases/latest"},
		{"v0.4.0", "https://api.github.com/repos/witnsby/aws-sso-login/releases/tags/v0.4.0"},
		{"v1.0.0/../x", "https://api.github.com/repos/witnsby/aws-sso-login/releases/tags/v1.0.0%2F..%2Fx"},
	}
	for _, tt := range tests {
		if got := releaseURL("https://api.github.com/", repo, tt.tag); got != tt.want {
			t.Errorf("releaseURL(%q) = %q, want %q", tt.tag, got, tt.want)
		}
	}
}

func fullRelease(tag string) Release {
	r := Release{TagName: tag}
	for _, name := range wantedAssets {
		r.Assets = append(r.Assets, Asset{Name: name, BrowserDownloadURL: "https://example.com/" + tag + "/" + name})
	}
	r.Assets = append(r.Assets, Asset{Name: "aws-sso-login_checksums.txt", BrowserDownloadURL: "https://example.com/sums"})
	return r
}

// releaseServer serves releases by tag and "latest", and 404 otherwise.
func releaseServer(t *testing.T, releases map[string]Release) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "/repos/" + repo + "/releases/"
		key := strings.TrimPrefix(r.URL.Path, prefix)
		key = strings.TrimPrefix(key, "tags/")
		rel, ok := releases[key]
		if !strings.HasPrefix(r.URL.Path, prefix) || !ok {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(rel); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFetchReleaseByTag(t *testing.T) {
	srv := releaseServer(t, map[string]Release{
		"latest": fullRelease("v0.5.0"),
		"v0.4.0": fullRelease("v0.4.0"),
	})
	rel, err := fetchRelease(srv.URL, repo, "v0.4.0")
	if err != nil {
		t.Fatal(err)
	}
	if rel.TagName != "v0.4.0" {
		t.Errorf("tag = %s, want v0.4.0", rel.TagName)
	}
}

func TestFetchLatestRelease(t *testing.T) {
	srv := releaseServer(t, map[string]Release{"latest": fullRelease("v0.5.0")})
	rel, err := fetchRelease(srv.URL, repo, "")
	if err != nil {
		t.Fatal(err)
	}
	if rel.TagName != "v0.5.0" {
		t.Errorf("tag = %s, want v0.5.0", rel.TagName)
	}
}

func TestFetchUnknownTag(t *testing.T) {
	srv := releaseServer(t, map[string]Release{"latest": fullRelease("v0.5.0")})
	_, err := fetchRelease(srv.URL, repo, "v9.9.9")
	if err == nil || !strings.Contains(err.Error(), "release v9.9.9 not found") {
		t.Fatalf("err = %v, want release v9.9.9 not found", err)
	}
}

func TestFetchReleaseWithoutAssets(t *testing.T) {
	srv := releaseServer(t, map[string]Release{"v0.4.0": {TagName: "v0.4.0"}})
	_, err := fetchRelease(srv.URL, repo, "v0.4.0")
	if err == nil || !strings.Contains(err.Error(), "has no assets") {
		t.Fatalf("err = %v, want has no assets", err)
	}
}

func TestSelectAssets(t *testing.T) {
	rel := fullRelease("v0.4.0")
	assets, err := selectAssets(&rel)
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != len(wantedAssets) {
		t.Errorf("got %d assets, want %d", len(assets), len(wantedAssets))
	}
}

func TestSelectAssetsMissingBinary(t *testing.T) {
	rel := fullRelease("v0.4.0")
	rel.Assets = rel.Assets[1:] // drop darwin_amd64
	_, err := selectAssets(&rel)
	if err == nil || !strings.Contains(err.Error(), "lacks aws-sso-login_darwin_amd64") {
		t.Fatalf("err = %v, want missing darwin_amd64", err)
	}
}
