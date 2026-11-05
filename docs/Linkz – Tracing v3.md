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
