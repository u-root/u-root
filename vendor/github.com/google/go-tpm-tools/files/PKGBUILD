# Maintainer: Joe Richey <joerichey@google.com>
pkgname=gotpm
pkgver=0.1.2
pkgrel=1
pkgdesc='TPM2 command-line utility'
arch=('x86_64')
_reponame=go-tpm-tools
url="https://github.com/google/${_reponame}"
license=('APACHE')
depends=('glibc') # go-pie requires CGO, so we have to link against libc
makedepends=('go-pie')
source=("git+${url}.git#tag=v${pkgver}?signed")
validpgpkeys=('19CE40CEB581BCD81E1FB2371DD6D05AA306C53F')
sha256sums=('SKIP')

build() {
  cd ${_reponame}
  go build \
    -trimpath \
    -ldflags "-extldflags $LDFLAGS" \
    ./cmd/${pkgname}
}

package() {
  cd ${_reponame}

  install -Dm755 $pkgname "${pkgdir}/usr/bin/${pkgname}"
  install -Dm755 files/boot-unseal.sh "${pkgdir}/etc/${pkgname}/boot-unseal.sh"

  initcpio_name='encrypt-gotpm'
  install -Dm644 files/initcpio.hooks "${pkgdir}/usr/lib/initcpio/hooks/${initcpio_name}"
  install -Dm644 files/initcpio.install "${pkgdir}/usr/lib/initcpio/install/${initcpio_name}"

  install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}
