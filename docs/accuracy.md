# Calculation Accuracy

## Scope and verification

Prayer times are calculated for Gregorian calendar dates using the location's IANA time zone. The project rounds times upward to the next whole minute after calculation and configured adjustments. No prayer-time profiles other than Kemenag have an application-owned reference fixture. For that reason, every other method is marked **unverified** here even where its parameters match the `adhango v0.1.0` built-in preset.

Kemenag uses Fajr 20° and Isha 18°, with the existing application safety margin (default +2 minutes; Dhuhr +3 minutes and Sunrise −2 minutes). Three Lhokseumawe dates are compared to an October 2026 schedule listing Bimas Islam/Kemenag as its source. The source is an aggregator page, not a directly captured official Bimas Islam table; see [langkah.md](../langkah.md). Only values that match exactly are called verified at that tolerance. Ashar on 2026-10-16 and Maghrib on 2026-10-28 remain within ±1 minute and are not claimed as exact.

The “checked” date below is the date the local preset/profile and tests were reviewed (2026-10-29), not a claim of confirmation by a method authority. “Fixture count” counts schedule dates, not individual prayer entries.

## Profiles

| Method | Fajr angle | Isha angle / interval | Built-in offsets (Fajr, Sunrise, Dhuhr, Asr, Maghrib, Isha) | Verification | Source | Checked | Fixture dates |
| --- | ---: | --- | --- | --- | --- | --- | ---: |
| Kemenag | 20° | 18° | +2, −2, +3, +2, +2, +2 min (includes default ihtiyati) | Partial | Project convention and three Lhokseumawe schedule dates; official Bimas Islam table not independently verified | 2026-10-29 | 3 |
| MWL | 18° | 17° | 0, 0, +1, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Umm al-Qura | 18.5° | 90 min | 0, 0, 0, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Egyptian | 19.5° | 17.5° | 0, 0, +1, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Karachi | 18° | 18° | 0, 0, +1, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| ISNA | 15° | 15° | 0, 0, +1, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Singapore | 20° | 18° | 0, 0, +1, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Dubai | 18.2° | 18.2° | 0, −3, +3, +3, +3, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Moonsighting Committee | 18° | 18° | 0, 0, +5, 0, +3, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Kuwait | 18° | 17.5° | 0, 0, 0, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| Qatar | 18° | 90 min | 0, 0, 0, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |
| UOIF | 12° | 12° | 0, 0, 0, 0, 0, 0 min | Unverified | `adhango v0.1.0` built-in method preset | 2026-10-29 | 0 |

The Kemenag profile's offset tuple is the normal default ihtiyati result. Setting `ihtiyati_minutes` to zero removes those base adjustments; the independent user `offsets` values are then applied afterward. User offsets are also added after the method's built-in offsets for all other methods. Each configured offset is constrained to −30 through +30 minutes. Imsak remains Fajr minus 10 minutes.

## Regression and date coverage

`internal/prayer/testdata/golden.json` stores minute-level output from the pre-profile implementation for all 12 methods, six cities (Lhokseumawe, Jakarta, London, Makkah, Tokyo, Sydney), and two dates (2026-10-28 and 2026-06-21). These snapshots are regression checks, not accuracy fixtures: they preserve prior behavior and do not establish agreement with an authority.

Additional tests cover the London and Sydney DST transitions and calculation robustness for Tromsø on 2026-06-21 and 2026-12-21 across all methods. These tests verify time-zone transitions, ordering, valid results or returned errors, not local authority accuracy.

## Known differences and open data request

Unofficial observations for Lhokseumawe on 2026-10-29 (not test fixtures): an Android application reported Imsak 04:50, Fajr 05:00, Sunrise 06:17, Dhuha 06:41, Dhuhr 12:19, Asr 15:38, Maghrib 18:17, and Isha 19:27. Google with “Kemenag Jakarta Pusat” as its basis and no ihtiyati reported Fajr 04:58, Sunrise 06:17, Dhuhr 12:15, Asr 15:36, Maghrib 18:14, and Isha 19:24. This project produced Sunrise 06:15 and Maghrib 18:16. These are observations only; no source is asserted for them and they are not fixtures.

To resolve the remaining Sunrise/Maghrib differences, the project needs an official Bimas Islam Lhokseumawe schedule table for at least five dates in October 2026, including sunrise and Maghrib. Potential explanations are hypotheses only until such data is available. The Maghrib and Sunrise formulas have not been adjusted to fit these observations.
