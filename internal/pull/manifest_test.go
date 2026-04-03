package pull

import "testing"

func TestIsManifestList(t *testing.T) {
	if !isManifestList("application/vnd.oci.image.index.v1+json") {
		t.Fatal("expected OCI image index media type to be treated as a manifest list")
	}
	if !isManifestList("application/vnd.docker.distribution.manifest.list.v2+json") {
		t.Fatal("expected Docker manifest list media type to be treated as a manifest list")
	}
	if isManifestList("application/vnd.docker.distribution.manifest.v2+json") {
		t.Fatal("expected single-image manifest media type not to be treated as a manifest list")
	}
}
