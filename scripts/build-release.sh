#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version> (example: v0.1.0)" >&2
  exit 2
fi

version="$1"
if [[ ! "$version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
  echo "Invalid version: $version" >&2
  exit 2
fi

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
dist_dir="$root_dir/dist"
mkdir -p "$dist_dir"
rm -f "$dist_dir"/prayer-times_"${version}"_linux_*.tar.gz "$dist_dir/SHA256SUMS"

for arch in amd64 arm64; do
  archive="$dist_dir/prayer-times_${version}_linux_${arch}.tar.gz"
  temp_dir="$(mktemp -d)"
  trap 'rm -rf "$temp_dir"' EXIT

  (
    cd "$root_dir"
    CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build \
      -trimpath \
      -ldflags "-s -w -X main.version=$version" \
      -o "$temp_dir/prayer-times" .
  )
  install -m 644 "$root_dir/LICENSE" "$temp_dir/LICENSE"
  install -m 644 "$root_dir/README.md" "$temp_dir/README.md"
  tar -C "$temp_dir" -czf "$archive" prayer-times LICENSE README.md
  rm -rf "$temp_dir"
  trap - EXIT
done

(
  cd "$dist_dir"
  sha256sum prayer-times_"${version}"_linux_*.tar.gz > SHA256SUMS
)
printf 'Release artifacts created in %s\n' "$dist_dir"
