package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ammarramma/prayer-times/internal/prayer"
	"github.com/spf13/viper"
)

// Config stores the application location, calculation options, and locale.
type Config struct {
	Location       prayer.Location `json:"location" mapstructure:"location"`
	Options        prayer.Options  `json:"options" mapstructure:"options"`
	Language       string          `json:"language" mapstructure:"language"`
	MethodExplicit bool            `json:"-" mapstructure:"-"`
}

// Default returns the initial Kemenag configuration for Lhokseumawe.
func Default() Config {
	return Config{
		Location: prayer.DefaultLocation(),
		Options:  prayer.DefaultOptions(prayer.MethodKemenag),
		Language: "en",
	}
}

// DefaultPath returns the standard user config path.
func DefaultPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locating user configuration directory: %w", err)
	}
	return filepath.Join(base, "prayer-times", "config.json"), nil
}

// Load reads a Viper JSON config file when present; CLI flags may override its values.
func Load(path string) (Config, error) {
	cfg := Default()
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			return cfg, nil
		}
		return Config{}, fmt.Errorf("reading config %s: %w", path, err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config %s: %w", path, err)
	}
	cfg.MethodExplicit = v.IsSet("options.method")
	mergeDefaults(&cfg)
	return cfg, nil
}

// Save writes the JSON configuration to path and creates its directory as needed.
func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	v.Set("location", cfg.Location)
	v.Set("options", cfg.Options)
	v.Set("language", cfg.Language)
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("writing config %s: %w", path, err)
	}
	return os.Chmod(path, 0o600)
}

func mergeDefaults(cfg *Config) {
	defaults := Default()
	if cfg.Location.Name == "" {
		cfg.Location.Name = defaults.Location.Name
	}
	if cfg.Location.Timezone == "" {
		cfg.Location.Timezone = defaults.Location.Timezone
	}
	if cfg.Location.Latitude == 0 && cfg.Location.Longitude == 0 {
		cfg.Location.Latitude, cfg.Location.Longitude = defaults.Location.Latitude, defaults.Location.Longitude
	}
	if cfg.Language == "" {
		cfg.Language = defaults.Language
	}

	if cfg.Options.Method == "" {
		cfg.Options.Method = prayer.DefaultMethod(cfg.Location.Timezone)
		cfg.Options.IhtiyatiMinutes = prayer.DefaultOptions(cfg.Options.Method).IhtiyatiMinutes
	}
	if cfg.Options.Madhab == "" {
		cfg.Options.Madhab = defaults.Options.Madhab
	}
	// An unset high-latitude rule stays empty: the caller resolves it from the
	// location latitude unless the user chose a rule explicitly.
}
