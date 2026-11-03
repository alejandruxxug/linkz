#!/usr/bin/env bash
# Install Linkz on macOS or Linux from a GitHub release of alejandruxxug/linkz.
#
#   scripts/install.sh                  # latest release -> ~/.local/bin
#   scripts/install.sh --system         # -> /usr/local/bin (uses sudo)
#   scripts/install.sh --version v0.1.0 # a specific release
#   scripts/install.sh --dir ~/bin      # another folder
#   scripts/install.sh --from-file linkz_0.1.0_darwin_arm64.tar.gz [--checksums checksums.txt]
#
# Downloads come only from github.com/alejandruxxug/linkz releases. While that
# repo is private this works once a release is public, or with a token that can
# read the repo:  GITHUB_TOKEN=... scripts/install.sh
#
# --from-file installs an archive you already have (no network), for testing
# and offline machines. The archive is checked when --checksums is given.
set -euo pipefail

REPO="alejandruxxug/linkz"
BIN="linkz"
version=""
dir="$HOME/.local/bin"
system=0
from_file=""
checksums_file=""

die() { echo "install: $*" >&2; exit 1; }

usage() { sed -n '2,13p' "$0" | sed 's/^# \{0,1\}//'; }

while [ $# -gt 0 ]; do
  case "$1" in
    --version)   [ $# -ge 2 ] || die "--version needs a value"; version="$2"; shift 2 ;;
    --dir)       [ $# -ge 2 ] || die "--dir needs a value"; dir="$2"; shift 2 ;;
    --system)    system=1; dir="/usr/local/bin"; shift ;;
    --from-file) [ $# -ge 2 ] || die "--from-file needs a value"; from_file="$2"; shift 2 ;;
    --checksums) [ $# -ge 2 ] || die "--checksums needs a value"; checksums_file="$2"; shift 2 ;;
    -h|--help)   usage; exit 0 ;;
    *) die "unknown option: $1 (try --help)" ;;
  esac
done

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

# Which archive is ours?
case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *) die "unsupported system $(uname -s); Windows users: scripts/install.ps1" ;;
esac
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64|amd64)  arch=amd64 ;;
  *) die "unsupported CPU $(uname -m)" ;;
esac
case "$os/$arch" in
  darwin/arm64|darwin/amd64|linux/amd64) ;;
  *) die "no Linkz build for $os/$arch" ;;
esac

sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1
  else shasum -a 256 "$1" | cut -d' ' -f1; fi
}

# check_archive ARCHIVE_PATH ARCHIVE_NAME CHECKSUMS_PATH
check_archive() {
  local want got
  want="$(awk -v n="$2" '$2 == n || $2 == "*"n { print $1 }' "$3")"
  [ -n "$want" ] || die "$2 is not listed in the checksums"
  got="$(sha256_of "$1")"
  [ "$want" = "$got" ] || die "checksum mismatch for $2 (wanted $want, got $got)"
  echo "checksum ok"
}

auth=()
[ -n "${GITHUB_TOKEN:-}" ] && auth=(-H "Authorization: Bearer $GITHUB_TOKEN")

if [ -n "$from_file" ]; then
  [ -f "$from_file" ] || die "no such file: $from_file"
  archive="$tmp/$(basename "$from_file")"
  cp "$from_file" "$archive"
  if [ -n "$checksums_file" ]; then
    check_archive "$archive" "$(basename "$from_file")" "$checksums_file"
  else
    echo "note: no --checksums given, the archive is not verified" >&2
  fi
else
  command -v curl >/dev/null 2>&1 || die "curl is needed"
  api="https://api.github.com/repos/$REPO/releases"
  if [ -n "$version" ]; then rel="$api/tags/$version"; else rel="$api/latest"; fi
  json="$tmp/release.json"
  curl -fsSL ${auth[@]+"${auth[@]}"} -H "Accept: application/vnd.github+json" "$rel" -o "$json.raw" ||
    die "could not read the release ($rel). The repo is private until a release is public; set GITHUB_TOKEN to a token that can read it."
  # One key per line, whether or not GitHub pretty-printed the JSON.
  tr ',{' '\n\n' < "$json.raw" > "$json"
  tag="$(sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' "$json" | head -n1)"
  [ -n "$tag" ] || die "the release has no tag"
  ver="${tag#v}"
  name="${BIN}_${ver}_${os}_${arch}.tar.gz"
  echo "installing $BIN $tag for $os/$arch"

  # asset_url NAME: the API url of a release asset (works for private repos).
  asset_url() {
    awk -v n="$1" '
      /"url": *"[^"]*\/releases\/assets\/[0-9]+"/ { gsub(/.*"url": *"|".*/, ""); u = $0 }
      /"name": *"/ { gsub(/.*"name": *"|".*/, ""); if ($0 == n) { print u; exit } }
    ' "$json"
  }
  fetch() { # NAME DEST
    local url
    url="$(asset_url "$1")"
    [ -n "$url" ] || die "release $tag has no asset $1"
    curl -fsSL ${auth[@]+"${auth[@]}"} -H "Accept: application/octet-stream" "$url" -o "$2"
  }
  archive="$tmp/$name"
  fetch "$name" "$archive"
  fetch "checksums.txt" "$tmp/checksums.txt"
  check_archive "$archive" "$name" "$tmp/checksums.txt"
fi

mkdir -p "$tmp/x"
tar -xzf "$archive" -C "$tmp/x"
[ -f "$tmp/x/$BIN" ] || die "the archive has no $BIN"
chmod +x "$tmp/x/$BIN"

if [ "$system" = 1 ]; then
  sudo mkdir -p "$dir"
  sudo install -m 755 "$tmp/x/$BIN" "$dir/$BIN"
else
  mkdir -p "$dir"
  install -m 755 "$tmp/x/$BIN" "$dir/$BIN"
fi
echo "installed $dir/$BIN"

echo
echo "Next steps"
case ":$PATH:" in
  *":$dir:"*) ;;
  *) echo "  - add $dir to your PATH:  export PATH=\"$dir:\$PATH\"" ;;
esac
echo "  - try it:  $BIN render examples/hello.lzdraw   (examples are in the archive)"
echo "  - open the editor:  $BIN edit"
if [ "$os" = darwin ]; then
  echo "  - if macOS says the app can't be opened (it is not notarized yet):"
  echo "      xattr -d com.apple.quarantine $dir/$BIN"
fi
