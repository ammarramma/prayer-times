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
- [x] Ihtiyati configurable (default 2 menit hanya Kemenag): Subuh/Ashar/Maghrib/Isya +2, Dzuhur ihtiyati dasar +1 (total +3), Terbit -2; opsional Imsak = Subuh -10
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
- Fixture jadwal Oktober 2026 bersumber dari tabel jadwal Lhokseumawe Tirto yang menyebut Bimas Islam/Kemenag sebagai sumber: [Tirto — jadwal Lhokseumawe](https://jadwal-sholat.tirto.id/kota-lhokseumawe). Situs resmi SIHAT bersifat dinamis; pembaca halaman tidak dapat mengambil datanya secara langsung. Fixture Subuh/Ashar/Maghrib/Isya memakai toleransi ±1 menit; Dzuhur wajib tepat (toleransi 0).
- Setelah koreksi ihtiyati Dzuhur menjadi dasar +1 menit, selisih terhadap fixture untuk Subuh/Ashar/Maghrib/Isya adalah: 2026-10-01 = 0/0/0/0 menit; 2026-10-16 = 0/−1/0/0 menit; 2026-10-28 = 0/0/−1/0 menit. Urutan kolom: Subuh/Ashar/Maghrib/Isya. Rumus Maghrib tidak diubah; selisih −1 menit pada 2026-10-28 tetap dicatat menunggu data jadwal resmi yang lebih kuat. Terbit tidak tercakup fixture ini, sehingga tidak ada selisih Terbit yang dapat dilaporkan atau dasar untuk mengubah rumusnya.
- `adhango` tidak mempertahankan timezone input untuk hasil astronomi; tanggal kalender lokal dipakai sebagai tanggal Gregorian setempat, dan hasil dikonversi ke timezone lokasi sebelum pembulatan waktu sipil.

---

## Fase 2: Interaktivitas & State (Bubble Tea) 🟡
**Tujuan**: Implementasi arsitektur ELM dan navigasi hari dengan konfigurasi lokasi/metode yang sama dari core.

### Checklist
- [x] Tambahkan Bubble Tea
- [x] Buat model yang menyimpan jadwal, tanggal, lokasi, dan opsi
- [x] Implementasi Init(), Update(), View()
- [x] Handle panah kiri/kanan untuk hari sebelumnya/berikutnya
- [x] Handle q dan Ctrl+C
- [x] Uji transisi navigasi dan quit
- [x] Verifikasi `go build ./...`, `go vet ./...`, `go test ./...`
- [x] Commit: "Fase 2: Bubble Tea & navigasi hari"

**Status**: ✅ Selesai (commit menyusul dicatat setelah pembuatan).

---

## Fase 3: Make it Beautiful (Lipgloss) 🟡
**Tujuan**: Styling, highlight waktu shalat berikutnya, dan countdown.

### Checklist
- [x] Install dan gunakan Lipgloss
- [x] Header tanggal/lokasi, 5 box shalat dan footer bantuan
- [x] Deteksi dan highlight shalat berikutnya
- [x] Countdown detik-ke-detik
- [x] Layout melebar/menyempit mengikuti ukuran terminal
- [x] Tes DOM string/layout lebar dan sempit
- [x] Countdown setelah Isya mengarah ke Subuh besok dan tetap positif
- [x] Sembunyikan label berikutnya/countdown saat menampilkan tanggal selain hari ini
- [x] Samakan tinggi box, tampilkan Terbit sebagai info, dan batasi lebar baris metadata
- [x] Tombol `t` kembali ke hari ini menurut timezone lokasi
- [x] Mode `-fullscreen` memakai alt screen dan memusatkan konten; inline tetap default
- [x] Jalankan `go test ./...`, `go vet ./...`, dan `go build ./...`
- [ ] Commit Fase 3

**Status**: Implementasi dan verifikasi selesai; commit menyusul.

**Catatan keputusan UI**:
- `BERIKUTNYA` diletakkan di bawah deretan box, bukan di dalam box highlight, agar semua box punya dimensi sama.
- Countdown dan highlight hanya muncul saat tanggal terpilih sama dengan hari ini di timezone lokasi. Setelah Isya, jadwal Subuh dihitung ulang untuk tanggal kalender besok di timezone lokasi.
- Baris metode/madhab/ihtiyati dibungkus oleh Lipgloss pada lebar terminal; Terbit tetap berupa baris info terpisah dan tidak menjadi target countdown.
- Mode inline tetap default. `-fullscreen` hanya mengaktifkan alt screen dan pemusatan menggunakan ukuran terminal yang diterima dari Bubble Tea.

---

## Fase 4: Konfigurasi User 🔵
**Tujuan**: User bisa menyimpan lokasi dan seluruh opsi kalkulasi sendiri.

### Checklist
- [x] Install Viper
- [x] Struktur config pada `~/.config/prayer-times/config.json` (path mengikuti `os.UserConfigDir`)
- [x] Simpan lokasi, metode, madhab, high-latitude rule, ihtiyati dan Imsak
- [x] Baca config dengan fallback default Lhokseumawe jika belum ada
- [x] Flags CLI menimpa config; `-set-config` menyimpan pilihan flag/config
- [x] Tulis file config dengan permission owner-only (0600)
- [x] Tes default, round-trip, path, permission
- [x] Jalankan `go build ./...`, `go vet ./...`, `go test ./...` dan uji CLI simpan/muat config
- [x] Commit: "Fase 4: User configuration dengan Viper"

**Status**: ✅ Selesai (commit menyusul dicatat setelah pembuatan).

---

## Fase 5: Persiapan rilis v0.1.0 🔵
**Tujuan**: Siapkan binary Linux, package Arch, lisensi, dan dokumentasi; unggah rilis atau submit AUR dilakukan manual oleh pemilik repo.

### Checklist
- [x] Tambahkan lisensi MIT (2026, Ammar Ramma), `.gitignore`, dan flag `-version`
- [x] README instalasi (Go, binary, PKGBUILD, AUR segera hadir), catatan akurasi, dan placeholder screenshot
- [x] Skrip `scripts/build-release.sh` menghasilkan linux/amd64, linux/arm64, dan SHA256SUMS; uji dengan v0.1.0-test
- [x] Verifikasi binary amd64: `-version` dan jadwal `-date 2026-10-29`
- [x] Tulis PKGBUILD dan hasilkan `.SRCINFO`
- [x] Tes package lokal dengan `makepkg -f` memakai tarball lokal sementara; tahap build/check/package lulus pada x86_64
- [ ] Jalankan namcap — tool tidak terpasang di lingkungan ini
- [x] Cek remote: `origin` = `https://github.com/ammarramma/prayer-times.git`
- [x] Cek Go repo Arch (`pacman -Si go`): `1.27.3-1` di `extra`, memenuhi `go 1.27.1` di go.mod
- [ ] Push branch, buat/push tag v0.1.0, dan buat GitHub Release — jangan dijalankan oleh agent; pemilik repo menjalankannya manual
- [ ] Submit ke AUR — menunggu pendaftaran akun AUR dibuka

**Catatan keputusan & keterbatasan**:
- Tidak ada perintah push, pembuatan tag, `gh release`, atau pengiriman lain yang dijalankan. Arsip lokal uji `v0.1.0-test` dihapus sebelum commit.
- `sha256sums=('SKIP')` sengaja sementara; setelah tag v0.1.0 tersedia di GitHub, jalankan `updpkgsums` dan commit pembaruan PKGBUILD/.SRCINFO.
- PKGBUILD lokal diuji terhadap tarball sementara `prayer-times-0.1.0.tar.gz` di `/tmp`, dibuat dari `git archive HEAD` ditambah perubahan working tree yang diperlukan untuk menyamakan `main.go` dan TUI. URL file lokal hanya dipakai pada salinan PKGBUILD sementara; PKGBUILD root tetap memakai URL tag GitHub. `makepkg -f --noconfirm` lulus untuk x86_64 dan menjalankan `go test ./...` pada tahap check.
- `namcap` belum terpasang sehingga pemeriksaan PKGBUILD/paket dilewati; jalankan manual ketika tersedia.
- Placeholder gambar README merujuk ke `docs/screenshot.png`; screenshot asli belum disediakan dan tidak dibuat palsu.
- Versi Go repo Arch memenuhi syarat saat dicek; `go.mod` tidak diturunkan.
- Binary lintas kompilasi dibuat untuk Linux amd64 dan arm64 dengan CGO dinonaktifkan. `dist/` dihapus setelah uji lokal agar artefak tidak ikut commit.

---

## Notes
- Selalu test manual setiap fase sebelum commit
- Jangan skip fase - ini bukan race
- Kalau ada bug, fix di fase yang sama, jangan bawa ke fase berikutnya
- Recovery point: setiap commit fase adalah safe point

---

**Last Updated**: 2026-10-29
