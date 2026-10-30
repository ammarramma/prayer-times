# Progress Tracking - Prayer Times TUI

**Status Keseluruhan**: ✅ Stabilisasi global dan dokumentasi akurasi selesai (2026-10-29), ditambah penyempurnaan rilis, default lintang tinggi, notifikasi desktop, dan persiapan rilis v0.1.1 (2026-10-30). Fase terdahulu dipertahankan sebagai riwayat; lihat bagian stabilisasi global, penyempurnaan rilis, notifikasi desktop, rilis v0.1.1, dan open items.

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
- Fixture jadwal Oktober 2026 bersumber dari tabel jadwal Lhokseumawe Tirto yang menyebut Bimas Islam/Kemenag sebagai sumber: [Tirto — jadwal Lhokseumawe](https://jadwal-sholat.tirto.id/kota-lhokseumawe). Situs resmi SIHAT bersifat dinamis; pembaca halaman tidak dapat mengambil datanya secara langsung. Fixture yang tepat memakai toleransi 0; Ashar 2026-10-16 dan Maghrib 2026-10-28 tetap ±1 menit karena selisih yang belum dapat dijelaskan.
- Selisih terhadap fixture untuk Fajr/Ashar/Maghrib/Isha: 2026-10-01 = 0/0/0/0 menit; 2026-10-16 = 0/−1/0/0 menit; 2026-10-28 = 0/0/−1/0 menit. Fajr dan Isha lulus ketat 0; Ashar 2026-10-16 dan Maghrib 2026-10-28 dibiarkan ±1 menit. Rumus Maghrib tidak diubah. Terbit tak termasuk tiga fixture tersebut; observasi nonresmi 2026-10-29 di bawah menunjukkan selisih yang masih perlu data resmi.
- `adhango` tidak mempertahankan timezone input untuk hasil astronomi; tanggal kalender lokal dipakai sebagai tanggal Gregorian setempat, dan hasil dikonversi ke timezone lokasi sebelum pembulatan waktu sipil.

---

## Stabilisasi global (2026-10-29)

### Checklist
- [x] Titik aman: baseline repository dikomit setelah `go build ./...`, `go vet ./...`, `go test ./...` lulus.
- [x] Snapshot baseline 12 metode × 6 kota × 2 tanggal di `internal/prayer/testdata/golden.json`; tes mencegah perubahan hasil menit.
- [x] Status NOW/NEXT murni berbasis waktu jadwal keluaran, jendela 15 menit, box NEXT biru selama NOW, dan teks eksplisit.
- [x] Locale `en` default, `id` opsional via `-lang` atau config.
- [x] Profil metode dan offset pengguna; default metode menurut zona; profil non-Kemenag tidak mengubah golden.
- [x] Tes fixture Kemenag ketat untuk semua nilai cocok; dua selisih tetap ±1 menit dengan komentar.
- [x] Tes DST London/Sydney serta Tromsø untuk seluruh metode.
- [x] Dokumentasi akurasi Inggris dan README Inggris.
- [x] `go build ./...`, `go vet ./...`, dan `go test ./...` lulus setelah stabilisasi.
- [x] Tolak argumen posisional tak terduga dengan pesan CLI yang jelas.

**Open items — belum diputuskan tanpa data resmi**

- Selisih Maghrib dan Terbit Lhokseumawe tidak ditambal; rumus Maghrib/Terbit tidak diubah. Observasi tidak resmi 2026-10-29 (bukan fixture): aplikasi Android: Imsak 04:50, Subuh 05:00, Terbit 06:17, Dhuha 06:41, Dzuhur 12:19, Ashar 15:38, Maghrib 18:17, Isya 19:27. Google, basis “Kemenag Jakarta Pusat”, tanpa ihtiyati: Fajr 04:58, Sunrise 06:17, Dhuhr 12:15, Asr 15:36, Maghrib 18:14, Isha 19:24. Proyek: Sunrise 06:15 dan Maghrib 18:16.
- Data yang diperlukan: tabel resmi Bimas Islam untuk Lhokseumawe sekurang-kurangnya 5 tanggal pada Oktober 2026, termasuk rincian Terbit/Maghrib. Hipotesis tentang perbedaan koordinat, zona, atau konvensi koreksi belum terbukti.
- Sumber preset metode adhango: API/preset dependency lokal; verifikasi terhadap otoritas lokal belum dilakukan. Dicatat unverified kecuali Kemenag partial.
- Positional argumen ditolak oleh CLI dengan pesan Inggris; awalnya diminta sebagai kerja opsional, tetapi ditambahkan sebagai perlindungan langsung terhadap kasus argumen `.` yang dijelaskan di brief.

**Keputusan stabilisasi**: baseline saat ini memiliki dua file tambahan tracked (`digest.txt`, `prayer-times-tui-0.1.0.tar.gz`); keduanya dipertahankan dan diamankan dalam commit baseline, bukan dibuang sebagai perubahan pengguna. Snapshot tanggal dipilih 2026-10-28 dan 2026-06-21 untuk mencakup musim dan tetap bisa direproduksi.

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
- `sha256sums=('SKIP')` sengaja sementara dengan komentar di PKGBUILD; setelah tag v0.1.0 tersedia di GitHub, jalankan `updpkgsums` dan commit pembaruan PKGBUILD/.SRCINFO.
- Artefak `digest.txt` dan `prayer-times-tui-0.1.0.tar.gz` dihapus dari tracking lewat `git rm --cached`; tarball build lokal dihapus dari disk, sementara `digest.txt` tetap ada di disk. Keduanya dan `.cli-integration-*/` diabaikan oleh `.gitignore`.
- PKGBUILD lokal diuji terhadap tarball sementara `prayer-times-0.1.0.tar.gz` di `/tmp`, dibuat dari `git archive HEAD` ditambah perubahan working tree yang diperlukan untuk menyamakan `main.go` dan TUI. URL file lokal hanya dipakai pada salinan PKGBUILD sementara; PKGBUILD root tetap memakai URL tag GitHub. `makepkg -f --noconfirm` lulus untuk x86_64 dan menjalankan `go test ./...` pada tahap check.
- `namcap` belum terpasang sehingga pemeriksaan PKGBUILD/paket dilewati; jalankan manual ketika tersedia.
- Placeholder gambar README kini berupa komentar HTML yang menyebut `docs/screenshot.png` menyusul; screenshot asli belum disediakan dan tidak dibuat palsu.
- Versi Go repo Arch memenuhi syarat saat dicek; `go.mod` tidak diturunkan.
- Binary lintas kompilasi dibuat untuk Linux amd64 dan arm64 dengan CGO dinonaktifkan. `dist/` dihapus setelah uji lokal agar artefak tidak ikut commit.

---

## Penyempurnaan rilis (2026-10-30)

### Checklist
- [x] README Inggris: bagian Install memuat `go install`, arsip rilis Linux amd64/arm64 beserta `SHA256SUMS`, Arch via clone lalu `makepkg -si` dengan PKGBUILD repo, dan status AUR; kalimat "Install the latest Go version:" menjadi "Install with Go:"; placeholder screenshot berupa komentar HTML agar tidak ada tautan gambar rusak
- [x] `docs/accuracy.md` diperiksa dengan `file`, pencarian byte NUL, dan validasi UTF-8: berkas sudah teks UTF-8 biasa tanpa NUL sehingga tidak perlu ditulis ulang
- [x] `main_test.go`: helper `runCLIResult` yang mengembalikan exit code dan output, plus test argumen posisional `.`, `-version` satu baris, `-lang id`, dan `-lang fr`
- [x] `main.go`: impor `_ "time/tzdata"`; `-version` memakai versi build info bila ldflags masih `dev`; doc comment `DefaultMethod` dipisahkan dari akhir `DefaultOptions`
- [x] `.github/workflows/ci.yml` pada push dan pull_request: `go build`, `go vet`, `go test ./...` dengan `actions/setup-go` (`go-version-file: go.mod`), tanpa secret dan tanpa langkah deploy
- [x] Default aturan lintang tinggi dari lintang: `seventh_of_the_night` bila |lintang| > 48 dan pengguna tidak memilih eksplisit, selain itu `middle_of_the_night`; test London 2026-06-21; README dan `docs/accuracy.md` diperbarui
- [x] `gofmt -l .` kosong, `go build ./...`, `go vet ./...`, `go test ./... -count=1` lulus setelah tiap langkah

**Catatan keputusan**:
- Keeksplisitan aturan lintang tinggi dilacak dengan nilai aturan kosong (belum dipilih) vs non-kosong (dipilih), sejalan dengan `Method` kosong dan tanpa flag tambahan. `DefaultOptions` tidak lagi mengisi aturan, `config.mergeDefaults` membiarkannya kosong, dan `main.go` menyimpan hasil resolusi saat `-set-config`. Flag `-high-latitude` atau `options.high_latitude_rule` yang terisi selalu menang.
- Golden London berubah pada kedua tanggal dan hanya pada Fajr/Isya, sedangkan Lhokseumawe, Jakarta, Makkah, Tokyo, dan Sydney identik byte demi byte. Karena hanya London yang berubah, langkah opsional tidak perlu dibatalkan.
- Default lintang tinggi diletakkan di domain (`DefaultHighLatitudeRule` dipanggil `Calculate`) agar pemanggil langsung dan golden memakai perilaku yang sama.
- Test London 2026-06-21 memastikan Isya−Maghrib dan Terbit−Subuh masing-masing di bawah 2 jam 30 menit dengan opsi default.

---

---

## Notifikasi desktop (2026-10-30)

### Checklist
- [x] Domain `internal/notify`: `Events` (lima shalat hari ini + Subuh besok), `Due` (lead/grace window, dedup state, pilih event terbaru), `Settings.Validate`, serta test sintetis dan kalkulasi nyata (London DST, Oslo musim panas)
- [x] `notifications` pada config Viper dengan default lima shalat, lead 0, dan grace 5 menit; file lama tetap memakai default tersebut
- [x] Label notifikasi `en`/`id` di `internal/i18n` lewat `Locale.Notification` dan `Locale.PrayerLabel`
- [x] Adapter D-Bus `org.freedesktop.Notifications` memakai `github.com/godbus/dbus/v5` dengan koneksi privat ke session bus
- [x] Mode `-notify-once` plus `-notify-prayers`, `-notify-lead`, `-notify-grace`, dan `-dry-run`; `-set-config` menyimpan pengaturan notifikasi
- [x] State dedup di `$XDG_STATE_HOME/prayer-times/notified.json` (default `~/.local/state`), tulisan atomik 0600 dan direktori 0700
- [x] Unit `contrib/systemd/user/prayer-times-notify.service` (oneshot) dan `.timer` per menit dengan `Persistent=true`
- [x] README Inggris: bagian Desktop notifications, flag baru, dan blok `notifications` pada contoh config
- [x] `gofmt -l .` kosong, `go build ./...`, `go vet ./...`, `go test ./... -count=1` lulus setelah tiap langkah

**Catatan keputusan**:
- Notifikasi memakai D-Bus langsung tanpa fallback `notify-send`, sehingga tidak ada dependency eksternal tambahan dan perilaku sama di X11 maupun Wayland.
- Kegagalan bus atau daemon hanya menjadi peringatan dan exit 0 agar timer systemd tidak masuk status gagal; state juga tidak ditulis bila pengiriman gagal.
- Hanya event terbaru yang jatuh tempo dikirim, sehingga suspend panjang tidak menghasilkan tumpukan notifikasi lama.
- `PKGBUILD` dan `.SRCINFO` tidak diubah; unit dipasang manual lewat dokumentasi README. `systemd-analyze verify` hanya melaporkan `/usr/bin/prayer-times` belum ada di lingkungan ini, bukan kesalahan sintaks unit.
- Tromsø pada 21 Juni tidak bisa dipakai sebagai fixture notifikasi karena kalkulasi mengembalikan NaN saat matahari tidak terbenam; test domain memakai Oslo yang masih memiliki terbenam dan tetap > 48° lintang.

---

## Rilis v0.1.1 (2026-10-30)

### Checklist
- [x] `CHANGELOG.md` dengan entri 0.1.0 dan 0.1.1 (Added/Changed/Fixed)
- [x] README: contoh unduhan memakai `version=v0.1.1` dan catatan arsip arm64, tidak lagi mengunci v0.1.0
- [x] Arsip lokal v0.1.1 amd64/arm64 + `SHA256SUMS` dibuat lewat `scripts/build-release.sh`; `-version` mencetak `v0.1.1`, binary statis x86-64, dan jadwal 2026-10-29 cocok dengan fixture
- [x] `gofmt -l .` kosong, `go build ./...`, `go vet ./...`, `go test ./... -count=1` lulus
- [ ] Tag `v0.1.1`, push branch dan tag, GitHub Release beserta unggahan arsip — dilakukan pemilik repo
- [ ] `PKGBUILD` (`pkgver=0.1.1`, checksum nyata) dan `.SRCINFO` — dilakukan pemilik repo agar prosesnya dipahami; agent tidak menyentuhnya
- [ ] Submit AUR — menunggu akun AUR

**Catatan keputusan**:
- Checksum di `dist/SHA256SUMS` lokal hanya untuk verifikasi; checksum PKGBUILD harus diambil dari tarball tag GitHub (`updpkgsums`) karena arsip rilis memuat metadata tar yang tidak dijamin identik antar build.
- Arsip sebaiknya dibangun ulang dari checkout bertag yang bersih sebelum diunggah; arsip lokal saat ini dipakai untuk membuktikan skrip rilis dan perilaku binary.
- Completion shell dan subcommand (`-t`, `-tC`, `--tui`) ditunda ke v0.2.0 bersama adopsi cobra/pflag dan shim kompatibilitas `-nama`.
- Hanya `PKGBUILD` dan `.SRCINFO` yang tersisa untuk rilis paket Arch; keduanya diperbarui manual oleh pemilik repo.

---

## Notes
- Selalu test manual setiap fase sebelum commit
- Jangan skip fase - ini bukan race
- Kalau ada bug, fix di fase yang sama, jangan bawa ke fase berikutnya
- Recovery point: setiap commit fase adalah safe point

---

**Last Updated**: 2026-10-30
