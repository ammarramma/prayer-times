# Agent Context - Prayer Times TUI

## Tentang Proyek
Aplikasi jadwal waktu shalat global dengan koordinat dan timezone IANA. UI dan output CLI memakai bahasa Inggris sebagai default; `-lang id` memilih Indonesia. Kemenag dipilih otomatis untuk zona waktu Indonesia yang ditentukan, sedangkan zona lain memakai MWL bila metode tidak diatur eksplisit. Proyek Go ini memakai Bubble Tea, Lipgloss, Viper, adhango, dan godbus/dbus untuk notifikasi desktop.

## Tech Stack
- **Bahasa**: Go
- **TUI Framework**: Bubble Tea (arsitektur ELM: Model, Update, View)
- **Styling**: Lipgloss (Charmbracelet)
- **Config Management**: Viper (`~/.config/prayer-times/config.json`)
- **Target Platform**: Linux (khususnya Arch Linux via AUR)

## Lokasi dan Kalkulasi
- **Lokasi Default**: Lhokseumawe, Aceh (5.1801, 97.1406; `Asia/Jakarta`); nama, koordinat, dan timezone lain diberikan melalui config/flag.
- `internal/prayer.Location` memakai timezone IANA, bukan timezone mesin.
- **Metode**: Kemenag (Fajr 20°, Isha 18°), MWL, Umm al-Qura, Egyptian, Karachi, ISNA, Singapore, Dubai, Moonsighting Committee, Kuwait, Qatar, UOIF. Default tak eksplisit: Kemenag untuk `Asia/Jakarta`, `Asia/Pontianak`, `Asia/Makassar`, `Asia/Jayapura`; MWL selainnya. Default lokasi Lhokseumawe tetap Kemenag.
- **Ihtiyati Kemenag**: default +2 menit untuk Fajr, Asr, Maghrib, Isha; Dhuhr +3; Sunrise −2. Imsak opsional adalah Fajr −10 menit. `Options.Offsets` ditambahkan setelah preset untuk semua metode.
- Pembulatan waktu keluaran: detik/nanodetik lebih dari nol selalu dibulatkan naik ke menit berikutnya.
- Madhab Asr: Shafii (default) atau Hanafi. Aturan lintang tinggi: Seventh of the Night (default bila |lintang| > 48 dan aturan tidak dipilih eksplisit), Middle of the Night (default lainnya), atau Twilight Angle. Aturan kosong menandai "belum dipilih" dan diselesaikan `Calculate` dari lintang.
- `internal/i18n.Locale` mendukung English default dan Indonesian (`For("en"|"id")`).
- Status NOW/NEXT dihitung di `internal/prayer` dari waktu yang ditampilkan; NOW berlangsung 15 menit. Sunrise bukan target status.
- Notifikasi desktop: domain `internal/notify` (lima shalat hari ini + Subuh besok), jendela lead/grace (default 0/5 menit), dan dedup lewat `$XDG_STATE_HOME/prayer-times/notified.json`. `-notify-once` dipanggil timer systemd user tiap menit; pengiriman memakai `org.freedesktop.Notifications` dan kegagalan bus/daemon hanya menjadi peringatan exit 0. `-dry-run` mencetak jadwal JSON tanpa mengirim.

## Arsitektur
```
prayer-times/
├── agent.md          # Konteks untuk AI agent (file ini)
├── langkah.md        # Tracking progress per fase
├── main.go           # Entry point aplikasi dan flag versi
├── LICENSE           # Lisensi MIT
├── PKGBUILD          # Build package Arch Linux
├── .github/
│   └── workflows/ci.yml # CI: build, vet, dan test pada push/PR
├── scripts/
│   └── build-release.sh # Build arsip binary rilis linux/amd64 dan linux/arm64
├── internal/
│   ├── prayer/       # Core logic kalkulasi waktu shalat
│   ├── notify/       # Penjadwal notifikasi murni, adapter D-Bus, dan state
│   ├── tui/          # Bubble Tea components dan styling Lipgloss
│   ├── i18n/         # String label per locale
│   └── config/       # Struktur dan config JSON melalui Viper
├── contrib/systemd/user/ # Unit oneshot + timer notifikasi desktop
├── docs/accuracy.md  # Parameter profil, verifikasi, fixture, dan default lintang tinggi
├── go.mod
├── go.sum
└── .SRCINFO          # Metadata package Arch untuk AUR

Dist rilis `dist/` dibuat oleh `scripts/build-release.sh` dan tidak disimpan di Git. `main.go` menyediakan flag `-version`: nilai ldflags dipakai lebih dulu, lalu versi modul dari `debug.ReadBuildInfo()` (dipakai `go install`), dan `dev` bila keduanya kosong. Binary mengimpor `time/tzdata` sehingga zona IANA tetap tersedia tanpa zoneinfo sistem.
```

## Prinsip Development
1. **Iteratif per Fase**: Setiap fase harus selesai sempurna dan di-commit sebelum lanjut ke fase berikutnya
2. **Akurasi First**: Logika waktu shalat harus 100% akurat sebelum UI disentuh
3. **Clean Commits**: Setiap fase memiliki commit tersendiri agar mudah di-recover
4. **Documentation**: Selalu update langkah.md setiap menyelesaikan sub-task

## Notes untuk Agent
- README dan `docs/accuracy.md` berbahasa Inggris; `agent.md`, `langkah.md`, dan pesan commit berbahasa Indonesia.
- Saat fase 1, jangan sentuh UI sama sekali - fokus pada akurasi kalkulasi
- Model lemah bisa refer ke file ini untuk context proyek
- Setiap fase harus verified working sebelum lanjut
- Core domain murni berada di `internal/prayer` dan `internal/notify`; wiring entry point tetap tipis. Mode notifikasi memakai interface `notify.Notifier` agar dapat diuji tanpa session bus.
- Default lintang tinggi adalah keputusan domain (`DefaultHighLatitudeRule`), bukan hanya wiring CLI, supaya pemanggil `Calculate` langsung mendapat perilaku yang sama.
- Fixture akurasi Indonesia harus mencatat sumber/toleransi yang benar-benar diketahui. Selisih tak terjelaskan tidak boleh ditambal dengan perubahan rumus; rekam di `langkah.md` dan `docs/accuracy.md`.
- Snapshot golden menjaga baseline perhitungan dan bukan bukti akurasi. Metode non-Kemenag belum diverifikasi terhadap otoritas lokal.
