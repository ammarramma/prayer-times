# Progress Tracking - Prayer Times TUI

**Status Keseluruhan**: 🟡 Fase 1 direvisi: core worldwide dan uji akurasi selesai; CLI terhubung. Fase 2–4 menyusul sesuai prioritas.

---

## Fase 1: Core Logic (Worldwide + Kemenag) 🟡
**Tujuan**: Kalkulasi pure-domain yang akurat untuk kota dunia, termasuk koreksi ihtiyati Kemenag, sebelum UI.

### Checklist
- [x] Pisahkan logika `internal/prayer`, default/struktur di `internal/config`, dan wiring CLI tipis
- [x] Model Location dengan nama, lintang, bujur, dan timezone IANA; default Lhokseumawe (5.1801, 97.1406, Asia/Jakarta)
- [x] Preset Kemenag (Fajr 20°, Isha 18°) serta MWL, Umm al-Qura, Egyptian, Karachi, ISNA, Singapore, Dubai, Moonsighting Committee, Kuwait, Qatar, UOIF
- [x] Madhab Shafii/Hanafi dan aturan lintang tinggi Middle of the Night, Seventh of the Night, Twilight Angle
- [x] Ihtiyati configurable (default 2 menit hanya Kemenag): Subuh/Dzuhur/Ashar/Maghrib/Isya +2, Terbit -2; opsional Imsak = Subuh -10
- [x] Pembulatan ke menit ke atas ketika detik atau nanodetik bukan nol
- [x] Flag lokasi/metode/madhab/aturan lintang/ihitiyati/tanggal pada CLI; output menampilkan nama lokasi, metode, madhab, dan status koreksi
- [x] Label output Bahasa Indonesia dipusatkan di `internal/i18n`
- [x] Unit test jadwal Kemenag 3 tanggal, ihtiyati non-Kemenag, Jakarta/London/Makkah, Oslo musim panas, pembulatan, dan Imsak
- [x] `go build ./...`, `go vet ./...`, `go test ./...` (verifikasi setelah implementasi)
- [x] Catat keputusan API adhango dan sumber fixture Kemenag

**Status**: Implementasi fase direvisi. Penjadwalan berikutnya dimulai dari Fase 2 setelah rangkaian verifikasi dan commit fase ini.

**Catatan keputusan & akurasi**:
- `adhango v0.1.0` menyediakan preset metode, madhab, high-latitude rule, dan adjustment per shalat. API kalkulasi membulatkan ke menit terdekat; implementasi domain memakai komponen waktu astronomy dan membulatkan **ke atas** setelah adjustment agar syarat Kemenag terpenuhi.
- Kemenag menggunakan Fajr 20° dan Isha 18°; ihtiyati manual ditambahkan hanya untuk metode Kemenag. Preset lain mempertahankan koreksi bawaan masing-masing, tanpa ihtiyati aplikasi.
- Fixture jadwal Oktober 2026 bersumber dari tabel jadwal Lhokseumawe Tirto yang menyebut Bimas Islam/Kemenag sebagai sumber: [Tirto — jadwal Lhokseumawe](https://jadwal-sholat.tirto.id/kota-lhokseumawe). Situs resmi SIHAT bersifat dinamis; pembaca halaman tidak dapat mengambil datanya secara langsung. Fixture per tanggal diverifikasi terhadap sumber publik tersebut dengan toleransi 1 menit.
- `adhango` tidak mempertahankan timezone input untuk hasil astronomi; tanggal kalender lokal dipakai sebagai tanggal Gregorian setempat, dan hasil dikonversi ke timezone lokasi sebelum pembulatan waktu sipil.

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

**Last Updated**: 2026-10-28
