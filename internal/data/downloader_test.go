package data

import (
	"path/filepath"
	"testing"
)

func TestBundledDataIntegrity(t *testing.T) {
	dataDir, err := filepath.Abs(filepath.Join("..", "..", "data"))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("REALITYCHECK_DATA_DIR", dataDir)
	if err := NewDownloader().EnsureDataFiles(); err != nil {
		t.Fatalf("bundled data verification failed: %v", err)
	}
}
