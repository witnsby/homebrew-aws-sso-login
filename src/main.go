package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	repo        = "witnsby/aws-sso-login"
	formulaPath = "./Formula/aws-sso-login.rb"
	githubAPI   = "https://api.github.com"
	// releaseTagEnv names the release to package; empty means the latest
	// release. The Fetch Release workflow sets it from its "tag" input.
	releaseTagEnv = "RELEASE_TAG"
)

// wantedAssets are the release binaries the formula installs.
var wantedAssets = []string{
	"aws-sso-login_darwin_amd64",
	"aws-sso-login_darwin_arm64",
	"aws-sso-login_linux_amd64",
	"aws-sso-login_linux_arm64",
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// releaseURL returns the GitHub API URL of the release with the given tag,
// or of the latest release when tag is empty.
func releaseURL(apiBase, repo, tag string) string {
	base := strings.TrimRight(apiBase, "/") + "/repos/" + repo + "/releases/"
	if tag == "" {
		return base + "latest"
	}
	return base + "tags/" + url.PathEscape(tag)
}

// fetchRelease fetches the release with the given tag, or the latest release
// when tag is empty.
func fetchRelease(apiBase, repo, tag string) (*Release, error) {
	name := tag
	if name == "" {
		name = "latest"
	}
	resp, err := http.Get(releaseURL(apiBase, repo, tag))
	if err != nil {
		return nil, fmt.Errorf("fetch release %s of %s: %w", name, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("release %s not found in %s", name, repo)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch release %s of %s: %s", name, repo, resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode release %s of %s: %w", name, repo, err)
	}

	if release.TagName == "" {
		return nil, errors.New("no tag_name in release")
	}
	if tag != "" && release.TagName != tag {
		return nil, fmt.Errorf("asked for release %s, got %s", tag, release.TagName)
	}
	if len(release.Assets) == 0 {
		return nil, fmt.Errorf("release %s has no assets", release.TagName)
	}

	return &release, nil
}

// selectAssets returns the wanted binaries of the release by name and fails
// when any of them is missing, so the formula is never written with an
// empty URL.
func selectAssets(release *Release) (map[string]Asset, error) {
	byName := make(map[string]Asset, len(release.Assets))
	for _, a := range release.Assets {
		byName[a.Name] = a
	}
	selected := make(map[string]Asset, len(wantedAssets))
	var missing []string
	for _, w := range wantedAssets {
		a, ok := byName[w]
		if !ok {
			missing = append(missing, w)
			continue
		}
		selected[w] = a
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("release %s lacks %s", release.TagName, strings.Join(missing, ", "))
	}
	return selected, nil
}

func downloadFile(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func generateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func updateFormula(filePath, version string, binaries map[string]string) error {
	getVals := func(key string) (url, sha string) {
		val, ok := binaries[key]
		if !ok {
			return "", ""
		}
		parts := strings.SplitN(val, "#", 2)
		if len(parts) != 2 {
			return "", ""
		}
		return parts[0], parts[1]
	}

	macIntelURL, macIntelSHA := getVals("darwin_amd64")
	macArmURL, macArmSHA := getVals("darwin_arm64")
	linuxIntelURL, linuxIntelSHA := getVals("linux_amd64")
	linuxArmURL, linuxArmSHA := getVals("linux_arm64")

	formula := fmt.Sprintf(`class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "%s"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "%s"
      sha256 "%s"
    end
    on_arm do
      url "%s"
      sha256 "%s"
    end
  end

  on_linux do
    on_intel do
      url "%s"
      sha256 "%s"
    end
    on_arm do
      url "%s"
      sha256 "%s"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
`,
		strings.TrimPrefix(version, "v"),
		macIntelURL, macIntelSHA,
		macArmURL, macArmSHA,
		linuxIntelURL, linuxIntelSHA,
		linuxArmURL, linuxArmSHA,
	)

	return os.WriteFile(filePath, []byte(formula), 0644)
}

func run() error {
	release, err := fetchRelease(githubAPI, repo, os.Getenv(releaseTagEnv))
	if err != nil {
		return err
	}
	assets, err := selectAssets(release)
	if err != nil {
		return err
	}

	binaries := make(map[string]string, len(assets))
	for name, asset := range assets {
		key := strings.TrimPrefix(name, "aws-sso-login_")
		tmpFile := "./" + asset.Name
		if err := downloadFile(asset.BrowserDownloadURL, tmpFile); err != nil {
			return fmt.Errorf("download %s: %w", asset.Name, err)
		}
		sum, err := generateSHA256(tmpFile)
		os.Remove(tmpFile)
		if err != nil {
			return fmt.Errorf("SHA256 of %s: %w", asset.Name, err)
		}
		binaries[key] = asset.BrowserDownloadURL + "#" + sum
	}

	if err := updateFormula(formulaPath, release.TagName, binaries); err != nil {
		return fmt.Errorf("update formula: %w", err)
	}
	fmt.Printf("Formula updated to %s.\n", release.TagName)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
