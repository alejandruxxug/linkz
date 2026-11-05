#!/bin/sh
# Times the line-art tracer stage by stage without tasks and with them, and
# checks that both give the same code. See tools/trace_speed.liz.
#   scripts/trace_speed.sh [images...]
# LIZARDSCRIPT names the lizardscript to use (default: the one on PATH).
set -eu
cd "$(dirname "$0")/.."
LZ="${LIZARDSCRIPT:-lizardscript}"
exec "$LZ" --vm tools/trace_speed.liz "$@"
