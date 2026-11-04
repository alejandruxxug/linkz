# B017 + B028 – Palette sidebar, color popover, color wheel and harmonies (v1.1, wave 2)

Design: R002 palette sidebar + color popover frames, R003-5 color wheel + harmonies; Style Guide R002/R003 rows "palette sidebar", "color popover", "Color wheel", "Harmony"; the approval note (wheel in the popover; ring/square toggle in the sidebar and fill editor, ring by default; harmony icons get tooltips; "color content" exception). Rules: `docs/Batches/_v11-rules.md`.

## Build
- **Color variables** are `let colors = {"sky": #f2a65a, …}` at the top of the drawing (create it when first needed, after the header comments). A variable can hold a gradient (`colors.dusk`), interchangeably. Picking a variable writes `colors.sky` (a live reference); editing a variable changes every shape that uses it.
- **Palette sidebar** (right, 106 wide, collapses to a 14 px rail; rail by default under 1100 px window width): the picker (ring + triangle by default, square toggle), recent colors, VARIABLES (swatches incl. gradients, rename, delete with "used by N shapes" refusal unless confirmed), "+ new variable".
- **Color popover** next to any clicked color (a code-row swatch, a props field, the fill editor's stop color): variables first, recent, the wheel, hex, alpha, old/new, "save as variable" (also wire B016's dim "Save as variable…" in the fill editor and its swatch row to the variables).
- **Color wheel**: hue ring + saturation/value triangle (Krita/Painter style), keyboard accessible, exact round-trip with hex.
- **Harmonies**: complementary, split-complementary, analogous, triadic, tetradic/square, monochrome; markers on the ring, swatches clickable, "+ add all" writes one `colors.<base>_2…` entry per swatch in one undo step. Color math in its own file (`lib/ui/color_harmony.liz` or `lib/draw/`), tested (complementary of #3fa34d is hue + 180°).
- Every edit one undo step; en + es.
