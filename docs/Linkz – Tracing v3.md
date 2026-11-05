# Linkz – Tracing v3: smaller, smoother, faster

> "Build everything and well look up anime girls references with somewhat clean colors and also messy colors!" — the developer, 2026-11-04

Builds on [[Linkz – Tracing v2]]. Where it stood (2026-11-04): ink IoU 0.70–0.86, mean color error 5–13 on the six Wikipe-tan references, a full line-art trace 3–4 s, ~220 KB of code per picture.

**Gate for every step:** `scripts/trace_score.sh` (step 0) on every reference before and after: mean color error, ink IoU, shapes, code bytes, ms. A step lands only if no reference gets clearly worse (error +0.3 or IoU −0.01 is the line; explain anything near it). Clean-color and messy-color references are both in the set.

| # | Step | Lane | Model |
|---|---|---|---|
| R | New references: anime-style girls, some clean flat colors, some messy/painted, freely licensed, with attribution | refs | Sonnet |
| 0 | `scripts/trace_score.sh`: one table for all references (and the synthetic set); the ink mask of the original is computed by a fixed, documented rule, not the frozen pre-B062 mask | fills | Sonnet |
| 1 | Merge pass: touching regions of one palette color joined; shapes nobody can see dropped; stack instead of holes (a region with islands is drawn whole, the islands on top); a same-color shape inside the one under it absorbed; near-equal palette colors (ΔE < 2) snapped to one entry | fills | Sonnet |
| 2 | Gradients: groups of neighbouring regions whose colors ramp smoothly → one shape with a `linear(...)`/`radial(...)` fill (least-squares fit), or stepped (`steps = n`) when the bands are few and even; hard cel shadows stay flat | fills | Sonnet |
| 3 | Strokes as curves: `brush` accepts Bézier items with a width/pressure per on-curve point; lineart strokes are fitted and written that way | strokes | Sonnet |
| 5 | Detail where it matters: smaller regions and a tighter fit where local detail is high (eyes, faces, hands, edges of high contrast), harder simplification in large flat areas | fills | Sonnet |
| 4 | Speed: per-stage timings; trace tiles in parallel with `task` and stitch them; then the slowest stage | after 1–3 | Sonnet |
| 6 | Eye check: before/after sheets of every reference in Figma for the developer | last | design |

## Step 4 result: speed (2026-11-04)

**Same drawing, sooner.** All 18 references trace to byte-identical code (digest checked against the code from before the step, and `scripts/trace_speed.sh` checks serial against parallel); `scripts/trace_score.sh --compare docs/Art/Checks/trace/v3/after-detail.json` shows no change in error, IoU or bytes. Total tracer time over the 18: 134.7 s before, 74.8 s after (`before-speed.*`, `after-speed.*` in `docs/Art/Checks/trace/v3/`; the machine was shared, so read the ratios, not the seconds).

- **Timings.** `stats.stages` has the milliseconds of every stage, with `mask_stages` (luma, paper, classify, flood) and `fill_stages` (solid, detail_map, regions, merge, gradients) inside; `stats.ms` is still the total. `scripts/trace_speed.sh [images]` prints them serial and parallel side by side.
- **Option `parallel`** (default true). The independent pieces run in `task`s (`lineart_parallel.liz`; the pieces are in `lineart_task.liz`): the box sums behind the mask's paper level and paper color, and the outline of each fill region (regions dealt out to 12 pieces, biggest first). Each piece is the same function the serial path calls, so the result cannot differ. Small pictures stay serial (a task costs about 15 ms and a copy of its data). With no `task` (inside a task, a sandbox that doesn't hand it over, `"parallel": false`) everything runs serially.
- **Why not bands.** The plan said horizontal bands with an overlap. The mask sums are float prefix sums (a band restarts the sum and rounds differently), the distance transform and the thinning follow raster order, and the palette and the region labels are global, so bands would not be identical. Splitting by independent computation (one sum each, one region each) is exact.
- **Serial speed-ups** (all exact): the paper sums share one luma sum and one count (6 box sums, not 8, with their loops hoisted); a gray pixel as light as its paper leaves the mask loop at once and a colored one skips the Lab distance unless it is dark enough to matter; `lineart_region_mask` fills its grid by rows and `lineart_fill_holes` reads the region grid instead of a second copy; `lineart_blur` has no edge tests in its inner loops. Serial alone: about 20 to 40 percent less time.
- **The editor.** The Trace image dialog runs the trace in one `task`, and a task can't start tasks, so there the trace is serial (faster than before, not fanned out). `linkz trace` passes `task` along (`--parallel false` turns it off). Language need logged: tasks starting tasks.
