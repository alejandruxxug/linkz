#!/bin/sh
# Scores the tracer on the reference pictures: error, ink IoU, shapes, size
# and time per picture (the gate of Tracing v3). See tools/trace_score.liz
# for the options (--mode, --json, --out, --md, --compare) and the ink rule.
#   scripts/trace_score.sh [--mode lineart|photo] [images...]
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
exec "$LZ" --vm tools/trace_score.liz "$@"
