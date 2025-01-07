package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	repo             = "witnsby/aws-sso-login"
	formulaPath      = "./Formula/aws-sso-login.rb"
	latestReleaseAPI = "https://api.github.com/repos/%s/releases/latest"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func fetchLatestRelease(repo string) (*Release, error) {
	url := fmt.Sprintf(latestReleaseAPI, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	if release.TagName == "" {
		return nil, errors.New("no tag_name in release")
	}
	if len(release.Assets) == 0 {
		return nil, errors.New("no assets in release")
	}

	return &release, nil
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

func main() {
	release, err := fetchLatestRelease(repo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	wanted := []string{
		"aws-sso-login_darwin_amd64",
		"aws-sso-login_darwin_arm64",
		"aws-sso-login_linux_amd64",
		"aws-sso-login_linux_arm64",
	}

	binaries := make(map[string]string)

	for _, asset := range release.Assets {
		for _, w := range wanted {
			if asset.Name == w {
				key := strings.TrimPrefix(w, "aws-sso-login_")
				tmpFile := "./" + asset.Name
				if err := downloadFile(asset.BrowserDownloadURL, tmpFile); err != nil {
					fmt.Printf("Error downloading %s: %v\n", asset.Name, err)
					return
				}
				defer os.Remove(tmpFile)

				sum, err := generateSHA256(tmpFile)
				if err != nil {
					fmt.Printf("Error generating SHA256 for %s: %v\n", asset.Name, err)
					return
				}
				binaries[key] = asset.BrowserDownloadURL + "#" + sum
			}
		}
	}

	for _, w := range wanted {
		key := strings.TrimPrefix(w, "aws-sso-login_")
		if _, ok := binaries[key]; !ok {
			fmt.Printf("Warning: %s not found in release assets\n", w)
		}
	}

	if err := updateFormula(formulaPath, release.TagName, binaries); err != nil {
		fmt.Println("Error updating formula:", err)
		return
	}

	fmt.Println("Formula updated successfully!")
}
