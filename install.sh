#!/bin/sh
set -e

RELEASES_URL="https://github.com/What-If-I/tlog/releases"
FILE_BASENAME="tlog"

test -z "$VERSION" && VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" |
  rev |
  cut -f1 -d'/' |
  rev)"

test -z "$VERSION" && {
  echo "Unable to get tlog version." >&2
  exit 1
}

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="$TMPDIR/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"

(
  cd "$TMPDIR"
  echo "Downloading tlog $VERSION..."
  curl -v -fLo "$TAR_FILE" \
    "$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"
  curl -sfLo "checksums.txt" "$RELEASES_URL/download/$VERSION/checksums.txt"
  curl -sfLo "checksums.txt.sig" "$RELEASES_URL/download/$VERSION/checksums.txt.sig"
  echo "Verifying checksums..."
  sha256sum --ignore-missing --quiet --check checksums.txt
  if command -v cosign >/dev/null 2>&1; then
    echo "Verifying signatures..."
    COSIGN_EXPERIMENTAL=1 cosign verify-blob \
      --signature checksums.txt.sig \
      checksums.txt
  else
    echo "Could not verify signatures, cosign is not installed."
  fi
)

tar -xf "$TAR_FILE" -C "$TMPDIR"
mv "${TMPDIR}/tlog" "/usr/local/bin/"
