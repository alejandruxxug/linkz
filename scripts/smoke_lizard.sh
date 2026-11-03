#!/bin/sh
# Smoke test of the lizard port: render, paint and trace from the command
# line, then the native editor started and driven with `lizardscript drive`,
# a screenshot looked at, quit. Exits non-zero on the first failure. HOME is a
# folder of its own, so nothing touches the user's drawings or settings.
#   scripts/smoke_lizard.sh [screenshot.png]
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
TMP="$(mktemp -d "${TMPDIR:-/tmp}/linkz-smoke.XXXXXX")"
SHOT="${1:-$TMP/linkz-smoke.png}"
export HOME="$TMP/home"
mkdir -p "$HOME"

has() { case "$1" in *"$2"*) ;; *) echo "expected $2 in: $1"; exit 1;; esac; }

# The commands.
has "$("$LZ" run --vm . render examples/hello.lzdraw "$TMP/hello.svg")" "wrote $TMP/hello.svg"
has "$(cat "$TMP/hello.svg")" "<svg"
"$LZ" run --vm . paint lib/gfx/tests/star.liz "$TMP/star.png" 200x100
"$LZ" run --vm . trace "$TMP/star.png" "$TMP/star.lzdraw" --colors 2 > /dev/null
has "$(cat "$TMP/star.lzdraw")" "path("
has "$(printf '%s\n' '{"jsonrpc":"2.0","id":1,"method":"ping"}' | "$LZ" run --vm . mcp 2> /dev/null)" '"id":1'

# The editor, in our own window only: every drive command names its pid, so
# a window the developer has open is never touched.
LIZARD_AUTOMATION=1 "$LZ" run --vm . edit --no-splash &
PID=$!
cleanup() { kill "$PID" 2>/dev/null || true; }
trap cleanup EXIT

d() { "$LZ" drive --pid "$PID" "$@"; echo; }
i=0
until "$LZ" drive --pid "$PID" Wait > /dev/null 2>&1; do
  i=$((i + 1)); [ "$i" -gt 100 ] && { echo "editor did not start"; exit 1; }
  sleep 0.1
done
has "$(d Call Ping hello)" '"pong":true'
has "$(d Dump)" '"screen":"home"'
d Call Show editor
d Drag 100 160 300 280 10 --pressure 0.2,1.0
has "$(d Screenshot "$SHOT")" "$SHOT"
has "$(d Dump)" '"screen":"editor"'
d Stats
d Quit
wait "$PID" || true
trap - EXIT
echo "smoke ok: $SHOT"
