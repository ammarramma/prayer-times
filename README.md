# Prayer Times TUI

A worldwide prayer schedule for the terminal, written in Go with Bubble Tea and Lipgloss. Enter a location's name, coordinates, and IANA time zone to calculate schedules with one of 12 calculation methods. The interface and CLI output are English by default; Indonesian is available with `-lang id`.

The default location is Lhokseumawe, Aceh (`5.1801, 97.1406`, `Asia/Jakarta`). When no method is explicitly configured, Indonesian time zones use Kemenag and other locations use MWL. An existing config method or `-method` flag takes precedence. When no high-latitude rule is explicitly configured, locations beyond 48° latitude use seventh of the night and all others use middle of the night.

<!-- Screenshot: docs/screenshot.png is coming soon; no image is embedded yet. -->

## Install

Install with Go:

```bash
go install github.com/ammarramma/prayer-times@latest
```

Prebuilt Linux binaries for amd64 and arm64 are attached to every [GitHub release](https://github.com/ammarramma/prayer-times/releases) together with a `SHA256SUMS` file. Download the archive for your architecture, verify it, and extract it:

```bash
curl -LO https://github.com/ammarramma/prayer-times/releases/latest/download/prayer-times_v0.1.0_linux_amd64.tar.gz
curl -LO https://github.com/ammarramma/prayer-times/releases/latest/download/SHA256SUMS
sha256sum --check --ignore-missing SHA256SUMS
tar -xzf prayer-times_v0.1.0_linux_amd64.tar.gz
./prayer-times -version
```

On Arch Linux, clone the repository and build the PKGBUILD it ships with:

```bash
git clone https://github.com/ammarramma/prayer-times.git
cd prayer-times
makepkg -si
```

An AUR package is coming soon. Until it is published, the PKGBUILD in the repository is the supported Arch install path.

You can also run it from a source checkout with `go run .`.

## Use

Run the interactive schedule:

```bash
prayer-times
```

Navigate days with ←/→, press `t` to return to today at the selected location, and press `q` or Ctrl+C to quit. The TUI runs inline by default; `-fullscreen` switches to the alternate screen.

Print a one-time schedule with `-date`:

```bash
prayer-times \
  -name London -lat 51.5072 -lon -0.1276 -tz Europe/London \
  -method mwl -madhab hanafi -high-latitude seventh_of_the_night \
  -date 2026-10-29
```

Use `-lang id` for Indonesian CLI and TUI labels. The default language is English. A `language` entry in the config is also supported.

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-config` | OS config directory + `prayer-times/config.json` | JSON config path |
| `-name` | `Lhokseumawe` | Location display name |
| `-lat` | `5.1801` | Latitude, from −90 to 90 |
| `-lon` | `97.1406` | Longitude, from −180 to 180 |
| `-tz` | `Asia/Jakarta` | IANA time zone, such as `Europe/London` or `Asia/Riyadh` |
| `-method` | Kemenag for Indonesia time zones; MWL otherwise | One of the 12 methods listed below |
| `-madhab` | `shafii` | Asr school: `shafii` or `hanafi` |
| `-high-latitude` | `seventh_of_the_night` beyond 48° latitude, `middle_of_the_night` otherwise | `middle_of_the_night`, `seventh_of_the_night`, or `twilight_angle` |
| `-ihtiyati` | `2` for Kemenag; `0` otherwise | Kemenag safety margin in minutes; Dhuhr uses one additional minute and Sunrise subtracts the base margin |
| `-imsak` | `false` | Include Imsak, calculated as Fajr minus 10 minutes |
| `-fullscreen` | `false` | Use the alternate screen and center the TUI |
| `-date` | Today at the selected location | Print one Gregorian `YYYY-MM-DD` schedule and exit |
| `-set-config` | `false` | Save the resolved config and exit; the selected method is saved explicitly |
| `-lang` | `en` | Language: `en` or `id` |
| `-version` | `dev` | Print the application version |

Methods: `kemenag`, `mwl`, `umm_al_qura`, `egyptian`, `karachi`, `isna`, `singapore`, `dubai`, `moonsighting_committee`, `kuwait`, `qatar`, and `uoif`.

**Flag parsing note:** flags must come before positional arguments. Positional arguments are rejected with an error; for example, a stray `.` after flags is not silently ignored.

## Configuration

The JSON config stores `location`, `options`, and `language`. Flags override matching file values. Older config files that omit `language` and `options.offsets` remain supported.

```json
{
  "location": {
    "name": "London",
    "latitude": 51.5072,
    "longitude": -0.1276,
    "timezone": "Europe/London"
  },
  "options": {
    "method": "mwl",
    "madhab": "shafii",
    "high_latitude_rule": "seventh_of_the_night",
    "ihtiyati_minutes": 0,
    "include_imsak": false,
    "offsets": {
      "fajr": 0,
      "sunrise": 0,
      "dhuhr": 0,
      "asr": 0,
      "maghrib": 0,
      "isha": 0
    }
  },
  "language": "en"
}
```

Each manual offset is in minutes and must be between −30 and +30. It is added after the selected method's built-in offsets and applies to Kemenag too. The legacy `ihtiyati_minutes` setting continues to control Kemenag's base safety margin.

The London example above stores `seventh_of_the_night` because that is the resolved rule for a location beyond 48° latitude. Omit `options.high_latitude_rule` to let the latitude decide again; `-set-config` writes the resolved rule explicitly.

Save the currently resolved settings (including the method selected from the time zone when not explicitly supplied) with:

```bash
prayer-times -config ./config.json -name London -lat 51.5072 -lon -0.1276 -tz Europe/London -set-config
```

## NOW and NEXT

On the displayed current date, a prayer is marked **NOW** for the first 15 minutes after its displayed time, including configured offsets and safety margins. During that window, the following prayer is marked **NEXT** in blue. At other times, only the next prayer is highlighted in green. Both labels include text, so they do not depend on color alone. After Isha, the next event is tomorrow's Fajr; before today's Fajr, it is today's Fajr. Selecting a different calendar date hides NOW, NEXT, and the countdown.

Sunrise is displayed as information and is not a prayer or a NOW/NEXT target.

## Accuracy

Prayer schedules can differ by 1–2 minutes from other applications because of coordinates, calculation conventions, built-in method adjustments, minute rounding, and local safety-margin practice. Check the schedule against the relevant local authority before relying on it for official use. For Indonesia, consult [Bimas Islam, Ministry of Religious Affairs](https://bimasislam.kemenag.go.id/jadwalshalat).

Only the Kemenag profile has schedule-date comparisons in this repository, and its verification is partial. Other method profiles are **unverified** against local authorities. The stored golden schedules are regression snapshots of the prior implementation, not accuracy fixtures. See [docs/accuracy.md](docs/accuracy.md) for profile parameters, fixture coverage, and unresolved differences.

## Development

```bash
go build ./...
go vet ./...
go test ./...
```

The calculation domain is in `internal/prayer`, configuration in `internal/config`, locale strings in `internal/i18n`, and the Bubble Tea UI in `internal/tui`.
