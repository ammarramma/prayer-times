# Agent Context - Prayer Times TUI

## Tentang Proyek
Aplikasi jadwal waktu shalat worldwide yang mendukung kota mana pun dengan koordinat dan timezone IANA. Kemenag RI tersedia sebagai salah satu metode kalkulasi dan menjadi metode default yang menjaga akurasi untuk Indonesia. Proyek ini direncanakan dirilis ke AUR (Arch User Repository).

## Tech Stack
- **Bahasa**: Go
- **TUI Framework**: Bubble Tea (arsitektur ELM: Model, Update, View)
- **Styling**: Lipgloss (Charmbracelet)
- **Config Management**: Viper (`~/.config/prayer-times/config.json`)
- **Target Platform**: Linux (khususnya Arch Linux via AUR)

## Lokasi dan Kalkulasi
- **Lokasi Default**: Lhokseumawe, Aceh (5.1801, 97.1406; `Asia/Jakarta`); nama, koordinat, dan timezone lain diberikan melalui config/flag.
- `internal/prayer.Location` memakai timezone IANA, bukan timezone mesin.
- **Metode**: Kemenag (Fajr 20°, Isha 18°, default), MWL, Umm al-Qura, Egyptian, Karachi, ISNA, Singapore, Dubai, Moonsighting Committee, Kuwait, Qatar, UOIF.
- **Ihtiyati Kemenag**: default +2 menit untuk Subuh, Ashar, Maghrib, dan Isya; Dzuhur memakai ihtiyati dasar +1 menit (total +3 menit); Terbit dikurangi 2 menit. Imsak opsional adalah Subuh −10 menit. Metode lain tidak mendapat koreksi ihtiyati aplikasi.
- Pembulatan waktu keluaran: detik/nanodetik lebih dari nol selalu dibulatkan naik ke menit berikutnya.
- Madhab Asr: Shafii (default) atau Hanafi. Aturan lintang tinggi: Middle of the Night (default), Seventh of the Night, atau Twilight Angle.
- Label presentasi berada di `internal/i18n` agar mendukung lokalisasi.

## Arsitektur
```
prayer-times/
├── agent.md          # Konteks untuk AI agent (file ini)
├── langkah.md        # Tracking progress per fase
├── main.go           # Entry point aplikasi
├── internal/
│   ├── prayer/       # Core logic kalkulasi waktu shalat
│   ├── tui/          # Bubble Tea components dan styling Lipgloss
│   ├── i18n/         # String label per locale
│   └── config/       # Struktur dan config JSON melalui Viper
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
- Core domain murni berada di `internal/prayer`; wiring entry point tetap tipis.
- Fixture akurasi Indonesia harus memakai jadwal Kemenag yang dapat direproduksi dan mencatat sumber/toleransi di `langkah.md`. Metode non-Kemenag tidak diberi ihtiyati aplikasi.
