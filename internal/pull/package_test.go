package pull

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageWritesArchiveToRequestedDirectory(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "image")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte("[]"), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	outputDir := filepath.Join(tempDir, "dist")
	archivePath, err := Package(sourceDir, "library/nginx", "latest", "amd64", "linux", outputDir, nil)
	if err != nil {
		t.Fatalf("Package returned error: %v", err)
	}

	want := filepath.Join(outputDir, "library_nginx_latest_amd64_linux.tar.gz")
	if *archivePath != want {
		t.Fatalf("archive path = %q, want %q", *archivePath, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected archive to exist: %v", err)
	}
}
