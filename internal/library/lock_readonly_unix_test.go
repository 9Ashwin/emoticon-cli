//go:build !windows

package library

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadManifestUsesReadOnlyExistingLock(t *testing.T) {
	root := t.TempDir()
	library, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := library.WriteManifest(context.Background(), Manifest{SchemaVersion: 1, Collection: "personal", Items: []Item{}}); err != nil {
		t.Fatal(err)
	}
	stickerDir := filepath.Join(root, ".sticker")
	lockPath := filepath.Join(stickerDir, "write.lock")
	manifestPath := filepath.Join(root, ManifestName)
	t.Cleanup(func() {
		_ = os.Chmod(lockPath, 0o600)
		_ = os.Chmod(manifestPath, 0o600)
		_ = os.Chmod(stickerDir, 0o700)
		_ = os.Chmod(root, 0o700)
	})
	for path, mode := range map[string]os.FileMode{
		root:         0o500,
		stickerDir:   0o500,
		lockPath:     0o400,
		manifestPath: 0o400,
	} {
		if err := os.Chmod(path, mode); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := library.ReadManifest(context.Background()); err != nil {
		t.Fatalf("read-only manifest access failed: %v", err)
	}
}

func TestReadManifestDoesNotCreateLockDirectory(t *testing.T) {
	root := t.TempDir()
	library, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := library.ReadManifest(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".sticker")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read path created .sticker: %v", err)
	}
}
