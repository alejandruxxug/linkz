#!/bin/sh
# Parity of the lizard port with the Go host, on the same files: the SVG of
# every example (exact), paint of the graphics scenes (pixels within 1, the
# soft edges of the two pixel engines), trace of a picture that needs no
# shrinking (exact code) and an MCP session (the replies, without timings,
# step counts and pictures). Exits non-zero on the first difference.
#   scripts/parity.sh
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
TMP="$(mktemp -d "${TMPDIR:-/tmp}/linkz-parity.XXXXXX")"
export HOME="$TMP/home"
mkdir -p "$HOME"
GO="$TMP/linkz-go"
go build -o "$GO" .
port() { "$LZ" run --vm . "$@"; }

echo "render:"
for f in examples/*.lzdraw; do
  b=$(basename "$f" .lzdraw)
  "$GO" render "$f" "$TMP/$b-go.svg" > /dev/null
  port render "$f" "$TMP/$b.svg" > /dev/null
  cmp "$TMP/$b-go.svg" "$TMP/$b.svg"
  echo "  $b same"
done

echo "paint:"
for s in rect:120x80 circle:120x120 star:200x100 stroke:200x100 brush:200x100 text:200x100 gradient:220x120 relgradient:220x140 blend:200x40 tiles:60x40; do
  name=${s%%:*}; size=${s#*:}
  "$GO" paint "lib/gfx/tests/$name.liz" "$TMP/$name-go.png" "$size" > "$TMP/$name-go.out"
  port paint "lib/gfx/tests/$name.liz" "$TMP/$name.png" "$size" > "$TMP/$name.out"
  cmp "$TMP/$name-go.out" "$TMP/$name.out"
  printf '  %s: ' "$name"
  "$LZ" tools/pngdiff.liz "$TMP/$name-go.png" "$TMP/$name.png" 1
done

echo "trace:"
"$GO" paint lib/gfx/tests/star.liz "$TMP/picture.png" 200x150 > /dev/null
for opts in "" "--colors 2" "--colors 4 --detail 2.5 --min-area 2"; do
  # shellcheck disable=SC2086
  "$GO" trace "$TMP/picture.png" "$TMP/t-go.lzdraw" $opts > /dev/null
  # shellcheck disable=SC2086
  port trace "$TMP/picture.png" "$TMP/t.lzdraw" $opts > /dev/null
  cmp "$TMP/t-go.lzdraw" "$TMP/t.lzdraw"
  echo "  [$opts] same"
done

echo "mcp:"
"$GO" mcp < scripts/parity/mcp.jsonl > "$TMP/mcp-go.jsonl" 2> /dev/null
rm -rf "$HOME/Documents"
port mcp < scripts/parity/mcp.jsonl > "$TMP/mcp.jsonl" 2> /dev/null
printf '  '
"$LZ" tools/mcpdiff.liz "$TMP/mcp-go.jsonl" "$TMP/mcp.jsonl"
echo "parity ok"
