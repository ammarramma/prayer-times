pkgname=prayer-times-tui
pkgver=0.1.1
pkgrel=1
pkgdesc='Prayer times for cities worldwide in a terminal UI'
arch=('x86_64' 'aarch64')
url='https://github.com/ammarramma/prayer-times'
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/ammarramma/prayer-times/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('dd08c9cbc03a5d8a95d4a8738e1e6af48f472d858fc6e236b8866cd17a4c8939')

prepare() {
  cd "$srcdir/prayer-times-$pkgver"
  GOFLAGS=-modcacherw go mod download -x
}

build() {
  cd "$srcdir/prayer-times-$pkgver"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -ldflags "-linkmode external -extldflags \"$LDFLAGS\" -X main.version=$pkgver" -o prayer-times .
}

check() {
  cd "$srcdir/prayer-times-$pkgver"
  go test ./...
}

package() {
  cd "$srcdir/prayer-times-$pkgver"
  install -Dm755 prayer-times "$pkgdir/usr/bin/prayer-times"
  install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
