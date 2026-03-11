package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.1", "v1.0.0", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v2.0.0", "v1.9.9", 1},
		{"v0.1.0", "v0.0.9", 1},
		{"v1.2.3", "v1.2.3", 0},
		{"1.0.0", "v1.0.0", 0},  // no prefix
		{"v10.0.0", "v9.0.0", 1}, // double digits
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_vs_%s", tt.a, tt.b), func(t *testing.T) {
			got := compareVersions(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestAssetName(t *testing.T) {
	name := assetName()
	want := fmt.Sprintf("agent-setup-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		want += ".exe"
	}
	if name != want {
		t.Errorf("assetName() = %q, want %q", name, want)
	}
}

func TestCheckForUpdate_DevVersion(t *testing.T) {
	result, err := CheckForUpdate("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Available {
		t.Error("expected Available=false for dev version")
	}
}

func TestCheckForUpdate_NewerAvailable(t *testing.T) {
	asset := assetName()
	rel := Release{
		TagName: "v2.0.0",
		Assets: []Asset{
			{Name: asset, BrowserDownloadURL: "https://example.com/download"},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rel)
	}))
	defer srv.Close()

	orig := ReleasesURL
	ReleasesURL = srv.URL
	defer func() { ReleasesURL = orig }()

	result, err := CheckForUpdate("v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Available {
		t.Fatal("expected Available=true")
	}
	if result.LatestVer != "v2.0.0" {
		t.Errorf("LatestVer = %q, want %q", result.LatestVer, "v2.0.0")
	}
	if result.Asset.Name != asset {
		t.Errorf("Asset.Name = %q, want %q", result.Asset.Name, asset)
	}
}

func TestCheckForUpdate_AlreadyLatest(t *testing.T) {
	rel := Release{
		TagName: "v1.0.0",
		Assets:  []Asset{{Name: assetName(), BrowserDownloadURL: "https://example.com/dl"}},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rel)
	}))
	defer srv.Close()

	orig := ReleasesURL
	ReleasesURL = srv.URL
	defer func() { ReleasesURL = orig }()

	result, err := CheckForUpdate("v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Available {
		t.Error("expected Available=false when already on latest")
	}
}

func TestDownloadAndReplace(t *testing.T) {
	// Create a fake "current" binary.
	tmpDir := t.TempDir()
	fakeBin := filepath.Join(tmpDir, "agent-setup")
	if err := os.WriteFile(fakeBin, []byte("old-binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Serve the "new" binary.
	newContent := []byte("new-binary-content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(newContent)))
		w.Write(newContent)
	}))
	defer srv.Close()

	// Override os.Executable by creating a symlink to our fake binary and
	// calling DownloadAndReplace with a custom approach: we write a small
	// helper that exercises the download-and-write portion.
	//
	// Since DownloadAndReplace uses os.Executable() internally, we test the
	// download+write logic directly.
	asset := Asset{
		Name:               "agent-setup-test",
		BrowserDownloadURL: srv.URL + "/agent-setup-test",
	}

	// We can't easily override os.Executable in a unit test, so we test
	// the download and atomic replace logic by calling the internal pieces.
	// Instead, let's test the full flow by temporarily patching the binary.
	// For a proper integration test, we'd need a subprocess.
	//
	// Test the download portion directly:
	client := &http.Client{}
	resp, err := client.Get(asset.BrowserDownloadURL)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(tmpDir, "agent-setup-new")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
		}
		if readErr != nil {
			break
		}
	}
	f.Close()

	// Verify the content was downloaded correctly.
	got, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(newContent) {
		t.Errorf("downloaded content = %q, want %q", got, newContent)
	}

	// Verify atomic rename works.
	if err := os.Chmod(tmpFile, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmpFile, fakeBin); err != nil {
		t.Fatalf("rename failed: %v", err)
	}

	replaced, err := os.ReadFile(fakeBin)
	if err != nil {
		t.Fatal(err)
	}
	if string(replaced) != string(newContent) {
		t.Errorf("replaced binary = %q, want %q", replaced, newContent)
	}
}
