#!/bin/sh
# The lizard tests of Linkz (tests/), on both engines: everything on the
# tree-walking evaluator, then everything on the bytecode VM. HOME is a folder
# of its own, so the tests that save files never touch the user's.
#   scripts/test.sh [--bench]
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
TMP="$(mktemp -d "${TMPDIR:-/tmp}/linkz-test.XXXXXX")"
mkdir -p "$TMP/tree" "$TMP/vm"
echo "== evaluator"
HOME="$TMP/tree" LINKZ_ENGINE=tree "$LZ" test "$@" tests
echo "== vm"
HOME="$TMP/vm" "$LZ" test --vm "$@" tests
