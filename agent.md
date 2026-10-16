# Agent Context - Prayer Times TUI

## Tentang Proyek
Aplikasi TUI (Text User Interface) untuk menampilkan jadwal waktu shalat yang akurat sesuai standar Kemenag Indonesia. Proyek ini akan dirilis ke AUR (Arch User Repository) untuk pengguna Arch Linux dan turunannya.

## Tech Stack
- **Bahasa**: Go
- **TUI Framework**: Bubble Tea (arsitektur ELM: Model, Update, View)
- **Styling**: Lipgloss (Charmbracelet)
- **Config Management**: Viper
- **Target Platform**: Linux (khususnya Arch Linux via AUR)

## Lokasi Referensi
- **Lokasi Default**: Lhokseumawe
  - Latitude: 5.1801
  - Longitude: 97.1406
- **Standar Akurasi**: Kemenag Indonesia

## Arsitektur
```
prayer-times/
├── agent.md          # Konteks untuk AI agent (file ini)
├── langkah.md        # Tracking progress per fase
├── main.go           # Entry point aplikasi
├── internal/
│   ├── prayer/       # Core logic kalkulasi waktu shalat
│   ├── tui/          # Bubble Tea components
│   └── config/       # Config management (Viper)
├── go.mod
├── go.sum
└── PKGBUILD          # (Fase 5) Package build script untuk AUR
```

## Prinsip Development
1. **Iteratif per Fase**: Setiap fase harus selesai sempurna dan di-commit sebelum lanjut ke fase berikutnya
2. **Akurasi First**: Logika waktu shalat harus 100% akurat sebelum UI disentuh
3. **Clean Commits**: Setiap fase memiliki commit tersendiri agar mudah di-recover
4. **Documentation**: Selalu update langkah.md setiap menyelesaikan sub-task

## Notes untuk Agent
- Gunakan bahasa Indonesia untuk dokumentasi dan commit messages
- Saat fase 1, jangan sentuh UI sama sekali - fokus pada akurasi kalkulasi
- Model lemah bisa refer ke file ini untuk context proyek
- Setiap fase harus verified working sebelum lanjut
