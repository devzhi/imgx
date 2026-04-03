package pull

import (
	"strings"
	"testing"
)

func TestDownloadImageRejectsUnsupportedManifest(t *testing.T) {
	_, err := DownloadImage(&TokenResponse{}, &ManifestsResp{MediaType: "text/plain"}, "amd64", "linux", "nginx", "latest")
	if err == nil {
		t.Fatal("expected unsupported manifest type to return an error")
	}
	if !strings.Contains(err.Error(), "unsupported manifest type") {
		t.Fatalf("unexpected error: %v", err)
	}
}
