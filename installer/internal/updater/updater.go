package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// ReleasesURL is the GitHub API endpoint for the latest release.
// Override in tests via httptest.
var ReleasesURL = "https://api.github.com/repos/collectiveai-team/agent-skills/releases/latest"

// Asset represents a single file attached to a GitHub release.
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Release represents the relevant fields of a GitHub release.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// UpdateResult describes whether an update is available.
type UpdateResult struct {
	Available  bool
	CurrentVer string
	LatestVer  string
	Asset      Asset
}

// CheckForUpdate queries GitHub for the latest release and compares it to
// currentVersion. Returns Available: false for "dev" builds or when no
// newer version exists.
func CheckForUpdate(currentVersion string) (*UpdateResult, error) {
	if currentVersion == "dev" {
		return &UpdateResult{Available: false, CurrentVer: currentVersion}, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(ReleasesURL)
	if err != nil {
		return nil, fmt.Errorf("checking for update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decoding release: %w", err)
	}

	if compareVersions(rel.TagName, currentVersion) <= 0 {
		return &UpdateResult{
			Available:  false,
			CurrentVer: currentVersion,
			LatestVer:  rel.TagName,
		}, nil
	}

	want := assetName()
	for _, a := range rel.Assets {
		if a.Name == want {
			return &UpdateResult{
				Available:  true,
				CurrentVer: currentVersion,
				LatestVer:  rel.TagName,
				Asset:      a,
			}, nil
		}
	}

	return nil, fmt.Errorf("no asset matching %q in release %s", want, rel.TagName)
}

// DownloadAndReplace downloads the asset and atomically replaces the current
// executable. Progress is printed to stderr with \r overwrites.
func DownloadAndReplace(asset Asset) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("resolving symlinks: %w", err)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("downloading update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	dir := filepath.Dir(exePath)
	tmp, err := os.CreateTemp(dir, "agent-setup-update-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	// Clean up temp file on error.
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	total := resp.ContentLength
	written := int64(0)
	buf := make([]byte, 32*1024)

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := tmp.Write(buf[:n]); wErr != nil {
				tmp.Close()
				err = fmt.Errorf("writing update: %w", wErr)
				return err
			}
			written += int64(n)
			if total > 0 {
				pct := float64(written) / float64(total) * 100
				fmt.Fprintf(os.Stderr, "\r  Downloading... %.0f%%", pct)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			tmp.Close()
			err = fmt.Errorf("reading update: %w", readErr)
			return err
		}
	}
	fmt.Fprintln(os.Stderr) // newline after progress

	if err = tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	if err = os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	// Atomic replace. On Windows, rename the old binary first.
	if runtime.GOOS == "windows" {
		oldPath := exePath + ".old"
		os.Remove(oldPath) // ignore error; may not exist
		if err = os.Rename(exePath, oldPath); err != nil {
			return fmt.Errorf("moving old binary: %w", err)
		}
	}

	if err = os.Rename(tmpPath, exePath); err != nil {
		return fmt.Errorf("replacing binary: %w", err)
	}

	return nil
}

// compareVersions compares two semver strings (with optional "v" prefix).
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func compareVersions(a, b string) int {
	pa := parseVersion(a)
	pb := parseVersion(b)

	for i := 0; i < 3; i++ {
		if pa[i] < pb[i] {
			return -1
		}
		if pa[i] > pb[i] {
			return 1
		}
	}
	return 0
}

// parseVersion extracts [major, minor, patch] from a version string.
func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		n, _ := strconv.Atoi(parts[i])
		result[i] = n
	}
	return result
}

// assetName returns the expected binary name for the current platform,
// e.g. "agent-setup-darwin-arm64" or "agent-setup-windows-amd64.exe".
func assetName() string {
	name := fmt.Sprintf("agent-setup-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}
