package config

import "github.com/ammarramma/prayer-times/internal/prayer"

// Config menyimpan lokasi serta opsi kalkulasi aplikasi.
type Config struct {
	Location prayer.Location `json:"location" mapstructure:"location"`
	Options  prayer.Options  `json:"options" mapstructure:"options"`
}

// Default mengembalikan konfigurasi awal Kemenag di Lhokseumawe.
func Default() Config {
	return Config{
		Location: prayer.DefaultLocation(),
		Options:  prayer.DefaultOptions(prayer.MethodKemenag),
	}
}
