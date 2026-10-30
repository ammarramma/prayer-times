package notify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ammarramma/prayer-times/internal/prayer"
)

func TestStatePathHonorsXDGStateHome(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)
	path, err := StatePath()
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, "prayer-times", "notified.json"); path != want {
		t.Fatalf("StatePath = %s, want %s", path, want)
	}
}

func TestLoadStateMissingFileReturnsNil(t *testing.T) {
	state, err := LoadState(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("missing state should not be an error: %v", err)
	}
	if state != nil {
		t.Fatalf("state = %#v, want nil", state)
	}
}

func TestSaveAndLoadState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "notified.json")
	want := State{Date: "2026-10-28", Prayer: prayer.Maghrib}
	if err := SaveState(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := LoadState(path)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || *got != want {
		t.Fatalf("state = %#v, want %#v", got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("state permissions = %o, want 600", info.Mode().Perm())
	}
	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if dirInfo.Mode().Perm() != 0o700 {
		t.Errorf("state directory permissions = %o, want 700", dirInfo.Mode().Perm())
	}
}

func TestLoadStateRejectsInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notified.json")
	if err := os.WriteFile(path, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadState(path); err == nil {
		t.Fatal("invalid state should return an error")
	}
}
