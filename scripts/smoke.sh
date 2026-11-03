#!/bin/sh
# Smoke test of the native editor: start it, drive it through the backdoor,
# look at a screenshot, quit. Exits non-zero on the first failure.
#   scripts/smoke.sh [screenshot.png]
set -eu
cd "$(dirname "$0")/.."
TMP="${TMPDIR:-/tmp}"; TMP="${TMP%/}"
SHOT="${1:-$TMP/linkz-smoke.png}"
BIN="$TMP/linkz-smoke-bin"
go build -o "$BIN" . 

# Our own window only: the debug file is per pid, and every command below
# names it, so a window the developer has open is never touched.
"$BIN" edit --no-splash &
PID=$!
DBG="$HOME/.linkz/debug-$PID.json"
cleanup() { kill "$PID" 2>/dev/null || true; rm -f "$DBG"; }
trap cleanup EXIT

# wait for the backdoor
i=0
while [ ! -f "$DBG" ]; do
  i=$((i + 1)); [ "$i" -gt 100 ] && { echo "editor did not start"; exit 1; }
  sleep 0.1
done

d() { "$BIN" debug --pid "$PID" "$@"; echo; }
d Wait
has() { case "$1" in *"$2"*) ;; *) echo "expected $2 in: $1"; exit 1;; esac; }
has "$(d Call Ping hello)" '"pong":true'
d Drag 100 160 300 280 10 --pressure 0.2,1.0
has "$(d Screenshot "$SHOT")" "$SHOT"
has "$(d Dump)" '"held":false'
d Stats
d Quit
wait "$PID" || true
trap - EXIT
rm -f "$DBG"
echo "smoke ok: $SHOT"
