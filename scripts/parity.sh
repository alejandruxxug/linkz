#!/bin/sh
# Parity of Linkz with the last release of the Go host (B042), whose outputs
# are kept in scripts/parity/golden: the SVG of every example (exact), paint
# of the graphics scenes (pixels within 1, the soft edges of the two pixel
# engines), trace of a picture that needs no shrinking (exact code) and an
# MCP session (the replies, without timings, step counts and pictures).
# Exits non-zero on the first difference.
#   scripts/parity.sh
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
BASE="${TMPDIR:-/tmp}"
TMP="$(mktemp -d "${BASE%/}/linkz-parity.XXXXXX")"
export HOME="$TMP/home"
mkdir -p "$HOME"
GOLD=scripts/parity/golden
port() { "$LZ" run --vm . "$@"; }

echo "render:"
for f in examples/*.lzdraw; do
  b=$(basename "$f" .lzdraw)
  port render "$f" "$TMP/$b.svg" > /dev/null
  cmp "$GOLD/render/$b.svg" "$TMP/$b.svg"
  echo "  $b same"
done

echo "paint:"
for s in rect:120x80 circle:120x120 star:200x100 stroke:200x100 brush:200x100 text:200x100 gradient:220x120 relgradient:220x140 blend:200x40 tiles:60x40; do
  name=${s%%:*}; size=${s#*:}
  port paint "lib/gfx/tests/$name.liz" "$TMP/$name.png" "$size" > "$TMP/$name.out"
  cmp "$GOLD/paint/$name.out" "$TMP/$name.out"
  printf '  %s: ' "$name"
  "$LZ" tools/pngdiff.liz "$GOLD/paint/$name.png" "$TMP/$name.png" 1
done

echo "trace:"
for t in "default:" "colors2:--colors 2" "colors4:--colors 4 --detail 2.5 --min-area 2"; do
  name=${t%%:*}; opts=${t#*:}
  # shellcheck disable=SC2086
  port trace "$GOLD/trace/picture.png" "$TMP/$name.lzdraw" $opts > /dev/null
  cmp "$GOLD/trace/$name.lzdraw" "$TMP/$name.lzdraw"
  echo "  [$opts] same"
done

echo "mcp:"
port mcp < scripts/parity/mcp.jsonl 2> /dev/null | sed "s#$HOME#\$HOME#g" > "$TMP/mcp.jsonl"
printf '  '
"$LZ" tools/mcpdiff.liz "$GOLD/mcp.jsonl" "$TMP/mcp.jsonl"
echo "parity ok"
