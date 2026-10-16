# Progress Tracking - Prayer Times TUI

**Status Keseluruhan**: 🟢 Fase 1 Selesai - Siap ke Fase 2

---

## Fase 1: Core Logic (Terminal Biasa Dulu) ✅
**Tujuan**: Pastikan kalkulasi waktu shalat 100% akurat sebelum sentuh UI

### Checklist
- [x] Setup project Go (go mod init)
- [x] Implementasi fungsi kalkulasi waktu shalat
  - [x] Hardcode koordinat Lhokseumawe (Lat 5.1801, Long 97.1406)
  - [x] Kalkulasi 5 waktu shalat (Subuh, Dzuhur, Ashar, Maghrib, Isya)
  - [x] Support metode Kemenag
- [x] Print ke terminal dengan fmt.Println
- [x] Verifikasi akurasi dengan jadwal resmi Kemenag
- [x] Commit: "Fase 1: Core logic waktu shalat"

**Status**: ✅ Selesai (Commit: d4e6bf2)

**Notes**:
- Menggunakan library `adhango` (high precision, Jean Meeus algorithm)
- Parameter Kemenag: Fajr 20°, Isha 18°
- Madhab Shafii untuk Asr
- Verifikasi dengan jadwalsholat.org: selisih ~2 menit (dalam toleransi)

---

## Fase 2: Interaktivitas & State (Bubble Tea) 🔵
**Tujuan**: Implementasi arsitektur ELM dan navigasi hari

### Checklist
- [ ] Install Bubble Tea (`go get github.com/charmbracelet/bubbletea`)
- [ ] Buat struktur Model
  - [ ] Field untuk menyimpan jadwal hari ini
  - [ ] Field untuk tracking tanggal yang ditampilkan
- [ ] Implementasi fungsi Init()
- [ ] Implementasi fungsi Update()
  - [ ] Handle key press (Arrow Left: kemarin, Arrow Right: besok)
  - [ ] Handle quit (q, Ctrl+C)
- [ ] Implementasi fungsi View() (render sederhana dulu)
- [ ] Test navigasi maju-mundur hari
- [ ] Commit: "Fase 2: Bubble Tea & navigasi hari"

**Status**: Belum Dimulai

---

## Fase 3: Make it Beautiful (Lipgloss) 🔵
**Tujuan**: Styling dan highlight waktu shalat berikutnya

### Checklist
- [ ] Install Lipgloss (`go get github.com/charmbracelet/lipgloss`)
- [ ] Design layout UI
  - [ ] Header (tanggal, lokasi)
  - [ ] Grid waktu shalat (5 box)
  - [ ] Footer (help text)
- [ ] Implementasi highlight waktu shalat berikutnya
  - [ ] Logic untuk deteksi waktu shalat mana yang akan datang
  - [ ] Apply style hijau pada box yang aktif
- [ ] Tambahkan progress bar countdown ke waktu shalat berikutnya
- [ ] Polish colors & spacing
- [ ] Commit: "Fase 3: Lipgloss styling & highlight"

**Status**: Belum Dimulai

---

## Fase 4: Konfigurasi User 🔵
**Tujuan**: User bisa set lokasi & metode kalkulasi sendiri

### Checklist
- [ ] Install Viper (`go get github.com/spf13/viper`)
- [ ] Buat struktur config
  - [ ] Path: `~/.config/prayer-times/config.json`
  - [ ] Fields: latitude, longitude, calculation_method
- [ ] Implementasi read config dengan fallback ke default (Lhokseumawe)
- [ ] Buat command untuk set config (misal: `prayer-times config set`)
- [ ] Update core logic untuk gunakan config user
- [ ] Dokumentasi config di README
- [ ] Commit: "Fase 4: User configuration dengan Viper"

**Status**: Belum Dimulai

---

## Fase 5: Build & Publish ke AUR 🔵
**Tujuan**: Release aplikasi ke AUR

### Checklist
- [ ] Cleanup & final testing
- [ ] Buat README.md lengkap
  - [ ] Cara install
  - [ ] Cara konfigurasi
  - [ ] Screenshots/Demo
- [ ] Push ke GitHub
- [ ] Buat GitHub Release
  - [ ] Tag version (v1.0.0)
  - [ ] Build binary untuk release
  - [ ] Upload binary ke release
- [ ] Tulis PKGBUILD
  - [ ] Test build lokal dengan makepkg
  - [ ] Verify PKGBUILD valid
- [ ] Submit ke AUR
  - [ ] Buat AUR account jika belum ada
  - [ ] Push PKGBUILD ke AUR git
- [ ] Test install: `yay -S prayer-times`
- [ ] Commit: "Fase 5: AUR release ready"

**Status**: Belum Dimulai

---

## Notes
- Selalu test manual setiap fase sebelum commit
- Jangan skip fase - ini bukan race
- Kalau ada bug, fix di fase yang sama, jangan bawa ke fase berikutnya
- Recovery point: setiap commit fase adalah safe point

---

**Last Updated**: 2026-10-16
