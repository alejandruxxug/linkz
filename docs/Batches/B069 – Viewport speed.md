# B069 – Viewport speed (quick wins, all lizard)

> "Can we optimize the viewport rendering?" … "A now !" — the developer, 2026-11-04

Measured before (VM, `tests/editor/editor_bench.liz`): keystroke 2k 14 ms (canvas 6–7), 10k 33 ms (canvas 23); first frame after a zoom change ~250 ms at 2k/400%; pan at 400% 45 ms at 10k; Alt-isolation slow at 1600%.

**Part 1 (agent 1)**
- **A5 – profile first:** the VM profiler on 10k keystroke and on the first frame at a new zoom; write the top costs into the report before changing anything.
- **A1 – never freeze on zoom:** on a zoom change, show the current tiles scaled at once (scaled blit), then re-render sharp tiles progressively within ~8 ms per frame, nearest the pointer first; panning keeps working meanwhile; an edit during refinement still shows correctly. Export/PNG and parity are untouched (they don't use tiles).

**Part 2 (agent 2, after part 1)**
- **A2 – zoom-aware flattening:** curve and circle flattening tolerance follows the screen scale (sharp at 3200%, fewer segments at 25%). Export keeps today's tolerance (parity byte-exact).
- **A3 – skip what can't be seen:** shapes whose on-screen box is under ~0.5 px are skipped (or one pixel) at that zoom; Alt-isolation rasterizes only the visible part.
- **A4 – a cache per zoom level:** spans of the last 2–3 zoom levels kept (memory cap, LRU), so zooming back is instant.

**Rules:** `_v11-rules.md`. Bench before/after at fit, 400%, 1600%, 2k and 10k, plus a new case "first frame after zoom" and "frames until sharp". Fit numbers must not regress > 5%. `scripts/test.sh`, `parity.sh`, `smoke.sh` green. Rig tests: zoom shows a placeholder then becomes exact (compare to a fresh render), edits during refinement, cache hits.
