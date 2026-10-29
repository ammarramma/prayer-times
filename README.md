# Prayer Times TUI 🕌

Aplikasi waktu shalat berbasis Go yang mendukung kota di seluruh dunia melalui koordinat dan timezone IANA. Kemenag RI tetap menjadi metode default, dengan koreksi ihtiyati untuk mempertahankan konvensi waktu Indonesia.

## Status proyek

Core worldwide dan kalkulasi selesai. Aplikasi membuka TUI Bubble Tea untuk navigasi tanggal; opsi `-date` tetap menghasilkan satu jadwal non-interaktif. Lihat [langkah.md](./langkah.md) untuk status detail.

## Fitur kalkulasi

- Lokasi dengan nama, lintang, bujur, timezone IANA. Default: Lhokseumawe (5.1801, 97.1406, `Asia/Jakarta`).
- Metode kalkulasi: Kemenag RI (Fajr 20°, Isha 18°), MWL, Umm al-Qura, Egyptian, Karachi, ISNA, Singapore, Dubai, Moonsighting Committee, Kuwait, Qatar, dan UOIF.
- Madhab Asr: Shafii (default) atau Hanafi.
- Aturan lintang tinggi: Middle of the Night, Seventh of the Night, Twilight Angle.
- Ihtiyati khusus Kemenag: default +2 menit untuk Subuh, Ashar, Maghrib, dan Isya; Dzuhur mendapat +3 menit (ihtiyati dasar +1), sedangkan Terbit dikurangi 2 menit. Untuk metode lain tidak ada koreksi ihtiyati dari aplikasi.
- Imsak opsional dihitung 10 menit sebelum Subuh.
- Waktu yang memiliki detik dibulatkan ke atas ke menit berikutnya.

## Instalasi

1. **Go:** pasang langsung dengan Go:

   ```bash
   go install github.com/ammarramma/prayer-times@latest
   ```

2. **Binary Linux:** unduh arsip untuk arsitektur `amd64` atau `arm64` dari [GitHub Releases](https://github.com/ammarramma/prayer-times/releases), ekstrak, lalu jalankan binary `prayer-times`.
3. **Arch Linux dari PKGBUILD:** setelah source tag rilis tersedia di GitHub:

   ```bash
   git clone https://github.com/ammarramma/prayer-times.git
   cd prayer-times
   makepkg -si
   ```

4. **AUR:** segera hadir (pendaftaran akun AUR sedang ditutup).

## Menjalankan

```bash
go run .
```

Di TUI gunakan panah kiri/kanan untuk memilih hari sebelumnya/berikutnya, `t` untuk kembali ke hari ini menurut timezone lokasi, dan `q` atau Ctrl+C untuk keluar. TUI berjalan inline secara default agar hasilnya tetap berada di scrollback. Gunakan `-fullscreen` untuk menjalankannya di alt screen dan menempatkan konten di tengah layar. Lokasi, metode, madhab, aturan lintang, ihtiyati, dan tanggal dapat diberikan lewat flag yang sama seperti contoh di bawah; seluruh opsi tetap memakai kalkulasi core yang sama.

Flag opsional:

```bash
go run . \
  -name London -lat 51.5072 -lon -0.1276 -tz Europe/London \
  -method mwl -madhab hanafi -high-latitude seventh_of_the_night \
  -date 2026-10-28
```

Daftar flag:

| Flag | Default | Keterangan |
| --- | --- | --- |
| `-config` | `~/.config/prayer-times/config.json` | Path config JSON; mengikuti direktori config OS |
| `-name` | `Lhokseumawe` | Nama kota/lokasi |
| `-lat` | `5.1801` | Lintang desimal, −90 sampai 90 |
| `-lon` | `97.1406` | Bujur desimal, −180 sampai 180 |
| `-tz` | `Asia/Jakarta` | Zona IANA (mis. `Europe/London`, `Asia/Riyadh`) |
| `-method` | `kemenag` | Salah satu metode di atas |
| `-madhab` | `shafii` | `shafii` atau `hanafi` |
| `-high-latitude` | `middle_of_the_night` | `middle_of_the_night`, `seventh_of_the_night`, `twilight_angle` |
| `-ihtiyati` | `2` | Menit koreksi khusus Kemenag; nilai metode lain diabaikan |
| `-imsak` | `false` | Tampilkan Imsak = Subuh −10 menit |
| `-fullscreen` | `false` | Jalankan TUI pada alt screen dan pusatkan konten |
| `-date` | Hari ini di zona lokasi | Tanggal Gregorian `YYYY-MM-DD`; menghasilkan output sekali jalan |
| `-set-config` | `false` | Simpan nilai config/flag ke config lalu keluar |

Buat atau perbarui konfigurasi melalui flags, misalnya:

```bash
go run . -set-config -name London -lat 51.5072 -lon -0.1276 -tz Europe/London -method mwl
```

Berkas config memakai JSON bertingkat, misalnya:

```json
{
  "location": { "name": "London", "latitude": 51.5072, "longitude": -0.1276, "timezone": "Europe/London" },
  "options": { "method": "mwl", "madhab": "shafii", "high_latitude_rule": "middle_of_the_night", "ihtiyati_minutes": 0, "include_imsak": false }
}
```

Flags yang ditentukan selalu mengesampingkan nilai file config; ketika file belum ada, default kembali ke Lhokseumawe.

Contoh metode lainnya: `umm_al_qura`, `egyptian`, `karachi`, `isna`, `singapore`, `dubai`, `moonsighting_committee`, `kuwait`, `qatar`, `uoif`.

## Akurasi

Jadwal default dihitung dengan metode Kemenag dan koreksi ihtiyati aplikasi. Hasilnya dapat berbeda sekitar 1 sampai 2 menit dari beberapa aplikasi lain, terutama untuk Maghrib dan Terbit, karena perbedaan lokasi/koordinat, metode, dan konvensi koreksi. Untuk penggunaan resmi, cocokkan jadwal dengan sumber resmi [Bimas Islam Kementerian Agama RI](https://bimasislam.kemenag.go.id/jadwalshalat).

## Screenshot

<!-- Screenshot menyusul; file docs/screenshot.png belum tersedia. -->

![Screenshot Prayer Times TUI](docs/screenshot.png)

## Pengembangan

```bash
go test ./...
go vet ./...
go build ./...
```

Core kalkulasi murni berada di `internal/prayer`, struktur/default konfigurasi di `internal/config`, dan label Bahasa Indonesia di `internal/i18n`. Entry point melakukan wiring input/output; aplikasi tidak memakai timezone mesin untuk menentukan tanggal/waktu lokasi.

## Roadmap

1. **Fase 1**: Core worldwide, ihtiyati Kemenag, metode, dan validasi akurasi.
2. **Fase 2**: Bubble Tea untuk navigasi tanggal, kembali ke hari ini (`t`), dan keluar dengan `q`/Ctrl+C.
3. **Fase 3**: Lipgloss, header, box waktu shalat, highlight dan countdown.
4. **Fase 4**: Konfigurasi user dengan Viper.
5. **Fase 5**: Packaging dan release AUR.

## Dokumentasi

- [agent.md](./agent.md) — konteks arsitektur dan prinsip development.
- [langkah.md](./langkah.md) — checklist fase, keputusan dan verifikasi.

## License

MIT
