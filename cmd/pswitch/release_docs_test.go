package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseWorkflowTargetsVersionTags(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	text := string(content)
	if !strings.Contains(text, "tags:") {
		t.Fatalf("release workflow does not define tag trigger")
	}
	if !strings.Contains(text, "- \"v*\"") {
		t.Fatalf("release workflow does not target version tags")
	}
}

func TestReleaseWorkflowBuildsExpectedArchives(t *testing.T) {
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	text := string(content)
	for _, snippet := range []string{
		"pswitch_${version}_${{ matrix.goos }}_${{ matrix.goarch }}",
		"dist/*.tar.gz",
		"dist/*.zip",
		"checksums.txt",
	} {
		if !strings.Contains(text, snippet) {
			t.Fatalf("release workflow missing %q", snippet)
		}
	}
}
