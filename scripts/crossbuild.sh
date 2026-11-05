#!/bin/sh
# Builds Linkz for every system we ship to, with `lizardscript build`: one
# executable per system, the lizardscript runtime with Linkz's files (main.liz,
# every module it uses, the data files of bundle.txt) appended, on the VM.
# Each archive is named like the install scripts expect, with checksums.txt.
#   scripts/crossbuild.sh [out_dir]         (default: dist)
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH). The
# other systems' runtimes are built from lizard-script's source with Go: put
# it next to linkz (../lizard-script) or say where with LIZARD_SRC.
# The icon (assets/icon, from docs/Art/Figma-gen/r006c.py) goes in with
# --icon: macOS archives hold Linkz.app (signed ad hoc when codesign is there,
# not notarized), Windows gets it in the .exe, Linux gets icon.png and a
# .desktop template beside the binary.
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
OUT="${1:-dist}"
if [ -z "${LIZARD_SRC:-}" ] && [ -d ../lizard-script ]; then
  LIZARD_SRC="$(cd ../lizard-script && pwd)"
  export LIZARD_SRC
fi
VERSION=$(sed -n 's/^let linkz_version = "\(.*\)"$/\1/p' lib/app/version.liz)
[ -n "$VERSION" ] || { echo "no linkz_version in lib/app/version.liz"; exit 1; }
mkdir -p "$OUT"
for target in darwin/arm64 darwin/amd64 linux/amd64 windows/amd64; do
  os="${target%/*}"; arch="${target#*/}"
  name="linkz_${VERSION}_${os}_${arch}"
  dir="$OUT/$name"
  mkdir -p "$dir"
  bin="linkz"
  [ "$os" = windows ] && bin="linkz.exe"
  [ "$os" = darwin ] && bin="Linkz.app"
  echo "== $os/$arch"
  "$LZ" build --vm --os "$os" --arch "$arch" --icon assets/icon/linkz-1024.png \
    --icon-small assets/icon/linkz-16.png . -o "$dir/$bin" > /dev/null
  cp README.md LICENSE "$dir/"
  mkdir -p "$dir/examples" && cp examples/*.lzdraw "$dir/examples/"
  if [ "$os" = windows ]; then
    (cd "$dir" && zip -qr "../$name.zip" .)
  else
    COPYFILE_DISABLE=1 tar -czf "$OUT/$name.tar.gz" -C "$dir" .
  fi
  rm -rf "$dir"
done
(cd "$OUT" && shasum -a 256 linkz_* > checksums.txt)
echo "crossbuild ok: $OUT"
