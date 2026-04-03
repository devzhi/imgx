package load

import "testing"

func TestBuildRemotePathUsesBaseName(t *testing.T) {
	got := buildRemotePath("/tmp/imgx", `C:\images\nginx latest.tar.gz`)
	want := "/tmp/imgx/nginx latest.tar.gz"
	if got != want {
		t.Fatalf("remote path = %q, want %q", got, want)
	}
}
