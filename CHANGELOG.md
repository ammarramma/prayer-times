# Changelog

All notable changes to this project are documented in this file. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses [semantic versioning](https://semver.org/).

## [Unreleased]

## [0.1.1] - 2026-10-30

### Added

- Desktop notifications: `-notify-once`, `-dry-run`, `-notify-prayers`, `-notify-lead`, and `-notify-grace`, delivered through the freedesktop D-Bus notification interface.
- `notifications` config block with `enabled`, `prayers`, `lead_minutes`, and `grace_minutes`; delivery state is kept in `$XDG_STATE_HOME/prayer-times/notified.json`.
- systemd user units in `contrib/systemd/user/` for a once-a-minute notification check.
- Indonesian locale with `-lang id`; English remains the default.
- Calculation method profiles and per-prayer user offsets in the config.
- NOW and NEXT status on the displayed current date, with a 15-minute NOW window.
- Golden schedule snapshots covering 12 methods, six cities, and two dates, plus DST and high-latitude tests.
- GitHub Actions CI running build, vet, and tests on push and pull requests.
- Embedded IANA time zone database, so zone lookups work without system zoneinfo.

### Changed

- The high-latitude rule now defaults to `seventh_of_the_night` beyond 48° latitude and stays `middle_of_the_night` elsewhere; an explicit flag or config value still wins.
- `-version` falls back to the module version recorded by `go install` when no release version is injected.
- README installation section covers Go, release archives with `SHA256SUMS`, the Arch PKGBUILD, and the pending AUR package.
- Unrelated positional arguments are rejected with a clear CLI error.

### Fixed

- The `DefaultMethod` doc comment is no longer attached to the end of `DefaultOptions`.

## [0.1.0] - 2026-10-29

### Added

- Initial release: worldwide prayer schedules with 12 calculation methods, Kemenag ihtiyati, madhab choice, and optional Imsak.
- Bubble Tea terminal UI with day navigation, Lipgloss styling, next-prayer highlight, and countdown.
- Viper configuration with CLI flags, `-date` one-shot output, and `-set-config`.
- MIT license, Arch PKGBUILD, and `scripts/build-release.sh` for Linux amd64 and arm64 archives with `SHA256SUMS`.
