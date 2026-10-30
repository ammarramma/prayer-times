package notify

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StatePath returns the notification state file, following XDG_STATE_HOME.
func StatePath() (string, error) {
	if base := os.Getenv("XDG_STATE_HOME"); base != "" {
		return filepath.Join(base, "prayer-times", "notified.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locating home directory: %w", err)
	}
	return filepath.Join(home, ".local", "state", "prayer-times", "notified.json"), nil
}

// LoadState reads the last delivered notification. A missing file is not an
// error and returns a nil state.
func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading notification state: %w", err)
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing notification state %s: %w", path, err)
	}
	return &state, nil
}

// SaveState writes the last delivered notification atomically with owner-only
// permissions.
func SaveState(path string, state State) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating notification state directory: %w", err)
	}
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encoding notification state: %w", err)
	}
	temp, err := os.CreateTemp(dir, ".notified-*.json")
	if err != nil {
		return fmt.Errorf("creating notification state file: %w", err)
	}
	tempName := temp.Name()
	defer func() { _ = os.Remove(tempName) }()
	if err := temp.Chmod(0o600); err != nil {
		_ = temp.Close()
		return fmt.Errorf("setting notification state permissions: %w", err)
	}
	if _, err := temp.Write(append(data, '\n')); err != nil {
		_ = temp.Close()
		return fmt.Errorf("writing notification state: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("closing notification state file: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("replacing notification state: %w", err)
	}
	return nil
}
