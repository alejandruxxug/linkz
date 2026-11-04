# B062 – Trace: weak outlines on belts and thin colored bands

Model: Sonnet. Rules: `_v11-rules.md`. Read only: `lib/trace/lineart.liz` (the pipeline order), then the stages you touch (`lineart_mask`, `lineart_fills`, `lineart_outline`), `tests/trace/lineart*`, and `docs/Art/References/trace/`.

On the Wikipe-tan references, belts, straps and other thin colored bands come out weak: broken or missing ink, and fills that bleed across the band. Find out, by measuring which stage loses them, why that happens: the mask threshold or hysteresis on colored lines, the thinning, line_min, or the fill growth.

- Fix it in the stage that loses them, without hurting the other images. Score all 6 references before and after with the existing render difference (mean color error, ink IoU), and put a table in the report.
- Add a synthetic test: a thin colored band crossing a fill, outlined in ink.
- Save before/after shots in `docs/Art/Checks/trace/belts/`.
- Don't touch `lib/ui`.
