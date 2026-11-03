# B037 – Canvas batching and line-count edits (M10, Linkz side)

Repo: `linkz` · Lane: lizard (`lib/ui/view_scene.liz`, `lib/gfx/tiles.liz`, `lib/gfx/raster.liz`, `lib/ui/doc_*.liz`, `editor_*`) + a host span-batch primitive only if the profile proves it · Needs: B036

## Why (B036's report, ms per tick)
- Keystroke at 2k is 29 ms, at 10k 101 ms: **canvas 19.7 / 87.7 ms**. The benchmark turns a radius 20 into 201, so ~400×400 px repaint, which today paints every shape once **per tile** it touches (~2,700 shape paints; `fb.spans` called 2,862 times).
- Full repaint 22.7 ms at 2k, 88.5 ms at 10k.
- Edits that change the line count (Enter, deleting a line) still rebuild every statement and row after the edit (O(n) at 10k); not in the benchmarks yet.
- ~6 ms per tick outside the phases: two whole-document line diffs per edit (`lines_common`, `rerun_changed_line`, could share one), 4–5 whole-file `doc_text` joins per edit, `view_reveal`, panel painting.
- Broken code falls back to a full scan on every keystroke (~30 ms at 2k).

## Build
1. **Paint per dirty region, not per tile**: union the dirty boxes, query the index once per region, paint each shape once into the region (clipped), then split into tiles; or send a whole tile's spans to the host in one call. Pick by measuring. Keep the golden test (`TestSceneMatchesFullRepaint`) byte-identical.
2. **Line-count edits incremental**: shift line numbers lazily (an offset per region / a line-delta list) instead of rebuilding statements and rows after the edit. Add benchmark cases: Enter in the middle, delete a line, paste 20 lines, at 2k and 10k.
3. **One diff per edit**, shared by the splice and the rerun; **no whole-file joins** on the hot path (keep the text joined once per version in the cache).
4. **Broken code**: when the text doesn't parse, rescan only the statement around the edit when the rest is unchanged (the previous good parse plus the broken region), or at least cache the last full scan per version.

## Done when
2k: every benchmark case (incl. the new line-count ones) < 16 ms; 10k: < 50 ms. Golden and randomized splice tests green, VM and `LINKZ_ENGINE=tree`, smoke green. Before/after table with per-phase ms.

## Rules
Clean code; LEP 0001 style (operators at line ends); lizard decides pixels (a Go primitive only to move spans in bulk, justified by numbers). Testers: temp HOME, `--no-splash`. Don't commit. linkz builds against `../.lizard-stable` (c703e44). Report language-level slowness with numbers.
