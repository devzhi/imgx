package load

import "testing"

func TestDockerLoadCommandQuotesArguments(t *testing.T) {
	got := dockerLoadCommand("/usr/bin/docker", "/tmp/nginx latest.tar.gz")
	want := "sudo '/usr/bin/docker' load -i '/tmp/nginx latest.tar.gz'"
	if got != want {
		t.Fatalf("docker load command = %q, want %q", got, want)
	}
}

func TestShellQuoteEscapesSingleQuotes(t *testing.T) {
	got := shellQuote("a'b")
	want := `'a'"'"'b'`
	if got != want {
		t.Fatalf("quoted string = %q, want %q", got, want)
	}
}
