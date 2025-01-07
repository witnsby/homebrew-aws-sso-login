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
	"regexp"
	"strings"
)

// Constants
const (
	repo        = "witnsby/aws-sso-login"      // GitHub repository
	formulaPath = "./Formula/aws-sso-login.rb" // Formula file path
	baseURL     = "https://github.com/%s/archive/refs/tags/%s.tar.gz"
	apiURL      = "https://api.github.com/repos/%s/releases/latest"
)

// Struct to parse GitHub API response
type Release struct {
	TagName string `json:"tag_name"`
}

// Fetch the latest release tag using GitHub API
func fetchLatestReleaseTag(repo string) (string, error) {
	url := fmt.Sprintf(apiURL, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch release: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", errors.New("no tag name found")
	}

	return release.TagName, nil
}

// Download tarball file
func downloadTarball(repo, tag, outputPath string) (string, error) {
	url := fmt.Sprintf(baseURL, repo, tag)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download tarball: %s", resp.Status)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return url, nil
}

// Generate SHA256 checksum for a file
func generateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Update formula file: replace URL and SHA256
func updateFormula(filePath, url, sha256, version string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	text := string(content)

	// Replace URL
	re := regexp.MustCompile(`url ".*?"`)
	text = re.ReplaceAllString(text, fmt.Sprintf(`url "%s"`, url))

	// Replace SHA256
	re = regexp.MustCompile(`sha256 ".*?"`)
	text = re.ReplaceAllString(text, fmt.Sprintf(`sha256 "%s"`, sha256))

	// Replace version (if exists in formula)
	re = regexp.MustCompile(`version ".*?"`)
	text = re.ReplaceAllString(text, fmt.Sprintf(`version "%s"`, strings.TrimPrefix(version, "v")))

	err = os.WriteFile(filePath, []byte(text), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Fetching latest release tag...")
	tag, err := fetchLatestReleaseTag(repo)
	if err != nil {
		fmt.Println("Error fetching release tag:", err)
		return
	}
	fmt.Println("Latest tag:", tag)

	tempFilePath := "temp.tar.gz"

	fmt.Println("Downloading tarball...")
	tarballURL, err := downloadTarball(repo, tag, tempFilePath)
	if err != nil {
		fmt.Println("Error downloading tarball:", err)
		return
	}
	defer os.Remove(tempFilePath)

	fmt.Println("Generating SHA256 checksum...")
	sha256, err := generateSHA256(tempFilePath)
	if err != nil {
		fmt.Println("Error generating SHA256:", err)
		return
	}
	fmt.Println("SHA256:", sha256)

	fmt.Println("Updating formula...")
	err = updateFormula(formulaPath, tarballURL, sha256, tag)
	if err != nil {
		fmt.Println("Error updating formula:", err)
		return
	}

	fmt.Println("Formula updated successfully!")
}
