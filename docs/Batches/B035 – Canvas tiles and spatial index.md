# B035 – Canvas tiles, dirty rectangles and a spatial index (M10 Performance, Linkz side)

Repo: `linkz` · Lane: lizard (`lib/gfx`, `lib/ui/view_canvas.liz`, maybe `lib/draw`) + host primitives only if the profile proves they're needed · Needs: B029 (VM)

## Measure first
Profile the editor with 2,000 and 10,000 shapes (generate fixtures in `docs/Batches/fixtures/` or testdata): Go pprof (`go test -cpuprofile` on a host benchmark that loads a document, repaints, moves one shape, types one key) and per-function timing in lizard. Report where the 55 ms repaint and the 37 ms keystroke go **before** changing anything. (lizard-script is getting a VM profiler in B033; until it lands, use time.ms() around lizard phases.)

## Build
1. **Spatial index** in lizard (`lib/gfx/spatial.liz` or `lib/draw`): a uniform grid of buckets over the canvas; insert/remove/update a shape by bounds; query by rect and by point. Used for hit-testing (click, hover, lasso later) and tile invalidation.
2. **Tiles**: the canvas raster is kept as 256×256 tiles (host keeps the pixels as offscreen buffers; lizard decides what to draw in each). A change (edit of a row, drag of a shape or anchor, a new stroke) computes the **dirty rectangle** (old bounds ∪ new bounds of what changed, plus stroke width/AA margin), marks the tiles it touches, and only those are re-rasterized (only the shapes the index says touch each tile, in paint order). Untouched tiles are blitted from cache.
3. **Layer caches**: blend modes/opacity still composite correctly per tile (the per-layer buffers from B020 work per tile).
4. **Re-run avoidance**: when only one literal changed (two-way edit, drag), update that shape's data from the edit instead of re-running the whole program when it's provably safe (the shape's row has only literals); otherwise re-run, but diff the shape list (by line tag + value) so only changed shapes dirty tiles.
5. Zoom and pan: while a gesture is running, scale/translate the cached tiles; re-rasterize sharp tiles when it stops.

## Done when
Before/after table with fixtures of 2k and 10k shapes: repaint after a small edit, drag of one shape (per frame), full repaint, keystroke, memory. Target: small-edit repaint and drag frames < 16 ms at 2k shapes, comfortable at 10k; recommended object count re-measured (`linkz calibrate`). Pixel-identical output vs a full repaint (golden test). All tests and `scripts/smoke.sh` green on the VM, and with `LINKZ_ENGINE=tree`.

## Rules
Clean code (small, documented, no hacks). The Linkz rule holds: lizard decides pixels; a new Go primitive only for pure pixel moving (blit/scale of tile buffers), justified by the profile. Testers use `HOME=$(mktemp -d)`. Don't commit. linkz builds against `../.lizard-stable` (go.work); don't change it.
