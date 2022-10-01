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

V_NUMBER=$(echo $VERSION | cut -c2-)
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="$TMPDIR/${FILE_BASENAME}_${V_NUMBER}_$(uname -s)_$(uname -m).tar.gz"

(
  cd "$TMPDIR"
  echo "Downloading tlog $VERSION..."
  curl -sfLo "$TAR_FILE" \
    "$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_${V_NUMBER}_$(uname -s)_$(uname -m).tar.gz"
)

tar -xf "$TAR_FILE" -C "$TMPDIR"
EXEC_PATH="/usr/local/bin/"
mv "${TMPDIR}/tlog" "/usr/local/bin/"
echo "Saved to" $EXEC_PATH
