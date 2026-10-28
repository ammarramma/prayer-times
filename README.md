# Prayer Times TUI 🕌

Aplikasi waktu shalat berbasis Go yang mendukung kota di seluruh dunia melalui koordinat dan timezone IANA. Kemenag RI tetap menjadi metode default, dengan koreksi ihtiyati untuk mempertahankan konvensi waktu Indonesia.

## Status

Fase 1 (core worldwide dan kalkulasi) sedang direvisi; fase interaktivitas TUI menyusul. Lihat [langkah.md](./langkah.md) untuk status detail.

## Fitur kalkulasi

- Lokasi dengan nama, lintang, bujur, timezone IANA. Default: Lhokseumawe (5.1801, 97.1406, `Asia/Jakarta`).
- Metode kalkulasi: Kemenag RI (Fajr 20°, Isha 18°), MWL, Umm al-Qura, Egyptian, Karachi, ISNA, Singapore, Dubai, Moonsighting Committee, Kuwait, Qatar, dan UOIF.
- Madhab Asr: Shafii (default) atau Hanafi.
- Aturan lintang tinggi: Middle of the Night, Seventh of the Night, Twilight Angle.
- Ihtiyati khusus Kemenag: default +2 menit untuk Subuh, Dzuhur, Ashar, Maghrib, Isya, dan −2 menit untuk Terbit. Untuk metode lain tidak ada koreksi ihtiyati dari aplikasi.
- Imsak opsional dihitung 10 menit sebelum Subuh.
- Waktu yang memiliki detik dibulatkan ke atas ke menit berikutnya.

## Menjalankan

```bash
go run .
```

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
| `-name` | `Lhokseumawe` | Nama kota/lokasi |
| `-lat` | `5.1801` | Lintang desimal, −90 sampai 90 |
| `-lon` | `97.1406` | Bujur desimal, −180 sampai 180 |
| `-tz` | `Asia/Jakarta` | Zona IANA (mis. `Europe/London`, `Asia/Riyadh`) |
| `-method` | `kemenag` | Salah satu metode di atas |
| `-madhab` | `shafii` | `shafii` atau `hanafi` |
| `-high-latitude` | `middle_of_the_night` | `middle_of_the_night`, `seventh_of_the_night`, `twilight_angle` |
| `-ihtiyati` | `2` | Menit koreksi khusus Kemenag; nilai metode lain diabaikan |
| `-imsak` | `false` | Tampilkan Imsak = Subuh −10 menit |
| `-date` | Hari ini di zona lokasi | Tanggal Gregorian `YYYY-MM-DD` |

Contoh metode lainnya: `umm_al_qura`, `egyptian`, `karachi`, `isna`, `singapore`, `dubai`, `moonsighting_committee`, `kuwait`, `qatar`, `uoif`.

## Pengembangan

```bash
go test ./...
go vet ./...
go build ./...
```

Core kalkulasi murni berada di `internal/prayer`, struktur/default konfigurasi di `internal/config`, dan label Bahasa Indonesia di `internal/i18n`. Entry point melakukan wiring input/output; aplikasi tidak memakai timezone mesin untuk menentukan tanggal/waktu lokasi.

## Roadmap

1. **Fase 1**: Core worldwide, ihtiyati Kemenag, metode, dan validasi akurasi.
2. **Fase 2**: Bubble Tea untuk navigasi tanggal kiri/kanan dan keluar dengan `q`/Ctrl+C.
3. **Fase 3**: Lipgloss, header, box waktu shalat, highlight dan countdown.
4. **Fase 4**: Konfigurasi user dengan Viper.
5. **Fase 5**: Packaging dan release AUR.

## Dokumentasi

- [agent.md](./agent.md) — konteks arsitektur dan prinsip development.
- [langkah.md](./langkah.md) — checklist fase, keputusan dan verifikasi.

## License

MIT
