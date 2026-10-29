package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ammarramma/prayer-times/internal/prayer"
	"github.com/spf13/viper"
)

// Config menyimpan lokasi serta opsi kalkulasi aplikasi.
type Config struct {
	Location prayer.Location `json:"location" mapstructure:"location"`
	Options  prayer.Options  `json:"options" mapstructure:"options"`
	Language string          `json:"language" mapstructure:"language"`
}

// Default mengembalikan konfigurasi awal Kemenag di Lhokseumawe.
func Default() Config {
	return Config{
		Location: prayer.DefaultLocation(),
		Options:  prayer.DefaultOptions(prayer.MethodKemenag),
		Language: "en",
	}
}

// DefaultPath mengembalikan lokasi file config standar XDG.
func DefaultPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("mencari direktori konfigurasi user: %w", err)
	}
	return filepath.Join(base, "prayer-times", "config.json"), nil
}

// Load membaca file JSON Viper bila tersedia; opsi flag dapat menimpa hasilnya.
func Load(path string) (Config, error) {
	cfg := Default()
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			return cfg, nil
		}
		return Config{}, fmt.Errorf("membaca config %s: %w", path, err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("mem-parsing config %s: %w", path, err)
	}
	mergeDefaults(&cfg)
	return cfg, nil
}

// Save menulis konfigurasi JSON ke path dan membuat direktori bila diperlukan.
func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("membuat direktori config: %w", err)
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("json")
	v.Set("location", cfg.Location)
	v.Set("options", cfg.Options)
	v.Set("language", cfg.Language)
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("menulis config %s: %w", path, err)
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
		cfg.Options.Method = defaults.Options.Method
	}
	if cfg.Options.Madhab == "" {
		cfg.Options.Madhab = defaults.Options.Madhab
	}
	if cfg.Options.HighLatitudeRule == "" {
		cfg.Options.HighLatitudeRule = defaults.Options.HighLatitudeRule
	}
}
