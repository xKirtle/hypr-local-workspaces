# Maintainer: Rodrigo Martins rodm.martins@proton.me
pkgname=hypr-local-workspaces
pkgver=0.1.0
pkgrel=1
pkgdesc="Make Hyprland workspaces local per monitor (Go)"
arch=('x86_64' 'aarch64')
url="https://github.com/xKirtle/hypr-local-workspaces"
license=('MIT')
depends=('glibc')
makedepends=('go')
options=('!debug' '!strip')
source=("$url/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('a193eb69633bccb75d93a1a28106c3eda8782fc1fcd16ecdcc03cabecd966fbf')

build() {
  cd "$srcdir/hypr-local-workspaces-$pkgver"
  export CGO_ENABLED=0
  go build -trimpath -ldflags "-s -w -X main.version=$pkgver" \
    -o hypr-local-workspaces .
}

package() {
  cd "$srcdir/hypr-local-workspaces-$pkgver"
  install -Dm755 "hypr-local-workspaces" "$pkgdir/usr/bin/hypr-local-workspaces"
  install -Dm644 "LICENSE" "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
