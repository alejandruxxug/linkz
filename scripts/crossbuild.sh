#!/bin/sh
# Build and vet Linkz for every system we ship to. Windows and Linux can't be
# run from here, so this is what keeps them compiling.
#   scripts/crossbuild.sh
set -eu
cd "$(dirname "$0")/.."
for target in windows/amd64 linux/amd64 darwin/arm64; do
  os="${target%/*}"; arch="${target#*/}"
  echo "== $os/$arch"
  GOOS="$os" GOARCH="$arch" go build ./...
  GOOS="$os" GOARCH="$arch" go vet ./...
done
echo "crossbuild ok"
